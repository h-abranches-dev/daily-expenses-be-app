package handlers

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
)

const (
	ok                  = "200 OK"
	created             = "201 Created"
	noContent           = "204 No Content"
	badRequest          = "404 Bad Request"
	methodNotAllowed    = "405 Method Not Allowed"
	internalServerError = "500 Internal Server Error"
)

func logResponse(reqMethod, reqPath, reqQuery, respStatus string, resp any) error {
	switch reflect.ValueOf(resp).Kind() {
	case reflect.String:
		log.Printf(">>> %s ; %s ; %v ; %s", reqMethod, reqPath, reqQuery, respStatus)
	case reflect.Struct:
		log.Printf(">>> %s ; %s ; %v ; %s", reqMethod, reqPath, reqQuery, respStatus)
	case reflect.Slice:
		log.Printf(">>> %s ; %s ; %v ; %s", reqMethod, reqPath, reqQuery, respStatus)
	case reflect.Pointer:
		log.Printf(">>> %s ; %s ; %v ; %s", reqMethod, reqPath, reqQuery, respStatus)
	default:
		return fmt.Errorf("unhandled default case")
	}
	return nil
}

func logDetailedError(detailedErr error) {
	log.Printf(">>> %s", internalServerError)
	err := fmt.Errorf("error: %s", detailedErr)
	log.Printf(">>> %s", err.Error())
}

func writeResponseWithError(w http.ResponseWriter, respStatusCode int, respStatus string) {
	err := fmt.Errorf("error: %s", respStatus)
	log.Printf(">>> %s", err.Error())
	http.Error(w, err.Error(), respStatusCode)
}

func writeResponseWithDetailedError(w http.ResponseWriter, respStatusCode int, respStatus string, detailedErr error) {
	err := fmt.Errorf("error: %s => %s", respStatus, detailedErr.Error())
	log.Printf(">>> %s", err.Error())
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Error(w, err.Error(), respStatusCode)
}
