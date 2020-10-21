package nilai

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Nilai struct
type Nilai struct {
	DB *sql.DB
}

type response struct {
	StatusCode int             `json:"status_code"`
	Message    string          `json:"message"`
	Nilai      []nilaiResponse `json:"nilai"`
}

type nilaiResponse struct {
	IDMahasiswa  string  `json:"id_mahasiswa"`
	IDMataKuliah string  `json:"id_mata_kuliah"`
	Nilai        float32 `json:"nilai"`
	Semester     int     `json:"semester"`
}

// GetAllNilai function
func (n Nilai) GetAllNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response response

		sql := "SELECT * FROM nilai"
		rows, err := n.DB.Query(sql)

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
				var nilaiResponse nilaiResponse

				rows.Scan(
					&nilaiResponse.IDMahasiswa, &nilaiResponse.IDMataKuliah,
					&nilaiResponse.Nilai, &nilaiResponse.Semester,
				)

				response.Nilai = append(response.Nilai, nilaiResponse)
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// GetNilai function
func (n Nilai) GetNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var response response
		var params = mux.Vars(r)

		sql := "SELECT * FROM nilai WHERE id_mahasiswa = ?"
		rows, err := n.DB.Query(sql, params["id"])

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
				var nilaiResponse nilaiResponse

				rows.Scan(
					&nilaiResponse.IDMahasiswa, &nilaiResponse.IDMataKuliah,
					&nilaiResponse.Nilai, &nilaiResponse.Semester,
				)

				response.Nilai = append(response.Nilai, nilaiResponse)
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// InsertNilai function
func (n Nilai) InsertNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var response response

		idMahasiswa := r.FormValue("id_mahasiswa")
		idMataKuliah := r.FormValue("id_mata_kuliah")
		nilai := r.FormValue("nilai")
		semester := r.FormValue("semester")

		sql := "INSERT INTO nilai VALUES (?,?,?,?)"
		stmt, err := n.DB.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()

			_, err := stmt.Exec(idMahasiswa, idMataKuliah, nilai, semester)

			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Something wrong"
			} else {
				w.WriteHeader(http.StatusCreated)
				response.StatusCode = http.StatusCreated
				response.Message = "Data created"
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// UpdateNilai
func (n Nilai) UpdateNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "PUT" {
		var response response
		var params = mux.Vars(r)

		idMataKuliah := r.FormValue("id_mata_kuliah")
		nilai := r.FormValue("nilai")
		semester := r.FormValue("semester")

		sql := `UPDATE nilai SET
		nilai = ?,
		semester = ?
		WHERE id_mahasiswa = ? AND id_mata_kuliah = ?`

		stmt, err := n.DB.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()

			if idMataKuliah == "" {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Set id_mata_kuliah in body request"
			} else {
				result, _ := stmt.Exec(nilai, semester, params["id"], idMataKuliah)
				rowAffected, _ := result.RowsAffected()

				if rowAffected == 0 {
					w.WriteHeader(http.StatusBadRequest)
					response.StatusCode = http.StatusBadRequest
					response.Message = "ID Mahasiswa or ID Mata Kuliah not found"
				} else {
					w.WriteHeader(http.StatusOK)
					response.StatusCode = http.StatusOK
					response.Message = "Update successfull"
				}
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// DeleteNilai function
func (n Nilai) DeleteNilai(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "DELETE" {
		var response response
		var params = mux.Vars(r)

		idMataKuliah := r.FormValue("id_mata_kuliah")

		sql := "DELETE FROM nilai WHERE id_mahasiswa = ? AND id_mata_kuliah = ?"
		stmt, err := n.DB.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()

			if idMataKuliah == "" {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "Set id_mata_kuliah in body request"
			} else {
				result, _ := stmt.Exec(params["id"], idMataKuliah)
				rowAffected, _ := result.RowsAffected()

				if rowAffected == 0 {
					w.WriteHeader(http.StatusBadRequest)
					response.StatusCode = http.StatusBadRequest
					response.Message = "ID Mahasiswa or ID Mata Kuliah not found"
				} else {
					w.WriteHeader(http.StatusOK)
					response.StatusCode = http.StatusOK
					response.Message = "Delete successful"
				}
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}
