package main

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
)

func main() {
	http.HandleFunc("/payments", paymentsHandler)
	http.ListenAndServe(":8080", nil)
}

func paymentsHandler(w http.ResponseWriter, req *http.Request) {

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.TODO()

	buf := new(bytes.Buffer)

	// Include a Validation logic here to sanitize the req.Body when working in a production environment

	buf.ReadFrom(req.Body)

	paymentDetails := buf.String()

	err := redisClient.RPush(ctx, "payments", paymentDetails).Err()

	if err != nil {
		fmt.Fprintf(w, err.Error()+"\r\n")
		return
	}
	fmt.Fprintf(w, "Payment details accepted successfully\r\n")
}
