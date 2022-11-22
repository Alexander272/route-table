package redis

import (
	"context"
	"encoding/json"

	"github.com/Alexander272/route-table/internal/models"
	"github.com/Alexander272/route-table/pkg/logger"
	"github.com/go-redis/redis/v8"
)

func UnMarshalBinary(str string) models.SessionData {
	var data models.SessionData
	json.Unmarshal([]byte(str), &data)
	return data
}

type SessionRepo struct {
	client redis.Cmdable
}

func NewSessionRepo(client redis.Cmdable) *SessionRepo {
	return &SessionRepo{
		client: client,
	}
}

func (r *SessionRepo) Create(ctx context.Context, sessionName string, data models.SessionData) error {
	if err := r.client.Set(ctx, sessionName, data, data.Exp).Err(); err != nil {
		return err
	}
	return nil
}

func (r *SessionRepo) Get(ctx context.Context, sessionName string) (data models.SessionData, err error) {
	cmd := r.client.Get(ctx, sessionName)
	if cmd.Err() != nil {
		if cmd.Err() == redis.Nil {
			return data, models.ErrSessionEmpty
		}
		logger.Error(cmd.Err())
		return data, cmd.Err()
	}

	str, err := cmd.Result()
	if err != nil {
		return data, err
	}
	return UnMarshalBinary(str), nil
}

func (r *SessionRepo) GetDel(ctx context.Context, sessionName string) (data models.SessionData, err error) {
	cmd := r.client.GetDel(ctx, sessionName)
	if cmd.Err() != nil {
		logger.Debug(cmd.Err())
		return data, cmd.Err()
	}

	str, err := cmd.Result()
	if err != nil {
		return data, err
	}
	return UnMarshalBinary(str), nil
}

func (r *SessionRepo) Remove(ctx context.Context, sessionName string) error {
	if err := r.client.Del(ctx, sessionName).Err(); err != nil {
		return err
	}
	return nil
}
