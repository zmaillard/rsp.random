package services

import (
	"context"
	"strconv"

	"github.com/dgraph-io/badger/v4"
	"rsp.random/db"
)

type CounterService interface {
	UpdateCounts(ctx context.Context) error
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

func (c *counterService) UpdateCounts(ctx context.Context) error {
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
