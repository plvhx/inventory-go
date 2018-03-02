package route

import(
	"fmt"
	"../api"
	"net/http"
	"github.com/gorilla/mux"
)

func RegisterRoute() (router *mux.Router) {
	r := mux.NewRouter()

	// barang
	r.HandleFunc("/api/barang/get", api.FetchBarang)
	r.HandleFunc("/api/barang/create", api.CreateBarang)
	r.HandleFunc("/api/barang/update/{id}", api.UpdateBarang)
	r.HandleFunc("/api/barang/delete/{id}", api.DeleteBarang)

	// jumlah barang
	r.HandleFunc("/api/jumlah_barang/get", api.FetchJumlahBarang)
	r.HandleFunc("/api/jumlah_barang/create", api.CreateJumlahBarang)
	r.HandleFunc("/api/jumlah_barang/update/{id}", api.UpdateJumlahBarang)
	r.HandleFunc("/api/jumlah_barang/delete/{id}", api.DeleteJumlahBarang)

	// barang masuk
	r.HandleFunc("/api/barang_masuk/get", api.FetchBarangMasuk)
	r.HandleFunc("/api/barang_masuk/create", api.CreateBarangMasuk)
	r.HandleFunc("/api/barang_masuk/update/{id}", api.UpdateBarangMasuk)
	r.HandleFunc("/api/barang_masuk/delete/{id}", api.DeleteBarangMasuk)

	// barang keluar
	r.HandleFunc("/api/barang_keluar/get", api.FetchBarangKeluar)
	r.HandleFunc("/api/barang_keluar/create", api.CreateBarangKeluar)
	r.HandleFunc("/api/barang_keluar/update/{id}", api.UpdateBarangKeluar)
	r.HandleFunc("/api/barang_keluar/delete/{id}", api.DeleteBarangKeluar)

	return r
}

func DeployRoute(r *mux.Router, p int) (e error) {
	err := http.ListenAndServe(fmt.Sprintf(":%d", p), r)

	return err
}