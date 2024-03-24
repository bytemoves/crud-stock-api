package router

import (
	"simple-stock-api/middleware"
	"github.com/gorilla/mux"
)


func Router() *mux.Router{
	router := mux.NewRouter()



	router.HandleFunc("/api/stock/{id}",middleware.GetStock).Methods("GET","OPTIONS")
	router.HandleFunc("/api/stock",middleware.GETALLStock).Methods("GET","OPTIONS")
	router.HandleFunc("/api/newstock",middleware.CreateStock).Methods("POST","OPTIONS")
	router.HandleFunc("/api/stock/{id}",middleware.UpdateStock).Methods("put","OPTIONS")
	router.HandleFunc("/api/deletestock/{id}",middleware.Deletestock).Methods("DELETE","OPTIONS")
	return router

}