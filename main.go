package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	sqltrace "github.com/DataDog/dd-trace-go/contrib/database/sql"
	redistrace "github.com/DataDog/dd-trace-go/contrib/go-redis/redis"
	muxtrace "github.com/DataDog/dd-trace-go/contrib/gorilla/mux"

	"github.com/go-redis/redis"
	"github.com/lib/pq"
)

func main() {
	r := newRouter()
	r.HandleFunc("/", r.handler)
	log.Fatal(http.ListenAndServe(":8080", r))
}

type Router struct {
	*muxtrace.Router
	redis *redistrace.Client
	pg    *sql.DB
}

func newRouter() *Router {
	r := muxtrace.NewSereMux("web-api")

	redis := redistrace.NewClient(&redis.Options{
		Addr: "redis:6379",
	}, "redis")

	pg, err := sqltrace.Open(&pq.Driver{}, "host=postgres user=postgres dbname=postgres sslmode=disable", "postgres")
	if err != nil {
		panic(err)
	}

	return &Router{r, redis, pg}
}

func (r *Router) handler(w http.ResponseWriter, req *http.Request) {
	var name, population string

	// Link this call to redis to the previous to the request
	r.redis.SetContext(req.Context())

	// Count the number of hits on this enpoint
	n := r.redis.Incr("counter").Val()

	// Get the city associated to this number of hits
	err := r.pg.QueryRowContext(req.Context(), "SELECT name, population FROM city WHERE id = $1", n%20+1).Scan(&name, &population)
	if err != nil {
		log.Print(err)
		return
	}

	// Return the name of the city and its population
	fmt.Fprintf(w, "(%v hits) - City: %v, %v inhabitants", n, name, population)
}
