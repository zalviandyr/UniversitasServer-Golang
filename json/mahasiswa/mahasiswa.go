package mahasiswa

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Mahasiswa struct to set DB
type Mahasiswa struct {
	DB *sql.DB
}

type response struct {
	StatusCode int                 `json:"status_code"`
	Message    string              `json:"message"`
	Mahasiswa  []mahasiswaResponse `json:"mahasiswa"`
}

type mahasiswaResponse struct {
	IDMahasiswa string `json:"id_mahasiswa"`
	Nama        string `json:"nama_mahasiswa"`
	Alamat      struct {
		Jalan     string `json:"jalan"`
		Kelurahan string `json:"kelurahan"`
		Kecamatan string `json:"kecamatan"`
		Kabupaten string `json:"kabupaten"`
		Provinsi  string `json:"provinsi"`
	} `json:"alamat"`
	Fakultas string `json:"fakultas"`
	Jurusan  string `json:"jurusan"`
}

// GetAllMahasiswa function
func (m Mahasiswa) GetAllMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response response

		sql := "SELECT * FROM mahasiswa"
		rows, err := m.DB.Query(sql)

		// error if no connection into database
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer rows.Close()

			w.WriteHeader(http.StatusOK)
			response.StatusCode = http.StatusOK
			response.Message = "Enjoy your result"

			for rows.Next() {
				var mahasiswaResponse mahasiswaResponse

				rows.Scan(
					&mahasiswaResponse.IDMahasiswa, &mahasiswaResponse.Nama,
					&mahasiswaResponse.Alamat.Jalan, &mahasiswaResponse.Alamat.Kelurahan,
					&mahasiswaResponse.Alamat.Kecamatan, &mahasiswaResponse.Alamat.Kabupaten,
					&mahasiswaResponse.Alamat.Provinsi, &mahasiswaResponse.Fakultas,
					&mahasiswaResponse.Jurusan,
				)

				response.Mahasiswa = append(response.Mahasiswa, mahasiswaResponse)
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// GetMahasiswa by id
func (m Mahasiswa) GetMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response response
		var params = mux.Vars(r)

		sql := "SELECT * FROM mahasiswa WHERE id_mahasiswa = ?"
		rows, err := m.DB.Query(sql, params["id"])

		// error if no connection into database
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer rows.Close()

			w.WriteHeader(http.StatusOK)
			response.StatusCode = http.StatusOK
			response.Message = "Enjoy your result"

			for rows.Next() {
				var mahasiswaResponse mahasiswaResponse

				rows.Scan(
					&mahasiswaResponse.IDMahasiswa, &mahasiswaResponse.Nama,
					&mahasiswaResponse.Alamat.Jalan, &mahasiswaResponse.Alamat.Kelurahan,
					&mahasiswaResponse.Alamat.Kecamatan, &mahasiswaResponse.Alamat.Kabupaten,
					&mahasiswaResponse.Alamat.Provinsi, &mahasiswaResponse.Fakultas,
					&mahasiswaResponse.Jurusan,
				)

				response.Mahasiswa = append(response.Mahasiswa, mahasiswaResponse)
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// InsertMahasiswa function:
// to insert mahasiswa
func (m Mahasiswa) InsertMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var response response

		idMahasiswa := r.FormValue("id_mahasiswa")
		nama := r.FormValue("nama")
		jalan := r.FormValue("jalan")
		kelurahan := r.FormValue("kelurahan")
		kecamatan := r.FormValue("kecamatan")
		kabupaten := r.FormValue("kabupaten")
		provinsi := r.FormValue("provinsi")
		fakultas := r.FormValue("fakultas")
		jurusan := r.FormValue("jurusan")

		sql := "INSERT INTO mahasiswa VALUES(?,?,?,?,?,?,?,?,?)"
		stmt, err := m.DB.Prepare(sql)

		// error jika tidak ada koneksi ke database
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()
			_, err := stmt.Exec(
				idMahasiswa, nama, jalan, kelurahan,
				kecamatan, kabupaten, provinsi, fakultas, jurusan,
			)

			// error jika data duplikat
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Data duplicate"
			} else {
				w.WriteHeader(http.StatusCreated)
				response.StatusCode = http.StatusCreated
				response.Message = "Data created"
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// UpdateMahasiswa function
func (m Mahasiswa) UpdateMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "PUT" {
		var response response
		var params = mux.Vars(r)

		nama := r.FormValue("nama")
		jalan := r.FormValue("jalan")
		kelurahan := r.FormValue("kelurahan")
		kecamatan := r.FormValue("kecamatan")
		kabupaten := r.FormValue("kabupaten")
		provinsi := r.FormValue("provinsi")
		fakultas := r.FormValue("fakultas")
		jurusan := r.FormValue("jurusan")

		sqlUpdate := `UPDATE mahasiswa SET
		nama = ?, jalan = ?, kelurahan = ?,
		kecamatan = ?, kabupaten = ?, provinsi = ?,
		fakultas = ?, jurusan = ?
		WHERE id_mahasiswa = ?
		`
		stmt, err := m.DB.Prepare(sqlUpdate)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()
			result, _ := stmt.Exec(
				nama, jalan, kelurahan, kecamatan,
				kabupaten, provinsi, fakultas, jurusan,
				params["id"],
			)
			rowAffected, _ := result.RowsAffected()

			if rowAffected == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "ID mahasiswa not found"
			} else {
				w.WriteHeader(http.StatusOK)
				response.StatusCode = http.StatusOK
				response.Message = "Update successful"
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// DeleteMahasiswa function
func (m Mahasiswa) DeleteMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "DELETE" {
		var response = response{}
		var params = mux.Vars(r)

		sql := "DELETE FROM mahasiswa WHERE id_mahasiswa = ?"
		stmt, err := m.DB.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()

			result, _ := stmt.Exec(params["id"])
			rowAffected, _ := result.RowsAffected()

			if rowAffected == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "ID mahasiswa not found"
			} else {
				w.WriteHeader(http.StatusOK)
				response.StatusCode = http.StatusOK
				response.Message = "Delete successful"
			}
		}

		json.NewEncoder(w).Encode(response)

	}
}
