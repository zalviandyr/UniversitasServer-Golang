package main

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

var db *sql.DB

type config struct {
	Connection struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Database string `yaml:"database"`
	} `yaml:"connection"`
}

type response struct {
	XMLName    xml.Name `xml:"responses"`
	StatusCode int      `xml:"status_code"`
	Message    string   `xml:"message"`
	Mahasiswa  struct {
		XMLName       xml.Name    `xml:"mahasiswa_list"`
		MahasiswaList []mahasiswa `xml:"mahasiswa"`
	}
}

type mahasiswa struct {
	IDMahasiswa string `xml:"id_mahasiswa,attr"`
	Nama        string `xml:"nama"`
	Alamat      struct {
		XMLName   xml.Name `xml:"alamat"`
		Jalan     string   `xml:"jalan"`
		Kelurahan string   `xml:"kelurahan"`
		Kecamatan string   `xml:"kecamatan"`
		Kabupaten string   `xml:"kabupaten"`
		Provinsi  string   `xml:"provinsi"`
	}
	Fakultas    string `xml:"fakultas"`
	Jurusan     string `xml:"jurusan"`
	NilaiDetail struct {
		XMLName         xml.Name      `xml:"nilai_detail"`
		NilaiDetailList []nilaiDetail `xml:"nilai"`
	}
}

type nilaiDetail struct {
	IDMahasiswa  string  `xml:"id_mahasiswa,attr"`
	IDMataKuliah string  `xml:"id_mata_kuliah,attr"`
	MataKuliah   string  `xml:"mata_kuliah"`
	Nilai        float32 `xml:"nilai"`
	Semester     int8    `xml:"semester"`
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}

func nilaiMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == "GET" {
		var response response
		var params = mux.Vars(r)
		var paramIdMahasiswa = params["idMahasiswa"]
		var paramIdMataKuliah = params["idMataKuliah"]
		var err error

		var sqlMahasiswa string
		var rowsMahasiswa *sql.Rows

		var sqlNilai string
		var rowsNilai *sql.Rows

		// memilih jika ada parameter atau tidak
		if paramIdMahasiswa == "" {
			sqlMahasiswa = "SELECT * FROM mahasiswa"
			rowsMahasiswa, err = db.Query(sqlMahasiswa)
		} else {
			sqlMahasiswa = "SELECT * FROM mahasiswa WHERE id_mahasiswa = ?"
			rowsMahasiswa, err = db.Query(sqlMahasiswa, params["idMahasiswa"])
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer rowsMahasiswa.Close()

			w.WriteHeader(http.StatusOK)
			response.StatusCode = http.StatusOK
			response.Message = "Enjoy your result"

			// mahasiswa
			for rowsMahasiswa.Next() {
				var mahasiswa mahasiswa

				rowsMahasiswa.Scan(
					&mahasiswa.IDMahasiswa, &mahasiswa.Nama, &mahasiswa.Alamat.Jalan,
					&mahasiswa.Alamat.Kelurahan, &mahasiswa.Alamat.Kecamatan, &mahasiswa.Alamat.Kabupaten,
					&mahasiswa.Alamat.Provinsi, &mahasiswa.Fakultas, &mahasiswa.Jurusan,
				)

				// mengecek jika ada param idMataKuliah
				if paramIdMataKuliah == "" {
					sqlNilai =
						`SELECT id_mahasiswa, id_mata_kuliah, nilai, semester
					FROM mahasiswa JOIN nilai USING(id_mahasiswa) WHERE id_mahasiswa = ?`
					rowsNilai, _ = db.Query(sqlNilai, mahasiswa.IDMahasiswa)
				} else {
					sqlNilai =
						`SELECT id_mahasiswa, id_mata_kuliah, nilai, semester
					FROM mahasiswa JOIN nilai USING(id_mahasiswa) WHERE id_mahasiswa = ? AND id_mata_kuliah = ?`
					rowsNilai, _ = db.Query(sqlNilai, mahasiswa.IDMahasiswa, params["idMataKuliah"])
				}

				// nilai
				for rowsNilai.Next() {
					var nilaiDetail nilaiDetail

					rowsNilai.Scan(
						&nilaiDetail.IDMahasiswa, &nilaiDetail.IDMataKuliah,
						&nilaiDetail.Nilai, &nilaiDetail.Semester,
					)

					sqlMataKuliah :=
						`SELECT nama
						FROM mata_kuliah JOIN nilai USING(id_mata_kuliah) WHERE id_mata_kuliah = ?`
					rowsMataKuliah, _ := db.Query(sqlMataKuliah, nilaiDetail.IDMataKuliah)

					// mata kuliah
					if rowsMataKuliah.Next() {
						rowsMataKuliah.Scan(&nilaiDetail.MataKuliah)
					}

					mahasiswa.NilaiDetail.NilaiDetailList = append(mahasiswa.NilaiDetail.NilaiDetailList, nilaiDetail)
				}

				response.Mahasiswa.MahasiswaList = append(response.Mahasiswa.MahasiswaList, mahasiswa)
			}
		}

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

func main() {
	// database
	db = getConnection()
	defer db.Close()

	// init router
	router := mux.NewRouter()

	// router url
	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/nilai-mahasiswa", nilaiMahasiswa).Methods("GET")
	router.HandleFunc("/nilai-mahasiswa/{idMahasiswa}", nilaiMahasiswa).Methods("GET")
	router.HandleFunc("/nilai-mahasiswa/{idMahasiswa}/{idMataKuliah}", nilaiMahasiswa).Methods("GET")

	// start server with port 8080
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getConnection() *sql.DB {
	file, err := ioutil.ReadFile("../config.yml")
	checkErr(err)

	var config config
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
