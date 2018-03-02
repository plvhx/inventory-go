package api

import(
	"log"
	"fmt"
	"io"
	"strings"
	//"time"
	"io/ioutil"
	"net/http"
	"../handler/db/sqlite3"
	"encoding/json"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type BarangMasuk struct {
	ID int `json:"id,omitempty"`
	Waktu string `json:"waktu,omitempty"`
	JumlahPemesanan int `json:"jumlah-pemesanan"`
	JumlahDiterima int `json:"jumlah-diterima"`
	HargaBeli string `json:"harga-beli"`
	Total string `json:"total"`
	NomorKwitansi string `json:"nomor-kwitansi"`
	Catatan string `json:"catatan,omitempty"`
	SKU string `json:"sku"`
}

func CreateBarangMasuk(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	var bmsk BarangMasuk
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&bmsk)
	if err != nil && nil != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("INSERT INTO barang_masuk(waktu, jumlah_pemesanan, jumlah_diterima, harga_beli, total, nomor_kwitansi, catatan, sku) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(
		bmsk.Waktu,
		bmsk.JumlahPemesanan,
		bmsk.JumlahDiterima,
		bmsk.HargaBeli,
		bmsk.Total,
		bmsk.NomorKwitansi,
		bmsk.Catatan,
		bmsk.SKU,
	)

	if err != nil {
		log.Fatal(err)
	}

	rf, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}

	buf := map[string]string{
		"id": string(bmsk.ID),
		"rows-affected": string(rf),
		"message": fmt.Sprintf("Data barang masuk dengan ID: %d sudah disimpan.", bmsk.ID),
	}

	json.NewEncoder(w).Encode(buf)
}

func FetchBarangMasuk(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT rowid, waktu, jumlah_pemesanan, jumlah_diterima, harga_beli, total, nomor_kwitansi, catatan, sku FROM barang_masuk")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var res [][]string
	bmsk := BarangMasuk{}

	for rows.Next() {
		err := rows.Scan(
			&bmsk.ID,
			&bmsk.Waktu,
			&bmsk.JumlahPemesanan,
			&bmsk.JumlahDiterima,
			&bmsk.HargaBeli,
			&bmsk.Total,
			&bmsk.NomorKwitansi,
			&bmsk.Catatan,
			&bmsk.SKU,
		)

		if err != nil {
			log.Fatal(err)
		}

		tmp := []string{
			fmt.Sprintf("%d", bmsk.ID),
			bmsk.Waktu,
			fmt.Sprintf("%d", bmsk.JumlahPemesanan),
			fmt.Sprintf("%d", bmsk.JumlahDiterima),
			bmsk.HargaBeli,
			bmsk.Total,
			bmsk.NomorKwitansi,
			bmsk.Catatan,
			bmsk.SKU,
		}

		res = append(res, tmp)
	}

	json.NewEncoder(w).Encode(res)
}

func UpdateBarangMasuk(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(405)
		buf := map[string]string{
			"message": "Method not allowed.",
		}
		json.NewEncoder(w).Encode(buf)
	}

	var bmsk BarangMasuk
	v := mux.Vars(r)
	rbuf, _ := ioutil.ReadAll(r.Body)
	err := json.NewDecoder(strings.NewReader(string(rbuf))).Decode(&bmsk)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	db, err := sqlite3.InitHandler("./db/inventory.db")
	if err != nil {
		log.Fatal(err)
	}

	defer sqlite3.DestroyHandler(db)

	stmt, err := db.Prepare("UPDATE barang_masuk SET waktu = ?, jumlah_pemesanan = ?, jumlah_diterima = ?, harga_beli = ?, total = ?, nomor_kwitansi = ?, catatan = ?, sku = ? WHERE sku = ?")
	if err != nil {
		log.Fatal(err)
	}

	res, err := stmt.Exec(
		bmsk.Waktu,
		int(bmsk.JumlahPemesanan),
		int(bmsk.JumlahDiterima),
		bmsk.HargaBeli,
		bmsk.Total,
		bmsk.NomorKwitansi,
		bmsk.Catatan,
		bmsk.SKU,
		bmsk.SKU,
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
		"message": fmt.Sprintf("Data barang masuk dengan ID: %s sudah diupdate.", v["id"]),
	}

	json.NewEncoder(w).Encode(buf)
}

func DeleteBarangMasuk(w http.ResponseWriter, r *http.Request) {
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

	stmt, err := db.Prepare("DELETE FROM barang_masuk WHERE sku = ?")
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
		"message": fmt.Sprintf("Data barang masuk dengan ID: %s sudah dihapus.", v["id"]),
	}

	json.NewEncoder(w).Encode(buf)
}