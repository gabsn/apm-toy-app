package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/DataDog/dd-trace-go/tracer"
	redistrace "github.com/DataDog/dd-trace-go/tracer/contrib/go-redis"
	"github.com/DataDog/dd-trace-go/tracer/contrib/net/httptrace"
	sqltrace "github.com/DataDog/dd-trace-go/tracer/contrib/sqltraced"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

func main() {
	r := newRouter()
	r.HandleFunc("/", r.handler)
	traceHandler := httptrace.NewHandler(r, "web-backend", tracer.DefaultTracer)
	log.Fatal(http.ListenAndServe(":8080", traceHandler))
}

type Router struct {
	*mux.Router
	redis *redistrace.TracedClient
	pg    *sql.DB
}

func newRouter() *Router {
	r := mux.NewRouter()

	redis := redistrace.NewTracedClient(&redis.Options{
		Addr: "redis:6379",
	}, tracer.DefaultTracer, "redis")

	pg, err := sqltrace.OpenTraced(&pq.Driver{}, "host=postgres user=postgres dbname=postgres sslmode=disable", "postgres", tracer.DefaultTracer)
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
