package api

import(
	"log"
	"fmt"
	"io"
	"strings"
	"io/ioutil"
	"net/http"
	"../handler/db/sqlite3"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type JumlahBarang struct {
	ID int `json:"id,omitempty"`
	Jumlah int `json:"jumlah"`
	SKU string `json:"sku"`
}

func CreateJumlahBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	var jum JumlahBarang
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&jum)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("INSERT INTO jumlah_barang VALUES(?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(jum.Jumlah, jum.SKU)
	if err != nil {
		log.Fatal(err)
	}

	rf, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	buf := map[string]string{
		"id": string(jum.ID),
		"rows-affected": string(rf),
		"message": fmt.Sprintf("Data jumlah barang dengan ID: %d sudah disimpan.", jum.ID),
	}

	json.NewEncoder(w).Encode(buf)
}

func FetchJumlahBarang(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT rowid,jumlah,sku FROM jumlah_barang")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var res [][]string
	jum := JumlahBarang{}

	for rows.Next() {
		err := rows.Scan(&jum.ID, &jum.Jumlah, &jum.SKU)
		if err != nil {
			log.Fatal(err)
		}

		tmp := []string{
			fmt.Sprintf("%d", jum.ID),
			fmt.Sprintf("%d", jum.Jumlah),
			jum.SKU,
		}

		res = append(res, tmp)
	}

	json.NewEncoder(w).Encode(res)
}

func UpdateJumlahBarang(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	var jum JumlahBarang
	v := mux.Vars(r)
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&jum)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("UPDATE jumlah_barang SET jumlah = ?, sku = ? WHERE sku = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(jum.Jumlah, jum.SKU, v["id"])
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
		"message": fmt.Sprintf("Data jumlah barang dengan ID: %s sudah diupdate.", v["id"]),
	}

	json.NewEncoder(w).Encode(buf)
}

func DeleteJumlahBarang(w http.ResponseWriter, r *http.Request) {
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

	stmt, err := db.Prepare("DELETE FROM jumlah_barang WHERE sku = ?")
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
		"message": fmt.Sprintf("Data jumlah barang dengan ID: %s sudah dihapus.", v["id"]),
	}

	json.NewEncoder(w).Encode(buf)
}