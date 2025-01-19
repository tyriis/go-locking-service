package main

import (
	"github.com/tyriis/rest-go/src/api"
	"github.com/tyriis/rest-go/src/dao"
	"github.com/tyriis/rest-go/src/service"
)

const (
	PORT = "8080"
)

func main() {
	api.NewHelloApi(service.NewHelloService(dao.NewHelloDAO()))
	server()
}
