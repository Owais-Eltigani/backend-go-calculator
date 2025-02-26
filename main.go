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

func main() {
	// only on test.
	fmt.Println("backend calculator using golan.")

	//
	router := http.NewServeMux()
	port := ":5000"

	router.HandleFunc("POST /login", Login) //! login first using owais:123 to access multi & div
	router.Handle("GET /add", BasicAuth(Add))
	router.Handle("GET /sub", BearTokenMiddleware(Sub))
	router.Handle("GET /multi", JWTverifiaction(Multi))
	router.HandleFunc("GET /div", JWTverifiaction(Div))

	// these endpoints use rate limiter
	// router.HandleFunc("GET /add", middleware(Add))
	// router.HandleFunc("GET /sub", middleware(Sub))
	// router.HandleFunc("GET /multi", middleware(Multi))
	// router.HandleFunc("GET /div", Middleware(Div))

	server := &http.Server{
		Handler:      router,
		Addr:         port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	fmt.Println("router...")
	log.Fatal(server.ListenAndServe())

}
