package main

import (
	"encoding/json"
	"fmt"
	"go_http_test/pkg/middleware"
	"go_http_test/pkg/router"
	"log"
	"net/http"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	City string `json:"city"`
}

func main() {
	router := router.NewRouter()

	data := map[string]interface{}{}
	data["1"] = 5
	data["2"] = "Hello world"
	// mux.HandleFunc("/", )

	router.POST("/home", middleware.New(loggerMiddleware, helloMiddleware).Then(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("This is post")
		w.Write([]byte("This is post"))
	}))

	router.POST("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("This is hello post"))
	})

	router.GET("/home", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("Hello")
		person := Person{Name: "Marat", Age: 20}

		bytes, err := json.Marshal(person)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte("Error"))
			return
		}
		w.Header().Set("Content-Type", "application/json")

		// w.Header().Add("Content-Type", "application/text")
		w.Write(bytes)

	})
	log.Println("Server starting on port: 8000")
	err := http.ListenAndServe(":8000", &router.Mux)
	log.Fatal(err)

}

func loggerMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Start handler")
		next(w, r)
		fmt.Println("End handler")

	})
}

func helloMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello middleware")
		next(w, r)
		fmt.Println("Hello end middleware")

	})
}
