package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
    "log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Struktur data sesuai tabel database
type Employee struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

var db *sql.DB

func main() {

	// Koneksi ke MAMP MySQL (sesuaikan user:pass dan port)
	var err error
	dsn := "root:root@tcp(host.docker.internal:8889)/db_idstar"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("Error konfigurasi DB:", err)
		return
	}

	// Tambahkan PING untuk memastikan database MAMP beneran nyambung
	err = db.Ping()
	if err != nil {
		fmt.Println("Gagal koneksi ke MySQL MAMP. Pastikan MAMP sudah ON! Error:", err)
		return
	}

	// Routing sederhana
	http.HandleFunc("/employees", getEmployees)
	http.HandleFunc("/employees/add", createEmployee)
	http.HandleFunc("/employees/update", updateEmployee)
	http.HandleFunc("/employees/delete", deleteEmployee)
    http.HandleFunc("/get-customers-from-employee", getRemoteCustomers)


	fmt.Println("Server jalan di http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func getEmployees(w http.ResponseWriter, r *http.Request) {
    // Tambahkan ini agar React bisa akses
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	rows, _ := db.Query("SELECT id, name, position FROM employees")
	defer rows.Close()

	var results []Employee
	for rows.Next() {
		var emp Employee
		rows.Scan(&emp.ID, &emp.Name, &emp.Position)
		results = append(results, emp)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
    // 1. Pastikan hanya menerima method POST
    if r.Method != "POST" {
        http.Error(w, "Method tidak diizinkan", http.StatusMethodNotAllowed)
        return
    }

    var emp Employee
    // 2. Decode JSON dari body request ke struct Employee
    err := json.NewDecoder(r.Body).Decode(&emp)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // 3. Masukkan ke MySQL (MAMP port 8889)
    query := "INSERT INTO employees (name, position) VALUES (?, ?)"
    _, err = db.Exec(query, emp.Name, emp.Position)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // 4. Kirim respon balik
    w.WriteHeader(http.StatusCreated)
    fmt.Fprintf(w, "Data karyawan berhasil ditambah!")
}

// Handler untuk Update Karyawan (PUT)
func updateEmployee(w http.ResponseWriter, r *http.Request) {
    // Validasi method
    if r.Method != "PUT" {
        http.Error(w, "Harus menggunakan PUT", http.StatusMethodNotAllowed)
        return
    }

    var emp Employee
    // Decode data dari body request
    err := json.NewDecoder(r.Body).Decode(&emp)
    if err != nil {
        http.Error(w, "JSON tidak valid", http.StatusBadRequest)
        return
    }

    // Eksekusi Raw SQL Update
    query := "UPDATE employees SET name = ?, position = ? WHERE id = ?"
    _, err = db.Exec(query, emp.Name, emp.Position, emp.ID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Karyawan ID %d berhasil diupdate!", emp.ID)
}
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
    if r.Method != "DELETE" {
        http.Error(w, "Harus menggunakan DELETE", http.StatusMethodNotAllowed)
        return
    }

    // Mengambil ID dari URL: /employees/delete?id=1
    id := r.URL.Query().Get("id")
    if id == "" {
        http.Error(w, "ID diperlukan", http.StatusBadRequest)
        return
    }

    // Eksekusi Raw SQL Delete
    query := "DELETE FROM employees WHERE id = ?"
    _, err := db.Exec(query, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Fprintf(w, "Karyawan ID %s berhasil dihapus!", id)
}

// Contoh fungsi di Employee Service yang memanggil Customer Service
func getInternalCustomerData(w http.ResponseWriter, r *http.Request) {
    // Memanggil service lain menggunakan NAMA SERVICE di docker-compose
    resp, err := http.Get("http://customer-srv:8081/customers")
    if err != nil {
        http.Error(w, "Gagal menghubungi Customer Service", http.StatusInternalServerError)
        return
    }
    defer resp.Body.Close()

    // Teruskan data dari Customer Service ke User
    fmt.Fprint(w, "Data dari Customer Service berhasil diambil lewat jaringan internal Docker!")
}

func getRemoteCustomers(w http.ResponseWriter, r *http.Request) {
    // 1. Panggil service-customer menggunakan NAMA SERVICE di docker-compose
    // Port 8081 adalah port internal yang didefinisikan di docker-compose
    resp, err := http.Get("http://customer-srv:8081/customers")
    if err != nil {
        log.Printf("Gagal memanggil Customer Service: %v", err)
        http.Error(w, "Service Customer sedang tidak tersedia", http.StatusServiceUnavailable)
        return
    }
    defer resp.Body.Close()

    // 2. Baca response body
    var customers []map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&customers); err != nil {
        http.Error(w, "Gagal decode data dari Customer Service", http.StatusInternalServerError)
        return
    }

    // 3. Berikan response balik ke user
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Data pelanggan berhasil diambil oleh Employee Service",
        "data":    customers,
    })
}