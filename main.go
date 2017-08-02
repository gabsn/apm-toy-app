package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
)

func main() {
	router := newRouter()
	http.HandleFunc("/redis", router.redisHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Router struct {
	redis *redis.Client
}

func newRouter() *Router {
	return &Router{
		redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDRESS"),
		}),
	}
}

func (r *Router) redisHandler(w http.ResponseWriter, req *http.Request) {
	n := r.redis.Incr("counter").Val()
	fmt.Fprintf(w, "/redis received %v hits.", n)
}
