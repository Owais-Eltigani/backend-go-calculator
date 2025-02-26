package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// struct type for generating jwt token, prove who the user who acclaims.
type claims struct {
	User string `json:"username"`
	jwt.RegisteredClaims
}

var jwtkey = []byte("my_secret_key") //? not working try to use my_secret_key
var jwtTokensArr []string            //? cookies jar DB.

//*==================================================================

// first middleware function to login user.
func JWTloggin(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//* ====================== basic auth ===================
		// 	user, pass, ok := r.BasicAuth()
		// 	if !ok {
		// 		fmt.Println("user name or password not correct")
		// 		http.Error(w, "user name or password not correct", http.StatusUnauthorized)
		// 		return
		// 	}
		// 	fmt.Println("user and pass", user, pass)
		// 	genToken, _ := generateJWT()
		// 	jwtTokensArr = append(jwtTokensArr, genToken)

		// 	next.ServeHTTP(w, r)
		// }

		// //* ====================== basic auth ===================

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "authorization must be Bearer token", http.StatusUnauthorized)
			return
		}

		// Generate new token
		genToken, err := generateJWT()
		if err != nil {
			http.Error(w, "failed to generate token", http.StatusInternalServerError)
			return
		}
		jwtTokensArr = append(jwtTokensArr, genToken)

		// Set token in response header
		w.Header().Set("Authorization", "Bearer "+genToken)

		next.ServeHTTP(w, r)
	}
}

func JWTverifiaction(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		//* code I wrote
		// // get the auth value from header
		// reqToken := r.Header.Get("Authorization")
		// fmt.Println("reqtokens: ", reqToken)

		// if reqToken == "" {
		// 	http.Error(w, "authorization header required", http.StatusUnauthorized)
		// 	return
		// }

		// token := strings.Split(reqToken, " ")[1] //? remember to use the [1] because it's an array.
		// claims := &claims{}                      //? to be filled by the parsewithClaims

		// tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// 	return jwtkey, nil
		// })

		// if err != nil {

		// 	fmt.Println("tkn is ", tkn)
		// 	if err == jwt.ErrSignatureInvalid {
		// 		fmt.Println("invalid user")
		// 		http.Error(w, "invalid user", http.StatusUnauthorized)
		// 		return
		// 	}

		// 	fmt.Println("bad request", err)
		// 	http.Error(w, "bad request", http.StatusBadRequest)
		// 	return
		// }

		// if !tkn.Valid {
		// 	fmt.Println("token not valid")
		// 	http.Error(w, "token not valid", http.StatusUnauthorized)
		// 	return
		// }

		// next.ServeHTTP(w, r)

		// * claude code

		// Get authorization header
		authHeader := r.Header.Get("Authorization")
		fmt.Println("auth header: ", authHeader)
		if authHeader == "" {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		// Extract the token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "authorization must be Bearer token", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Parse and validate JWT
		claims := &claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method is what we expect
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtkey, nil
		})

		// Handle parsing errors
		if err != nil {
			fmt.Println("JWT Error:", err)
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "invalid signature", http.StatusUnauthorized)
			} else {
				http.Error(w, "invalid token", http.StatusUnauthorized)
			}
			return
		}

		// Check if token is valid
		if !token.Valid {
			http.Error(w, "token not valid", http.StatusUnauthorized)
			return
		}

		// Add user info to context and proceed
		ctx := context.WithValue(r.Context(), "user", claims.User)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// helper function for JWTverification to generate the jwt.
func generateJWT() (string, error) {

	expirationTime := time.Now().Add(2 * time.Minute)
	claims := &claims{
		User: "owais",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return tokenString, nil
}

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		fmt.Println("wrong method type")
		http.Error(w, "wrong method", http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		fmt.Println("body is empty")
		http.Error(w, "body is empty", http.StatusNoContent)
		return
	}

	var creds Credentials
	json.NewDecoder(r.Body).Decode(&creds)

	if creds.User != "owais" || creds.Password != "123" {
		fmt.Println("wrong credentials")
		http.Error(w, "wrong credentials", http.StatusUnauthorized)
		return
	}

	token, _ := generateJWT()
	fmt.Println("token was created successfully.", token)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token)
	json.NewEncoder(w).Encode(token)
}

type Credentials struct {
	User     string `json:"username"`
	Password string `json:"password"`
}
