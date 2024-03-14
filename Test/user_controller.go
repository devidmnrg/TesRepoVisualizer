package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	m "pbp/Modul4/models"
	"strconv"

	"github.com/gorilla/mux"
)

// ini untuk mendapatkan semua user
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM users"
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]

	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name= '" + name[0] + "'"
	}

	if age != nil {
		if name[0] != "" {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " age= '" + age[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	if !rows.Next() {
		w.Header().Set("Content-Type", "application/json")
		var response m.UsersResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	var user m.User
	var users []m.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	if name == "" || age == 0 || address == "" {
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

	_, errQuery := db.Exec("INSERT INTO users (name, age, address) VALUES (?, ?, ?)", name, age, address)
	if errQuery != nil {
		http.Error(w, "Insert failed", http.StatusBadRequest)
		return
	}

	if errQuery == nil {
		sendSuccessResponse(w)
	} else {
		sendErrorResponse(w)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	userID := r.URL.Query().Get("id")

	if userID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	age, _ := strconv.Atoi(r.Form.Get("age"))
	address := r.Form.Get("address")

	if name == "" || age == 0 || address == "" {
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

	_, errQuery := db.Exec("UPDATE users SET name = ?, age = ?, address = ? WHERE id = ?", name, age, address, userID)
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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := r.URL.Query().Get("id")

	if userID == "" {
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

	_, errQuery := db.Exec("DELETE FROM users WHERE id = ?", userID)
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

func sendSuccessResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	var response m.UserResponse
	response.Status = 200
	response.Message = "Success"
	json.NewEncoder(w).Encode(response)
}
func sendErrorResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	var response m.UserResponse
	response.Status = 400
	response.Message = "Failed"
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransactionbyID(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := mux.Vars(r)["id"]
	query := "SELECT t.ID, u.ID, u.name, u.age, u.address, p.ID, p.name, p.price, t.quantity FROM transactions t INNER JOIN users u ON t.UserID = u.ID INNER JOIN products p ON t.ProductID = p.ID WHERE u.ID = ?"

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println(err)
		return
	}

	var transactions []m.Transaction
	for rows.Next() {
		var transaction m.Transaction
		var user m.User
		var product m.Product

		if err := rows.Scan(&transaction.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transaction.UserID = user
			transaction.ProductID = product

			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var response m.TransactionsResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func GetAllUserTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT t.ID, u.ID, u.name, u.age, u.address, p.ID, p.name, p.price, t.quantity FROM transactions t INNER JOIN users u ON t.UserID = u.ID INNER JOIN products p ON t.ProductID = p.ID"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var transactions []m.Transaction
	for rows.Next() {
		var transaction m.Transaction
		var user m.User
		var product m.Product

		if err := rows.Scan(&transaction.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transaction.UserID = user
			transaction.ProductID = product

			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var response m.TransactionsResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var credentials m.Credentials
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "SELECT ID, Name, Age, Address FROM Users WHERE Email=? AND Password=?"

	db := connect()
	defer db.Close()

	row := db.QueryRow(query, credentials.Email, credentials.Password)

	var user m.User

	err = row.Scan(&user.ID, &user.Name, &user.Age, &user.Address)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
