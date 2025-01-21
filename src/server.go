package main

import (
	"fmt"
	"net/http"

	"github.com/tyriis/rest-go/src/dto"
)

func server(config *dto.Config) {
	fmt.Printf("Server is running on http://localhost:%s\n", config.API.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.API.Port), nil)
}
