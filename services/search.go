package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v4"
	"rsp.random/config"
)

type SearchResult interface {
	GetRedirectUrl(c *config.Config) (string, error)
}

type SearchService interface {
	RandomSign() (SearchResult, error)
	RandomSignByState(state string) (SearchResult, error)
	RandomSignByCounty(state string, county string) (SearchResult, error)
	RandomSignByPlace(state string, place string) (SearchResult, error)
}

type searchService struct {
	httpClient *http.Client
	badgerDb   *badger.DB
	config     *config.Config
}

type searchBody struct {
	Query  string `json:"q"`
	Filter string `json:"filter"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func NewSearchService(httpClient *http.Client, badgerDb *badger.DB, cfg *config.Config) SearchService {
	return &searchService{httpClient: httpClient, badgerDb: badgerDb, config: cfg}
}

func (s searchService) getSearchBody(filter string, offset int) ([]byte, error) {
	b := searchBody{
		Query:  "",
		Filter: filter,
		Limit:  1,
		Offset: offset,
	}

	return json.Marshal(b)
}

func (s searchService) RandomSign() (SearchResult, error) {
	count := 0
	err := s.badgerDb.View(func(txn *badger.Txn) error {
		item, berr := txn.Get([]byte("allsigns"))
		if berr != nil {
			return berr
		}
		scount := item.String()
		count, berr = strconv.Atoi(scount)
		return berr
	})

	if count == 0 {
		return nil, errors.New("no signs found")
	}
	offset := rand.N[int](count)
	url, err := s.config.GetSearchUrl()
	if err != nil {
		return nil, err
	}
	searchBody, err := s.getSearchBody("", offset)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(searchBody))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.config.SearchApiKey))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var searchResult searchResult
	err = json.Unmarshal(bodyBytes, &searchResult)
	return searchResult, err
}

func (s searchService) RandomSignByState(state string) (SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchService) RandomSignByCounty(state string, county string) (SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

func (s searchService) RandomSignByPlace(state string, place string) (SearchResult, error) {
	//TODO implement me
	panic("implement me")
}

type searchResult struct {
	Hits []struct {
		Id  string `json:"id"`
		Geo struct {
			Lng float64 `json:"lng"`
			Lat float64 `json:"lat"`
		} `json:"_geo"`
		Country struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"country"`
		County struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"county"`
		DateTaken   string `json:"dateTaken"`
		Description string `json:"description"`
		Highways    []struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"highways"`
		Place struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"place"`
		State struct {
			Name string `json:"name"`
			Slug string `json:"slug"`
		} `json:"state"`
		Tags       interface{} `json:"tags"`
		Title      string      `json:"title"`
		Url        string      `json:"url"`
		Quality    int         `json:"quality"`
		DateTaken1 time.Time   `json:"date_taken"`
	} `json:"hits"`
	Query              string `json:"query"`
	ProcessingTimeMs   int    `json:"processingTimeMs"`
	Limit              int    `json:"limit"`
	Offset             int    `json:"offset"`
	EstimatedTotalHits int    `json:"estimatedTotalHits"`
	RequestUid         string `json:"requestUid"`
}

func (s searchResult) GetRedirectUrl(c *config.Config) (string, error) {
	return url.JoinPath(c.BaseUrl, "sign", s.Hits[0].Id)
}
