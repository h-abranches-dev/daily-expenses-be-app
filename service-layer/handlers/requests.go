package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
)

func debugRequest(r *http.Request) {
	fmt.Printf("REQUEST\n")
	fmt.Printf("Method: %s\n", r.Method)
	fmt.Printf("URL Path: %s\n", r.URL.Path)

	fmt.Printf("Headers\n")
	for k, v := range r.Header {
		fmt.Printf("Key: %s ; Value: %v\n", k, v)
	}

	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
		return
	}
	bodyBytes := buf.String()
	fmt.Printf("Body: %s\n", bodyBytes)
	os.Exit(0)
}
