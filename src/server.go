package main

import (
	"fmt"
	"net/http"
)

func server(config *Config) {
	fmt.Printf("Server is running on http://localhost:%s\n", config.API.Port)
	http.ListenAndServe(fmt.Sprintf(":%s", config.API.Port), nil)
}
