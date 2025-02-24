package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
	_ "golang.org/x/time/rate"
)

type Operands struct {
	Operand1 string `json:"a"`
	Operand2 string `json:"b"`
}

var (
	// limiter  = rate.NewLimiter(rate.Limit(2), 0) //?	for time interval rate limiting.
	vistiors = make(map[string]*rate.Limiter)
	mutex    sync.Mutex
)

var limiter = rate.NewLimiter(2, 3)
var visitors = make(map[string]*rate.Limiter, 0)

func checkIpAddressLimit(ip string) *rate.Limiter {
	mutex.Lock()
	defer mutex.Unlock()

	v, exist := visitors[ip]

	if !exist {

		v = rate.NewLimiter(rate.Limit(3), 2)
	}
	return v
}

func middleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr
		limit := checkIpAddressLimit(ip)

		if !limit.Allow() {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func main() {
	// only on test.
	fmt.Println("backend calculator using golan.")

	//
	router := http.NewServeMux()
	port := ":5000"

	router.HandleFunc("GET /add", Add)
	router.HandleFunc("GET /sub", Sub)
	router.HandleFunc("GET /multi", Multi)
	router.HandleFunc("GET /div", Div)

	server := &http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	fmt.Println("router...")
	log.Fatal(server.ListenAndServe())

}
