package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// rate limiting using number of requests from ip address
func CheckIpAddressLimit(ip string) *rate.Limiter {

	mutex.Lock() //? look the DB visitors.
	defer mutex.Unlock()

	v, exist := visitors[ip]
	if !exist {

		// v = rate.NewLimiter(0, 0)  //? for testing.
		v = rate.NewLimiter(rate.Limit(3), 2) //? 3 requests for user and 2
	}
	return v
}

// read the user ip and compare number of request left.
func Middleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ip := r.RemoteAddr
		limit := CheckIpAddressLimit(ip)
		fmt.Println("ip: ", ip, "limit: ", limit.Allow())

		if !limit.Allow() {
			http.Error(w, "too many requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// basic auth middleware function to check if user exist in DB or not.
func BasicAuth(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		username, password, ok := r.BasicAuth()
		if !ok {
			fmt.Println("not possiable")
			http.Error(w, "not authenticated", http.StatusUnauthorized)
			return
		}
		fmt.Println(username, password, ok)
		next.ServeHTTP(w, r)
	}
}

// implement token bearer auth
func BearTokenMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()
		if !ok {
			fmt.Println("not authorized")
			http.Error(w, "token not authorized", http.StatusUnauthorized)
			return
		}
		fmt.Println(user, pass)

		// generate hexadecimal string as string
		getToken, _ := generateToken(10)
		cookie := &http.Cookie{
			Value:    getToken,
			Path:     "/sub",
			Name:     "auth-token",
			HttpOnly: true,
			MaxAge:   24 * 3600, // 224 hours in seconds
			Expires:  time.Now().Add(24 * time.Hour),
		}

		http.SetCookie(w, cookie)
		fmt.Println("token ", getToken, "cookie: \n", cookie)
		next.ServeHTTP(w, r)

	}

}

// helper function for bearTokenMiddleware to generate the hexadecimal string.
func generateToken(hash int) (string, error) {

	token := make([]byte, hash)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}
