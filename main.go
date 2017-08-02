package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func main() {
	router := newRouter()
	http.HandleFunc("/redis", router.handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Router struct {
	redis *redis.Client
	pg    *sql.DB
}

func newRouter() *Router {
	redis := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	pg, err := sql.Open("postgres", "host=postgres user=postgres dbname=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	return &Router{redis, pg}
}

func (r *Router) handler(w http.ResponseWriter, req *http.Request) {
	var name, population string

	// Count the number of hits on this enpoint
	n := r.redis.Incr("counter").Val()

	// Get the city associated to this number of hits
	err := r.pg.QueryRow("SELECT name, population FROM city WHERE id = $1", n%20+1).Scan(&name, &population)
	if err != nil {
		log.Print(err)
		return
	}

	// Return the result
	fmt.Fprintf(w, "(%v hits) - City: %v, %v inhabitants", n, name, population)
}
