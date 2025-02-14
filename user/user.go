package user

import (
	"99-api-public/config"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

var ctx = context.Background()

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt int    `json:"created_at"`
	UpdatedAt int    `json:"updated_at"`
}

func GetAllUsers() ([]User, error) {
	msUrl := config.GetEnv("USER_SERVICE_URL", "http://localhost:8080") + "/users"
	resp, err := http.Get(msUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []User
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]json.RawMessage
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	body, ok := result["users"]
	if !ok {
		return nil, errors.New("users key not found in response")
	}

	if err := json.Unmarshal(result["users"], &users); err != nil {
		return nil, err
	}

	return users, nil
}

func CacheUsers(rdb *redis.Client, users []User) error {

	// if len(users) == 0 {
	// 	return errors.New("users is empty")
	// }

	for _, user := range users {
		userData, err := json.Marshal(user)
		if err != nil {
			return err
		}
		err = rdb.Set(ctx, "user:"+strconv.Itoa(user.ID), userData, 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetUserFromCache(rdb *redis.Client, id int) (*User, error) {
	userData, err := rdb.Get(ctx, "user:"+strconv.Itoa(id)).Result()
	if err == redis.Nil {
		return nil, errors.New("user not found")
	} else if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal([]byte(userData), &user); err != nil {
		return nil, err
	}

	return &user, nil
}

func SetUserToCache(rdb *redis.Client, user User) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = rdb.Set(ctx, "user:"+strconv.Itoa(user.ID), userData, 0).Err()
	if err != nil {
		return err
	}
	return nil
}
