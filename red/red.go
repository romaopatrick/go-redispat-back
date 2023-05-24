package red

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"redispat/repo"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func newClusterClient(c *repo.Connection) *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    c.Addresses,
		Password: c.Password,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		ClientName:      c.Name,
		ConnMaxIdleTime: 10 * time.Second,
		NewClient: func(opt *redis.Options) *redis.Client {
			return redis.NewClient(opt)
		},
		MaxRetries: 3,
	})
}

func Set(ctx context.Context, c *repo.Connection, key string, val interface{}) (err error) {
	ccl := newClusterClient(c)

	var b []byte

	switch t := val.(type) {
	case int:
		b = []byte(strconv.Itoa(t))
	case int64:
		b = []byte(strconv.FormatInt(t, 10))
	default:
		b, err = json.Marshal(t)
		if err != nil {
			log.Panicf("[redis] - %s", err.Error())
			return
		}
	}

	err = ccl.Set(ctx, key, b, redis.KeepTTL).Err()
	return
}

func Get(ctx context.Context, c *repo.Connection, key string, pointer *interface{}) {
	cl := newClusterClient(c)

	re, err := cl.Get(ctx, key).Result()

	if err == redis.Nil {
		return
	}

	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal([]byte(re), pointer); err != nil {
		*pointer = re
	}
}

func ListKeys(ctx context.Context, c *repo.Connection, contains string) []string {
	cl := newClusterClient(c)

	cmd := cl.Keys(ctx, contains)
	re, err := cmd.Result()

	if err != nil && err != redis.Nil {
		panic(err)
	}

	if len(re) > 50 {
		return re[:50]
	}
	return re
}

func DeleteKey(ctx context.Context, c *repo.Connection, key string) error {
	cl := newClusterClient(c)

	_, err := cl.Del(ctx, key).Result()
	if err != nil && err != redis.Nil {
		panic(err)
	}

	return nil
}
