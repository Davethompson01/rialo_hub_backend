package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type JsonResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func RespondWithJson(res http.ResponseWriter, StatusCode int, Success bool, Message string, data interface{}) {

	response := JsonResponse{
		Success: Success,
		Message: Message,
		Data:    data,
	}
	json, err := json.Marshal(response)
	if err != nil {
		log.Println("Failed to marshal JSON:", err)
		http.Error(res, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	res.Header().Add("Content-type", "application/json")
	res.WriteHeader(StatusCode)
	res.Write(json)
}

func RespondWithError(res http.ResponseWriter, code int, Message string) {
	if code > 499 {
		fmt.Println("Return with 500 ")
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}

	RespondWithJson(res, code, false, "Error", ErrorResponse{
		Error: Message,
	})
}

func InternalResponseModel(Success bool, Message string, data interface{}) JsonResponse {
	return JsonResponse{
		Success: Success,
		Message: Message,
		Data:    data,
	}
}
