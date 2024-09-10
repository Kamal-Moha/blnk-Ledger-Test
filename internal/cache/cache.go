/*
Copyright 2024 Blnk Finance Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	redis_db "github.com/jerry-enebeli/blnk/internal/redis-db"

	"github.com/jerry-enebeli/blnk/config"

	"github.com/go-redis/cache/v9"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, data interface{}) error
	Delete(ctx context.Context, key string) error
}

type RedisCache struct {
	cache *cache.Cache
}

func NewCache() (Cache, error) {
	cfg, err := config.Fetch()
	if err != nil {
		return nil, err
	}
	ca, err := newRedisCache([]string{fmt.Sprintf("redis://%s", cfg.Redis.Dns)})
	if err != nil {
		return nil, err
	}
	return ca, nil
}

const cacheSize = 128000

func newRedisCache(addresses []string) (*RedisCache, error) {
	client, err := redis_db.NewRedisClient(addresses)
	if err != nil {
		return nil, err
	}

	c := cache.New(&cache.Options{
		Redis:      client.Client(),
		LocalCache: cache.NewTinyLFU(cacheSize, 1*time.Minute),
	})

	r := &RedisCache{cache: c}

	return r, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, data interface{}, ttl time.Duration) error {
	return r.cache.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: data,
		TTL:   ttl,
	})
}

func (r *RedisCache) Get(ctx context.Context, key string, data interface{}) error {
	err := r.cache.Get(ctx, key, &data)
	if errors.Is(err, cache.ErrCacheMiss) {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.cache.Delete(ctx, key)
}
