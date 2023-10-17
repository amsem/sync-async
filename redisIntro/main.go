package main

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
)

func main()  {
    client := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        Password: "",
        DB: 0,
    })
    ctx := context.Background()
    p, _ := client.Ping(ctx).Result()
    fmt.Println(p)
}
