package services

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"
	"slices"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"rsp.random/config"
)

type SearchResultSlim struct {
	ImageId      string `json:"imageId"`
	HasProcessed bool   `json:"hasProcessed"`
}

type SearchResult interface {
	GetRedirectUrl(c *config.Config) (string, error)
	GetIdOnly() SearchResultSlim
}

type SearchService interface {
	RandomSign() (SearchResult, error)
	RandomSignByState(state string) (SearchResult, error)
	RandomSignByCounty(state string, county string) (SearchResult, error)
	RandomSignByPlace(state string, place string) (SearchResult, error)
	HasProcessed(imageId string) (bool, error)
}

type searchService struct {
	httpClient *http.Client
	badgerDb   *badger.DB
	config     *config.Config
}

func NewSearchService(httpClient *http.Client, badgerDb *badger.DB, cfg *config.Config) SearchService {
	return &searchService{httpClient: httpClient, badgerDb: badgerDb, config: cfg}
}

func (s SearchResultSlim) GetRedirectUrl(c *config.Config) (string, error) {
	return url.JoinPath(c.BaseUrl, "sign", s.ImageId)
}

func (s SearchResultSlim) GetIdOnly() SearchResultSlim {
	return s
}

func (s *searchService) RandomSign() (SearchResult, error) {
	var valCopy []byte
	err := s.badgerDb.View(func(txn *badger.Txn) error {
		key := "totalCount"
		item, berr := txn.Get([]byte(key))
		if berr != nil {
			return berr
		}
		berr = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return berr
	})
	if err != nil {
		return nil, err
	}

	var totalCount string
	err = gob.NewDecoder(bytes.NewReader(valCopy)).Decode(&totalCount)
	if err != nil {
		return nil, err
	}
	iTotalCount, err := strconv.Atoi(totalCount)
	if err != nil {
		return nil, err
	}

	offset := rand.N[int](iTotalCount)

	var resCopy []byte
	err = s.badgerDb.View(func(txn *badger.Txn) error {
		key := strconv.Itoa(offset)
		item, berr := txn.Get([]byte(key))
		if berr != nil {
			return berr
		}
		berr = item.Value(func(val []byte) error {
			resCopy = append([]byte{}, val...)
			return nil
		})
		return berr
	})
	if err != nil {
		return nil, err
	}
	var imageId string
	err = gob.NewDecoder(bytes.NewReader(resCopy)).Decode(&imageId)
	if err != nil {
		return nil, err
	}

	hasProcessed, err := s.HasProcessed(imageId)
	if err != nil {
		return nil, err
	}

	return SearchResultSlim{ImageId: imageId, HasProcessed: hasProcessed}, err
}

func (s *searchService) RandomSignByState(state string) (SearchResult, error) {
	var valCopy []byte
	err := s.badgerDb.View(func(txn *badger.Txn) error {
		item, berr := txn.Get([]byte(state))
		if berr != nil {
			return berr
		}
		berr = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return berr
	})
	if err != nil {
		return nil, err
	}

	var resultArray []string
	err = gob.NewDecoder(bytes.NewReader(valCopy)).Decode(&resultArray)
	if err != nil {
		return nil, err
	}

	offset := rand.N[int](len(resultArray))
	imageId := resultArray[offset]
	hasProcessed, err := s.HasProcessed(imageId)
	if err != nil {
		return nil, err
	}
	return SearchResultSlim{ImageId: imageId, HasProcessed: hasProcessed}, err

}

func (s *searchService) RandomSignByCounty(state string, county string) (SearchResult, error) {
	key := fmt.Sprintf("%s_%s", state, county)
	var valCopy []byte
	err := s.badgerDb.View(func(txn *badger.Txn) error {
		item, berr := txn.Get([]byte(key))
		if berr != nil {
			return berr
		}
		berr = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return berr
	})
	if err != nil {
		return nil, err
	}

	var resultArray []string
	err = gob.NewDecoder(bytes.NewReader(valCopy)).Decode(&resultArray)
	if err != nil {
		return nil, err
	}

	offset := rand.N[int](len(resultArray))
	imageId := resultArray[offset]
	hasProcessed, err := s.HasProcessed(imageId)
	if err != nil {
		return nil, err
	}
	return SearchResultSlim{ImageId: imageId, HasProcessed: hasProcessed}, err
}

func (s *searchService) RandomSignByPlace(state string, place string) (SearchResult, error) {
	key := fmt.Sprintf("%s_%s", state, place)
	var valCopy []byte
	err := s.badgerDb.View(func(txn *badger.Txn) error {
		item, berr := txn.Get([]byte(key))
		if berr != nil {
			return berr
		}
		berr = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return berr
	})
	if err != nil {
		return nil, err
	}

	var resultArray []string
	err = gob.NewDecoder(bytes.NewReader(valCopy)).Decode(&resultArray)
	if err != nil {
		return nil, err
	}

	offset := rand.N[int](len(resultArray))
	imageId := resultArray[offset]
	hasProcessed, err := s.HasProcessed(imageId)
	if err != nil {
		return nil, err
	}
	return SearchResultSlim{ImageId: imageId, HasProcessed: hasProcessed}, err
}

func (s *searchService) HasProcessed(imageId string) (bool, error) {
	var valCopy []byte
	err := s.badgerDb.View(func(txn *badger.Txn) error {
		item, berr := txn.Get([]byte("processed"))
		if berr != nil {
			return berr
		}
		berr = item.Value(func(val []byte) error {
			valCopy = append([]byte{}, val...)
			return nil
		})
		return berr
	})
	if err != nil {
		return false, err
	}

	var processed []string
	err = gob.NewDecoder(bytes.NewReader(valCopy)).Decode(&processed)
	if err != nil {
		return false, err
	}

	if slices.Contains(processed, imageId) {
		return true, nil
	}

	return false, nil
}
