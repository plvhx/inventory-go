package api

import(
	"net/http"
	"log"
	"fmt"
	"io"
	"strings"
	"io/ioutil"
	"encoding/json"
	"../handler/db/sqlite3"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Barang struct {
	SKU string `json:"sku"`
	NamaItem string `json:"nama-item"`
}

func CreateBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",			
		}
		json.NewEncoder(w).Encode(buf)
	}

	var brg Barang
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&brg)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("INSERT INTO barang(sku, nama_item) VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(brg.SKU, brg.NamaItem)
	if err != nil {
		log.Fatal(err)
	}

	rf, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	buf := map[string]string{
		"id": brg.SKU,
		"affected-rows": string(rf),
		"message": fmt.Sprintf("Barang dengan ID: %s sudah disimpan.", brg.SKU),
	}

	json.NewEncoder(w).Encode(buf)
}

func FetchBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	rows, err := db.Query("SELECT * FROM barang")
	if err != nil {
		log.Fatal(err)
	}

	var res [][]string
	brg := Barang{}

	for rows.Next() {
		err := rows.Scan(&brg.SKU, &brg.NamaItem)
		if err != nil {
			log.Fatal(err)
		}

		tmp := []string{
			brg.SKU,
			brg.NamaItem,
		}

		res = append(res, tmp)
	}

	defer rows.Close()

	json.NewEncoder(w).Encode(res)
}

func UpdateBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	var brg Barang
	v := mux.Vars(r)
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&brg)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("UPDATE barang SET nama_item = ? WHERE sku = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(brg.NamaItem, v["id"])
	if err != nil {
		log.Fatal(err)
	}

	af, err := res.RowsAffected()

	buf := map[string]string{
		"id": brg.SKU,
		"rows-affected": string(af),
		"message": fmt.Sprintf("Data barang dengan ID: %s sudah diupdate.", brg.SKU),
	}

	json.NewEncoder(w).Encode(buf)
}

func DeleteBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method != "DELETE" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	v := mux.Vars(r)
	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("DELETE FROM barang WHERE sku = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(v["id"])
	if err != nil {
		log.Fatal(err)
	}

	af, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	buf := map[string]string{
		"id": v["id"],
		"rows-affected": string(af),
		"message": fmt.Sprintf("Data barang dengan ID: %s sudah dihapus.", v["id"]),
	}

	json.NewEncoder(w).Encode(buf)
}