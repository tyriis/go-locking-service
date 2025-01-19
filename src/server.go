package main

import (
	"fmt"
	"net/http"
)

func server() {
	fmt.Printf("Server is running on http://localhost:%s\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%s", PORT), nil)
}
