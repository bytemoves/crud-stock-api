package main

import (
	
	"fmt"
	"log"
	"net/http"
	"simple-stock-api/router"
)

func main (){
	r := router.Router()
	fmt.Println("Startinf server on port 8080")

	log.Fatal(http.ListenAndServe(":8080",r))

}