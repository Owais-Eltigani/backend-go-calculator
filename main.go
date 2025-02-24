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

		// v = rate.NewLimiter(0, 0)  //? for testing.
		v = rate.NewLimiter(rate.Limit(3), 2) //? 3 requests for user and 2 burst.
	}
	return v
}

func middleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr
		limit := checkIpAddressLimit(ip)
		fmt.Println("ip: ", ip, "limit: ", limit.Allow())

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

	router.HandleFunc("GET /add", middleware(Add))
	router.HandleFunc("GET /sub", middleware(Sub))
	router.HandleFunc("GET /multi", middleware(Multi))
	router.HandleFunc("GET /div", middleware(Div))

	server := &http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	fmt.Println("router...")
	log.Fatal(server.ListenAndServe())

}
