package user

import (
	"fmt"
	routes "go_http_test/internal"
	"go_http_test/pkg/router"
)

type Person struct {
	Name *string `json:"name,omitempty" validate:"required" schema:"name"`
	Age  *int    `json:"age,omitempty" validate:"required" schema:"age"`
	City *string `json:"city,omitempty" validate:"required" schema:"city"`
}

func NewPerson() interface{} {
	return &Person{}
}

func init() {
	u := routes.R.SubRouter("/user")
	fmt.Println(u.BasePath)
	
	u.GET("/get_user", router.Params{StructureForValidateQueryDataCreator: NewPerson},
		func(ctx router.HandlerCtx) (router.Result, router.Error) {
			p := (*ctx.Structure).(*Person)
			fmt.Println("Name", *p.Name)
			return router.Result{Body: p}, nil
		})

	u.POST("/home", router.Params{
		StructureForValidateQueryDataCreator: NewPerson},
		func(ctx router.HandlerCtx) (router.Result, router.Error) {

			a := (*ctx.Structure).(*Person)
			fmt.Println("Structure from post", *(a.Name))
			// fmt.Println("Structure from post", a.Age)

			return router.Result{Body: a}, nil
		})

}
