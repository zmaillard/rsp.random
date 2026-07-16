package services

import (
	"bytes"
	"context"
	"encoding/gob"
	"log/slog"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"rsp.random/db"
)

type UpdateCounterProcess func(ctx context.Context) error

type CounterService interface {
	UpdateCounts(ctx context.Context) error
	UpdateData(ctx context.Context) error
}

type counterService struct {
	querier db.Querier
	kvStore *badger.DB
}

func NewCounterService(querier db.Querier, kvStore *badger.DB) CounterService {
	return &counterService{
		querier: querier,
		kvStore: kvStore,
	}
}
func (c *counterService) UpdateData(ctx context.Context) error {
	err := c.kvStore.Update(func(txn *badger.Txn) error {
		allSigns, berr := c.querier.GetSigns(ctx)

		if berr != nil {
			return berr
		}

		totalCount := len(allSigns)
		groupings := make(map[string][]string)
		for idx, sign := range allSigns {
			groupArr := groupings[sign.StateSlug]
			groupArr = append(groupArr, sign.Imageid)
			groupings[sign.StateSlug] = groupArr

			if sign.PlaceSlug != nil {
				ps := *sign.PlaceSlug
				placeArr := groupings[ps]
				placeArr = append(placeArr, sign.Imageid)
				groupings[ps] = placeArr
			}

			if sign.CountySlug != nil {
				cs := *sign.CountySlug
				countyArr := groupings[cs]
				countyArr = append(countyArr, sign.Imageid)
				groupings[cs] = countyArr
			}

			var signBuff bytes.Buffer
			enc := gob.NewEncoder(&signBuff)
			eerr := enc.Encode(sign.Imageid)
			if eerr != nil {
				return eerr
			}

			berr = txn.Set([]byte(strconv.Itoa(idx)), signBuff.Bytes())
			if berr != nil {
				return berr
			}
		}
		var totalCountBuff bytes.Buffer
		enc := gob.NewEncoder(&totalCountBuff)
		eerr := enc.Encode(strconv.Itoa(totalCount))
		if eerr != nil {
			return eerr
		}

		berr = txn.Set([]byte("totalCount"), totalCountBuff.Bytes())
		if berr != nil {
			return berr
		}

		for key, value := range groupings {
			var signBuff bytes.Buffer
			enc := gob.NewEncoder(&signBuff)
			eerr := enc.Encode(value)
			if eerr != nil {
				return eerr
			}

			berr = txn.Set([]byte(key), signBuff.Bytes())
			if berr != nil {
				return berr
			}

		}

		return nil

	})

	return err

}

func (c *counterService) UpdateCounts(ctx context.Context) error {
	slog.Warn("here")
	err := c.kvStore.Update(func(txn *badger.Txn) error {
		total, berr := c.querier.GetTotalCounts(ctx)
		if berr != nil {
			return berr
		}
		berr = txn.Set([]byte("allsigns"), []byte(strconv.Itoa(int(total))))
		if berr != nil {
			return berr
		}

		counties, berr := c.querier.GetCountyCounts(ctx)
		for _, county := range counties {
			berr = txn.Set([]byte(county.CountySlug), []byte(strconv.Itoa(int(county.Counter))))
			if berr != nil {
				return berr
			}
		}

		states, berr := c.querier.GetStateCounts(ctx)
		for _, state := range states {
			berr = txn.Set([]byte(state.Slug.String), []byte(strconv.Itoa(int(state.Counter))))
			if berr != nil {
				return berr
			}
		}

		places, berr := c.querier.GetPlaceCounts(ctx)
		for _, place := range places {
			berr = txn.Set([]byte(place.PlaceSlug), []byte(strconv.Itoa(int(place.Counter))))
			if berr != nil {
				return berr
			}
		}
		return nil
	})

	return err
}
