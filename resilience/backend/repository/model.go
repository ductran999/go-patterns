package repository

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

var requestGroup singleflight.Group
var ErrDB = errors.New("DB Error")

type modelRepository struct {
	cache      sync.Map
	allowCache time.Time
}

func NewModelRepository() *modelRepository {
	return &modelRepository{
		allowCache: time.Now().Add(5 * time.Second),
	}
}

func (r *modelRepository) SimulateQueryDB(apiKey string) (string, error) {
	if val, ok := r.cache.Load(apiKey); ok {
		return val.(string), nil
	}

	res, err, shared := requestGroup.Do(apiKey, func() (any, error) {
		if val, ok := r.cache.Load(apiKey); ok {
			return val.(string), nil
		}

		slog.Info(">>> DO QUERY DB (SELECT) CHO KEY: " + apiKey)
		time.Sleep(100 * time.Millisecond)

		data := "query data for " + apiKey
		if time.Now().After(r.allowCache) {
			r.cache.Store(apiKey, "cached data: "+data)
		}

		return data, nil
	})

	if err != nil {
		return "", err
	}

	if shared {
		return "reuse data" + res.(string), nil
	}

	return res.(string), nil
}
