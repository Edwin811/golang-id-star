package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Customer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getCustomers(w http.ResponseWriter, r *http.Request) {
    // Tambahkan ini agar React bisa akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    // 1. Koneksi (Pastikan port MAMP benar 8889)
    dsn := "root:root@tcp(host.docker.internal:8889)/db_idstar"
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Printf("DB Connection Error: %v", err)
        http.Error(w, "Gagal koneksi database", 500)
        return
    }
    defer db.Close()

    // 2. Query
    rows, err := db.Query("SELECT id, name, email FROM customers")
    if err != nil {
        log.Printf("Query Error: %v", err)
        http.Error(w, "Gagal query data", 500)
        return
    }
    defer rows.Close()

    // 3. Parsing Data (Inisialisasi slice kosong agar tidak return 'null')
    customers := []Customer{} 
    for rows.Next() {
        var c Customer
        if err := rows.Scan(&c.ID, &c.Name, &c.Email); err != nil {
            log.Printf("Scan Error: %v", err)
            continue
        }
        customers = append(customers, c)
    }

    // 4. Return JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(customers)
}
func main() {
	http.HandleFunc("/customers", getCustomers)
	fmt.Println("Customer Service jalan di port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}