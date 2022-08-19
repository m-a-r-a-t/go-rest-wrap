package main

import (
	"fmt"
	server "go_http_test/internal"
	_ "go_http_test/internal/controllers/user"
	"log"
	"net/http"
)

type Person struct {
	Name *string `json:"name,omitempty" validate:"required" schema:"name"`
	Age  *int    `json:"age,omitempty" validate:"required" schema:"age"`
	City *string `json:"city,omitempty" validate:"required" schema:"city"`
}

func NewPerson() interface{} {
	return &Person{}
}

func main() {

	// mux.HandleFunc("/", )
	// ! Middlewares: middleware.New(loggerMiddleware, helloMiddleware),
	//!======================================================================
	// r.POST("/home", router.Params{
	// 	StructureForValidateQueryDataCreator: NewPerson},
	// 	func(ctx router.HandlerCtx) (router.MyResponse, error) {

	// 		a := (*ctx.Structure).(*Person)
	// 		fmt.Println("Structure from post", *(a.Name))
	// 		// fmt.Println("Structure from post", a.Age)

	// 		return router.MyResponse{Body: a}, nil
	// 	})

	// r.GET("/hello", router.Params{StructureForValidateQueryDataCreator: NewPerson},
	// 	func(ctx router.HandlerCtx) (router.MyResponse, error) {
	// 		p := (*ctx.Structure).(*Person)
	// 		fmt.Println(*p.Name)
	// 		return router.MyResponse{Body: p}, nil
	// 	})
	//!===================================================
	// r.GET("/home", router.Params{}, func(w http.ResponseWriter, r *http.Request) (router.MyResponse, error) {
	// 	//! реализовать передачу валидированных данных сразу в  контексте
	// 	fmt.Println("Hello")
	// 	person := Person{Name: "Marat", Age: 20}

	// 	return router.MyResponse{Body: person}, nil

	// })

	log.Println("Server starting on port: 8000")
	err := http.ListenAndServe(":8000", &server.R.Mux)
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
