package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	m "pbp/Modul4/models"
	"strconv"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"

	//Read from Query Param
	name := r.URL.Query()["name"]
	price := r.URL.Query()["price"]
	if name != nil {
		fmt.Println(name[0])
		query += " WERE name = '" + name[0] + "'"
	}

	if price != nil {
		if name[0] != "" {
			query += " AND"
		} else {
			query += " WHERE"
		}
		query += " price='" + price[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		// send error response
		return
	}

	var product m.Product
	var products []m.Product

	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err)
			// send error response
			return
		} else {
			products = append(products, product)
		}
	}
	w.Header().Set("Content-Type", "application/json")

	var response m.ProductsResponse

	response.Status = 200
	response.Message = "Success"
	response.Data = products
	json.NewEncoder(w).Encode(response)
}

func InsertNewProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form data:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
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

	_, errQuery := db.Exec("INSERT INTO products (name, price) VALUES (?, ?)", name, price)
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

func PutProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	prodID := r.URL.Query().Get("id")

	if prodID == "" {
		log.Println("Error: ID missing")
		http.Error(w, "Bad Request: ID missing", http.StatusBadRequest)
		return
	}

	name := r.Form.Get("name")
	price, _ := strconv.Atoi(r.Form.Get("price"))

	if name == "" || price == 0 {
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

	_, errQuery := db.Exec("UPDATE products SET name = ?, price = ? WHERE id = ?", name, price, prodID)
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

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	id := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM transactions WHERE productID = ?", idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stmt, err := db.Prepare("DELETE FROM products WHERE id=?")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(idInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Product deleted successfully!")
}
