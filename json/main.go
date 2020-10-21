package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	m "universitas/json/mahasiswa"
	mK "universitas/json/mata-kuliah"
	n "universitas/json/nilai"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Connection struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"connection"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func main() {
	// database
	db := getConnection()
	defer db.Close()

	// init router
	router := mux.NewRouter()

	// index
	router.HandleFunc("/", index).Methods("GET")

	// init objek mahasiswa
	mahasiswa := m.Mahasiswa{DB: db}

	// route mahasiswa
	router.HandleFunc("/mahasiswa", mahasiswa.GetAllMahasiswa).Methods("GET")
	router.HandleFunc("/mahasiswa/{id}", mahasiswa.GetMahasiswa).Methods("GET")
	router.HandleFunc("/mahasiswa/{id}", mahasiswa.UpdateMahasiswa).Methods("PUT")
	router.HandleFunc("/mahasiswa/{id}", mahasiswa.DeleteMahasiswa).Methods("DELETE")
	router.HandleFunc("/mahasiswa", mahasiswa.InsertMahasiswa).Methods("POST")

	// init object nilai
	nilai := n.Nilai{DB: db}

	// route nilai
	router.HandleFunc("/nilai", nilai.GetAllNilai).Methods("GET")
	router.HandleFunc("/nilai/{id}", nilai.GetNilai).Methods("GET")
	router.HandleFunc("/nilai/{id}", nilai.UpdateNilai).Methods("PUT")
	router.HandleFunc("/nilai/{id}", nilai.DeleteNilai).Methods("DELETE")
	router.HandleFunc("/nilai", nilai.InsertNilai).Methods("POST")

	// init object mata kuliah
	mataKuliah := mK.MataKuliah{DB: db}

	// route mata kuliah
	router.HandleFunc("/mata-kuliah", mataKuliah.GetAllMataKuliah).Methods("GET")
	router.HandleFunc("/mata-kuliah/{id}", mataKuliah.GetMataKuliah).Methods("GET")
	router.HandleFunc("/mata-kuliah/{id}", mataKuliah.UpdateMataKuliah).Methods("PUT")
	router.HandleFunc("/mata-kuliah/{id}", mataKuliah.DeleteMataKuliah).Methods("DELETE")
	router.HandleFunc("/mata-kuliah", mataKuliah.InsertMataKuliah).Methods("POST")

	// start server with port 8080
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getConnection() *sql.DB {
	file, err := ioutil.ReadFile("../config.yml")
	checkErr(err)

	var config Config
	err = yaml.Unmarshal(file, &config)
	checkErr(err)

	dataSource := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		config.Connection.User,
		config.Connection.Password,
		config.Connection.Host,
		config.Connection.Port,
		config.Connection.Database,
	)
	db, err := sql.Open("mysql", dataSource)
	checkErr(err)

	return db
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
