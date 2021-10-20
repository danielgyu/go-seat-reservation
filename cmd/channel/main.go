package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

var N int = 210

func main() {
	log.Println("started client")

	for i := 0; i < N; i++ {
		res, err := http.Get("http://localhost:8000/channel/sarang")
		if err != nil {
			panic(err)
		}

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)

		fmt.Println(string(body))
	}
}
