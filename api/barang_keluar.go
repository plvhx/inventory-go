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

type BarangKeluar struct {
	ID int `json:"id,omitempty"`
	Waktu string `json:"waktu"`
	JumlahKeluar int `json:"jumlah-keluar"`
	HargaJual string `json:"harga-jual"`
	Total string `json:"total"`
	Catatan string `json:"catatan,omitempty"`
	SKU string `json:"sku"`
}

func CreateBarangKeluar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	var bkel BarangKeluar
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&bkel)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("INSERT INTO barang_keluar(id, waktu, jumlah_keluar, harga_jual, total, catatan, sku) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(
		0,
		bkel.Waktu,
		bkel.JumlahKeluar,
		bkel.HargaJual,
		bkel.Total,
		bkel.Catatan,
		bkel.SKU,
	)

	if err != nil {
		log.Fatal(err)
	}

	rf, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	buf := map[string]string{
		"id": string(bkel.ID),
		"rows-affected": string(rf),
		"message": fmt.Sprintf("Data jumlah barang dengan ID: %d sudah disimpan.", bkel.ID),
	}

	json.NewEncoder(w).Encode(buf)
}

func FetchBarangKeluar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	db, err := sqlite3.InitHandler("./inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	rows, err := db.Query("SELECT * FROM barang_keluar")
	if err != nil {
		log.Fatal(err)
	}

	var res [][]string
	bkel := BarangKeluar{}

	for rows.Next() {
		err := rows.Scan(
			&bkel.ID,
			&bkel.Waktu,
			&bkel.JumlahKeluar,
			&bkel.HargaJual,
			&bkel.Total,
			&bkel.Catatan,
			&bkel.SKU,
		)

		if err != nil {
			log.Fatal(err)
		}

		tmp := []string{
			string(bkel.ID),
			bkel.Waktu,
			string(bkel.JumlahKeluar),
			bkel.HargaJual,
			bkel.Total,
			bkel.Catatan,
			bkel.SKU,
		}

		res = append(res, tmp)
	}

	defer rows.Close()

	json.NewEncoder(w).Encode(res)
}

func UpdateBarangKeluar(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	var bkel BarangKeluar
	v := mux.Vars(r)
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&bkel)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("UPDATE barang_keluar SET waktu = ?, jumlah_keluar = ?, harga_jual = ?, total = ?, catatan = ?, sku = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(
		bkel.Waktu,
		int(bkel.JumlahKeluar),
		bkel.HargaJual,
		bkel.Total,
		bkel.Catatan,
		bkel.SKU,
	)

	if err != nil {
		log.Fatal(err)
	}

	rf, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	buf := map[string]string{
		"id": v["id"],
		"rows-affected": string(rf),
		"message": fmt.Sprintf("Data barang keluar dengan ID: %d sudah diupdate.", v["id"]),
	}

	json.NewEncoder(w).Encode(buf)
}

func DeleteBarangKeluar(w http.ResponseWriter, r *http.Request) {
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

	stmt, err := db.Prepare("DELETE FROM barang_keluar WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(v["id"])
	if err != nil {
		log.Fatal(err)
	}

	rf, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	buf := map[string]string{
		"id": v["id"],
		"rows-affected": string(rf),
		"message": fmt.Sprintf("Data barang keluar dengan ID: %d sudah dihapus.", v["id"]),
	}

	json.NewEncoder(w).Encode(buf)
}