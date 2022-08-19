package router

import (
	"encoding/json"
	"fmt"
	"go_http_test/pkg/errors"
	"go_http_test/pkg/middleware"
	"io/ioutil"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
)

type Router struct {
	BasePath    string
	Mux         http.ServeMux
	handleFuncs map[string]*Funcs
}

type Funcs struct {
	get  http.HandlerFunc
	post http.HandlerFunc
}

func NewRouter(basePath string) Router {
	return Router{
		BasePath:    basePath,
		Mux:         *http.NewServeMux(),
		handleFuncs: map[string]*Funcs{},
	}
}

// ! сделать глобальный error handler и добавить возможность передавать заголовки
func (r *Router) createHandleFunc(path string) {
	r.handleFuncs[path] = &Funcs{}
	r.Mux.HandleFunc(r.BasePath+path, func(write http.ResponseWriter, request *http.Request) {

		switch true {
		case request.Method == "GET" && r.handleFuncs[path].get != nil:
			r.handleFuncs[path].get(write, request)
		case request.Method == "POST" && r.handleFuncs[path].post != nil:
			r.handleFuncs[path].post(write, request)
		default:
			http.NotFound(write, request)

		}

	})

}

func (r *Router) SubRouter(prefixPath string) *Router {
	subRouter := r
	subRouter.BasePath = r.BasePath + prefixPath
	return subRouter
}

func (r *Router) GET(path string, params Params, fn MyHandleFunc) {

	if _, ok := r.handleFuncs[path]; !ok {
		r.createHandleFunc(path)
	}

	r.handleFuncs[path].get = GlobalHandler(fn, params)

}

func (r *Router) POST(path string, params Params, fn MyHandleFunc) {

	if _, ok := r.handleFuncs[path]; !ok {
		r.createHandleFunc(path)
	}

	r.handleFuncs[path].post = GlobalHandler(fn, params)
}

// type IResponse interface{
// 	nil | MyResponse
// }
type Error (interface{})

type MyHandleFunc func(ctx HandlerCtx) (Result, Error)

type HandlerCtx struct {
	Write     http.ResponseWriter
	Request   *http.Request
	Structure *interface{}
}

type Result struct {
	Body    interface{}
	Status  int32
	Headers map[string]string
}

func GlobalHandler(fn MyHandleFunc,
	params Params,
	//  structure T
) http.HandlerFunc {

	// ! если передана структура ,то валидируем ,если нет структуры ,то пропускаем
	// ! использовать валидатор go-playground/validator
	return params.Middlewares.Then(func(w http.ResponseWriter, r *http.Request) {
		var decoder = schema.NewDecoder()
		var bytes []byte
		var handlerEror Error
		var result Result

		w.Header().Add("Content-Type", "application/json")

		// fmt.Println("Authorization", r.Header.Get("Authorization"))
		switch r.Method {
		case "POST":
			if params.StructureForValidateQueryDataCreator != nil {
				validate := validator.New()

				structure := params.StructureForValidateQueryDataCreator()
				body, err := ioutil.ReadAll(r.Body)

				if err != nil {
					SendError(errors.NewInvalidJSONError(), w)
					return
					// Если не корректный json ,то отправить ошибку о не корректном jsonе
				}

				err = json.Unmarshal(body, &structure)

				if err != nil {
					// Если не корректный json ,то отправить ошибку о не корректном jsonе
					SendError(errors.NewInvalidJSONError(), w)
					return
				}

				err = validate.Struct(structure)
				// ! добавить swagger  swaggo/swag
				if err != nil {
					// Если не прошли json валидацию ,то отправить ошибку валидации
					SendError(errors.NewValidationError("json", err.(validator.ValidationErrors)), w)
					return
				}

				// результат функции и кастомная ошибка
				result, handlerEror = fn(HandlerCtx{Structure: &structure, Write: w, Request: r})

			} else {
				result, handlerEror = fn(HandlerCtx{Write: w, Request: r})
			}

		case "GET":
			if params.StructureForValidateQueryDataCreator != nil {
				validate := validator.New()

				structure := params.StructureForValidateQueryDataCreator()
				err := decoder.Decode(structure, r.URL.Query())

				if err != nil {
					SendError(errors.NewInvalidQueryParametersError(err), w)
					return
				}

				err = validate.Struct(structure)
				// ! добавить swagger  swaggo/swag
				if err != nil {
					SendError(errors.NewValidationError("query", err.(validator.ValidationErrors)), w)
					return
				}
				result, handlerEror = fn(HandlerCtx{Structure: &structure, Write: w, Request: r})
			} else {
				result, handlerEror = fn(HandlerCtx{Write: w, Request: r})
			}

		}

		if handlerEror != nil {
			fmt.Println(handlerEror)
			SendError(handlerEror, w)
			return
		}

		// fmt.Println(result)
		bytes, _ = json.Marshal(result.Body)
		w.Write(bytes)

	})
}

type Params struct {
	StructureForValidateQueryDataCreator func() interface{}
	Middlewares                          middleware.Chain
}

func SendError(e interface{}, w http.ResponseWriter) {
	bytes, _ := json.Marshal(e)
	w.WriteHeader(http.StatusBadRequest)
	w.Write(bytes)
}

/*
1. Для валидации в validate указываем required
2. Для парсинга query в get используем schema:"name" , где в кавычках указываем название query параметра
! type Example struct {
!	Name *string `json:"name,omitempty" validate:"required" schema:"name"`
!	Age  *int    `json:"age,omitempty" validate:"required" schema:"age"`
!	City *string `json:"city,omitempty" validate:"required" schema:"city"`
! }
*/
