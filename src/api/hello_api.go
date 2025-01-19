package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tyriis/rest-go/src/service"
)

type HelloApi struct {
	service *service.HelloService
	router  *mux.Router
}

func NewHelloApi(service *service.HelloService) *HelloApi {
	api := &HelloApi{
		service: service,
		router:  mux.NewRouter(),
	}
	api.router.HandleFunc("/", api.HandleHello)
	api.router.HandleFunc("/hello/{name}", api.HandlePersonalizedHello)
	http.Handle("/", api.router)
	return api
}

func (api *HelloApi) HandleHello(w http.ResponseWriter, r *http.Request) {
	message := api.service.GetHelloMessage()
	fmt.Fprintf(w, message)
}

func (api *HelloApi) HandlePersonalizedHello(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	message := api.service.GetPersonalizedHello(name)
	fmt.Fprintf(w, message)
}
