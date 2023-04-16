package main

import (
	"fmt"
	"net/http"
)

func main() {
	counter := 0

	err := http.ListenAndServe(":8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
		w.Write([]byte(fmt.Sprintf("This server has served %d requests (don't forget the favicon!)", counter)))
	}))

	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}
}
