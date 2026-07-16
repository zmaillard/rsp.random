package services

import (
	"bytes"
	"context"
	"encoding/gob"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"rsp.random/db"
)

type UpdateCounterProcess func(ctx context.Context) error

type CounterService interface {
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
