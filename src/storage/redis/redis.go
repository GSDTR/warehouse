package Redis

import (
	"github.com/go-redis/redis"
	"log"
	"strings"
)

const Redis_Pong  = "PONG"

var client *redis.Client

func Connect() error{
	client = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	is_connected := Test_connection()
	return is_connected
}

func Test_connection() error {
	pong, err := client.Ping().Result()
	if err != nil {
		return err
//		panic(err)
	}
	log.Println("Redis is active: ", pong)

	_, err = Get("key")
	if err != nil {
		log.Println("Redis not initialized")
	}

	err = Set("key", "value4")
	if err != nil {
		return err
	}

	_, err = Get("key")
	if err != nil {
		return err
	}

	if strings.Compare(pong, Redis_Pong) == 0 {
		return nil
	} else {
		return err
	}
}

func Set(key string, value interface{}) error {
	return client.Set(key, value, 0).Err()
}

func Get(key string) (string, error) {
	return client.Get(key).Result()
}
