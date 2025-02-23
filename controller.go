package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		fmt.Println("GET/add: the request body is empty.")
		w.Write([]byte("request body is empty"))
		return
	}

	var operands Operands
	json.NewDecoder(r.Body).Decode(&operands)

	if operands.isEmpty("GET/add") {
		fmt.Println(operands.Operand1, operands.Operand2)
		w.Write([]byte("one of the operands are empty."))
		// return
	}

	op1, err := strconv.ParseFloat(operands.Operand1, 36)
	if err != nil {
		fmt.Println("GET/add: some operand 1 is empty.", err)
		return
	}

	op2, err := strconv.ParseFloat(operands.Operand2, 36)
	if err != nil {
		fmt.Println("GET/add: some operand 2 is empty.", op1, op2)
		return
	}

	// operation := op1 + op2
	val := fmt.Sprintf("%.1f + %.1f = %.2f\n", op1, op2, op1+op2)
	fmt.Printf(val)
	json.NewEncoder(w).Encode(val)
}

func Sub(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		fmt.Println("GET/sub: the request body is empty.")
		w.Write([]byte("request body is empty"))
		return
	}

	var operands Operands
	json.NewDecoder(r.Body).Decode(&operands)

	if operands.isEmpty("GET/sub") {
		fmt.Println(operands.Operand1, operands.Operand2)
		w.Write([]byte("one of the operands are empty."))
		// return
	}

	op1, err := strconv.ParseFloat(operands.Operand1, 36)
	if err != nil {
		fmt.Println("GET/sub: some operand 1 is empty.", err)
		return
	}

	op2, err := strconv.ParseFloat(operands.Operand2, 36)
	if err != nil {
		fmt.Println("GET/sub: some operand 2 is empty.", op1, op2)
		return
	}

	// operation := op1 + op2
	val := fmt.Sprintf("%.1f - %.1f = %.2f\n", op1, op2, op1-op2)
	fmt.Printf(val)
	json.NewEncoder(w).Encode(val)
}

func Multi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		fmt.Println("GET/multi: the request body is empty.")
		w.Write([]byte("request body is empty"))
		return
	}

	var operands Operands
	json.NewDecoder(r.Body).Decode(&operands)

	if operands.isEmpty("GET/multi") {
		fmt.Println(operands.Operand1, operands.Operand2)
		w.Write([]byte("one of the operands are empty."))
		// return
	}

	op1, err := strconv.ParseFloat(operands.Operand1, 36)
	if err != nil {
		fmt.Println("GET/multi: some operand 1 is empty.", err)
		return
	}

	op2, err := strconv.ParseFloat(operands.Operand2, 36)
	if err != nil {
		fmt.Println("GET/multi: some operand 2 is empty.", op1, op2)
		return
	}

	// operation := op1 + op2
	val := fmt.Sprintf("%.1f * %.1f = %.2f\n", op1, op2, op1*op2)
	fmt.Printf(val)
	json.NewEncoder(w).Encode(val)
}

func Div(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Body == nil {
		fmt.Println("GET/div: the request body is empty.")
		w.Write([]byte("request body is empty"))
		return
	}

	var operands Operands
	json.NewDecoder(r.Body).Decode(&operands)

	if operands.isEmpty("GET/div") {
		fmt.Println(operands.Operand1, operands.Operand2)
		w.Write([]byte("one of the operands are empty."))
		// return
	}

	op1, err := strconv.ParseFloat(operands.Operand1, 36)
	if err != nil {
		fmt.Println("GET/div: some operand 1 is empty.", err)
		return
	}

	op2, err := strconv.ParseFloat(operands.Operand2, 36)
	if err != nil {
		fmt.Println("GET/div: some operand 2 is empty.", op1, op2)
		return
	}

	// operation := op1 + op2
	val := fmt.Sprintf("%.1f / %.1f = %.2f\n", op1, op2, op1/op2)
	fmt.Printf(val)
	json.NewEncoder(w).Encode(val)
}

func (op *Operands) isEmpty(source string) bool {

	if op.Operand1 == "" || op.Operand2 == "" {
		fmt.Println("%s: some operand is empty.", source)
		return true
	}

	return false

}
