package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	short "github.com/lrweck/shawty/shortener"
	"github.com/pkg/errors"
)

type redisRepo struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Creates a new Redis repository of data using a user supplied URL.
func NewRedisRepository(redisURL string) (short.RedirectRepository, error) {
	repo := &redisRepo{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *redisRepo) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

// Finds in storage the redirect the user queried by code
func (r *redisRepo) Find(code string) (*short.Redirect, error) {
	redirect := &short.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()

	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(short.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	redirect.CreatedAt = createdAt
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	return redirect, nil
}

// Stores the user supplied url to a redirect
func (r *redisRepo) Store(redirect *short.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
