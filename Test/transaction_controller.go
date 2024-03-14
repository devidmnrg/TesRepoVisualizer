package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	m "pbp/Modul4/models"
	"strconv"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		// send error response
		return
	}

	var transaction m.Transaction
	var transactions []m.Transaction

	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
			// send error response
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	w.Header().Set("Content-Type", "application/json")

	var response m.TransactionsResponse

	response.Status = 200
	response.Message = "Success"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func InsertNewTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("user_id"))
	productID, _ := strconv.Atoi(r.Form.Get("product_id"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	if userID == 0 || productID == 0 || quantity == 0 {
		log.Println("Error: Incomplete data provided")
		http.Error(w, "Bad Request: Incomplete data", http.StatusBadRequest)
		return
	}

	var transactions m.Transaction

	err = db.QueryRow("SELECT ID FROM Products WHERE ID = ?", transactions.ProductID.ID).Scan(&productID)

	if err == sql.ErrNoRows {
		_, err := db.Exec("INSERT INTO Products (ID, Name, Price) VALUES (?, '', 0)", transactions.ProductID.ID)
		if err != nil {
			log.Println(err)
			http.Error(w, "Failed to insert new product", http.StatusInternalServerError)
			return
		}
	}

	data, err := db.Begin()
	if err != nil {
		log.Println("Error database not found:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("INSERT INTO transactions (userid, productid, quantity) VALUES (?, ?, ?)", userID, productID, quantity)
	if errQuery != nil {
		http.Error(w, "Insert success", http.StatusBadRequest)
		log.Println(errQuery.Error())
		return
	}

	if errQuery == nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	transID := r.URL.Query().Get("id")

	if transID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	userID, _ := strconv.Atoi(r.Form.Get("user_id"))
	productID, _ := strconv.Atoi(r.Form.Get("product_id"))
	quantity, _ := strconv.Atoi(r.Form.Get("quantity"))

	if userID == 0 || productID == 0 || quantity == 0 {
		log.Println("Error: Incomplete data provided")
		http.Error(w, "Bad Request: Incomplete data", http.StatusBadRequest)
		return
	}

	data, err := db.Begin()
	if err != nil {
		log.Println("Error database not found:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("UPDATE transactions SET userid = ?, productid = ?, quantity = ? WHERE id = ?", userID, productID, quantity, transID)
	if errQuery != nil {
		http.Error(w, "Update failed", http.StatusBadRequest)
		return
	}

	if errQuery == nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	transID := r.URL.Query().Get("id")

	if transID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	data, err := db.Begin()
	if err != nil {
		log.Println("Error database not found:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer data.Rollback()

	_, errQuery := db.Exec("DELETE FROM transactions WHERE id = ?", transID)
	if errQuery != nil {
		http.Error(w, "Delete failed", http.StatusBadRequest)
		return
	}

	if errQuery == nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}
}
