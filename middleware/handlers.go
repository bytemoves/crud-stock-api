package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"simple-stock-api/models"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	
)

type response struct{
	ID int64 `json:"id,omitempty"` 	
	Message string `json:"message,omitempty"`
}

func createConnection () *sql.DB{
	err := godotenv.Load(".env")

	if err != nil{
		log.Fatal("error loading .env" )

	}

	db , err := sql.Open("postgres",os.Getenv("POSTGRES_URL"))

	if err != nil{
		panic(err)
	}
	 err = db.Ping()

	 if err != nil{
		panic(err)
	 }

	 fmt.Println("succesfull connected to postgres")

	 return db


}


func CreateStock (w http.ResponseWriter, r *http.Request){
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil{
		log.Fatalf("unable to decode the request body %v",err)
	}

	insertID := insertStock(stock)
	 res := response{
		ID: insertID,
		Message: "stock created succesfully",
	 }
	 json.NewEncoder(w).Encode(res)
	 

}


func GetStock (w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id,err := strconv.Atoi(params["id"])
	if err != nil{
		log.Fatal("unable to convert the string to int. &v",err)
	}
	stock ,err := getStock(int64(id))
	 if err !=nil {
		log.Fatalf("unable to get stock . %v",err)
	 }

	 json.NewEncoder(w).Encode(stock)

}

func GETALLStock(w http.ResponseWriter, r *http.Request){
	stocks , err := getALLStocks()

	if err != nil{
		log.Fatalf("unable to get all the stocks %v",err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock (w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id , err := strconv.Atoi(params["id"])
	if err != nil{
		log.Fatal("unable to convert the string to int. &v",err)
	}
	var stock models.Stock
	 err = json.NewDecoder(r.Body).Decode(&stock)
	 if err != nil{
		log.Fatalf("unable to decode the request body . %v", err)

	 }

	 updateRows := updateStock(int64(id),stock)

	 msg := fmt.Sprintf("stock updated succesfully. Total rows /record %v",updateRows)

	 res := response{
		ID: int64(id),
		Message: msg,
	 }

	 json.NewEncoder(w).Encode(res)
	
}

func Deletestock(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	id , err := strconv.Atoi(params["id"])
	if err != nil{
		log.Fatalf("Unable to convert string to int  %v", err)

	}

	deletedRows  := deleteStock(int64(id))
	msg := fmt.Sprintf("stock deleted succedfully total rows %v",deletedRows)

	res := response{
		ID: int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
	

}


func insertStock(stock  models.Stock)int64 {
	db := createConnection()

	defer db.Close()
	sqlStatement := `INSERT INTO  stocks(name,price,company) VALUES ($1,$2,$3) RETURNING stockid`
	var id int64

	err := db.QueryRow(sqlStatement,stock.Name,stock.Price,stock.Company).Scan(&id)
	if err != nil{
		log.Fatalf("unable to execute the query %v",err)
	}
	fmt.Printf("Inserted a single redord %v",id)
 
	return id

}

func getStock(id int64) (models.Stock, error){
	db := createConnection()
	defer db.Close()

	var stock models.Stock
	sqlStatement := `SELECT * FROM stocks  WHERE  stockid=$1`

	row := db.QueryRow(sqlStatement,id)

	err := row.Scan(&stock.StockID, &stock.Name,&stock.Price,&stock.Company)
	switch err{
	case sql.ErrNoRows:
		fmt.Println("no rows were returned")
		return stock , nil
	case nil:
		return stock ,nil

	default:
		log.Fatalf("unable to scan the row %v",err)
	}
	return stock,err

}

func getALLStocks() ([]models.Stock , error) {
	db := createConnection()
	defer db.Close()

	var stocks []models.Stock
	sqlStatement := `SELECT * FROM stocks`
	rows , err := db.Query(sqlStatement)
	
	if err != nil{
		log.Fatalf("unable to execute the query %v",err)
	}
	defer rows.Close()
	for rows.Next(){
		var stock models.Stock
		err = rows.Scan(&stock.StockID,&stock.Name,&stock.Price,&stock.Company)
		if err != nil{
			log.Fatalf("unable to scan the rows %v",err)

		}

		stocks = append(stocks,stock)
	}

	return stocks ,err
}

func updateStock(id int64, stock models.Stock) int64 {
	db := createConnection()
	defer db.Close()
	sqlStatement := `UPDATE stocks SET name=$2, price$3, company=$4 WHERE stockid=$1`
	res , err := db.Exec(sqlStatement,id,stock.Name,stock.Price,stock.Company)
	
	if err !=nil{
		log.Fatalf("unable to execute the query %v",err)

	}

	rowsAffected,err := res.RowsAffected()

	if err !=nil{
		log.Fatalf("error while checking affected rows %v",err)
	}
	fmt.Printf("TOtal rows affected %v",rowsAffected)
	return rowsAffected


} 

func deleteStock(id int64) int64{

	db := createConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM stocks WHERE  stockid=$1`
	res , err := db.Exec(sqlStatement,id)
	
	if err !=nil{
		log.Fatalf("unable to execute the query %v",err)

	} 
	rowsAffected,err := res.RowsAffected()

	if err !=nil{
		log.Fatalf("error while checking affected rows %v",err)
	}
	fmt.Printf("TOtal rows affected %v",rowsAffected)
	return rowsAffected

	  
}