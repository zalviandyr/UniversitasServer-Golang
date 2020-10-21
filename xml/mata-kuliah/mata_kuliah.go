package mata_kuliah

import (
	"database/sql"
	"encoding/xml"
	"net/http"

	"github.com/gorilla/mux"
)

// MataKuliah struct
type MataKuliah struct {
	DB *sql.DB
}

type response struct {
	XMLName    xml.Name `xml:"response"`
	StatusCode int      `xml:"status_code"`
	Message    string   `xml:"message"`
	MataKuliah struct {
		XMLName        xml.Name `xml:"mata_kuliah_list"`
		MataKuliahList []mataKuliahResponse
	}
}

type mataKuliahResponse struct {
	XMLName      xml.Name `xml:"mata_kuliah"`
	IDMataKuliah string   `xml:"id_mata_kuliah,attr"`
	Nama         string   `xml:"nama"`
}

// GetAllMataKuliah function
func (mK MataKuliah) GetAllMataKuliah(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == "GET" {
		var response response

		sql := "SELECT * FROM mata_kuliah"
		rows, err := mK.DB.Query(sql)

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
				var mataKuliahResponse mataKuliahResponse

				rows.Scan(&mataKuliahResponse.IDMataKuliah, &mataKuliahResponse.Nama)

				response.MataKuliah.MataKuliahList = append(response.MataKuliah.MataKuliahList, mataKuliahResponse)
			}
		}

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// GetMataKuliah function
func (mK MataKuliah) GetMataKuliah(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == "GET" {
		var response response
		var params = mux.Vars(r)

		sql := "SELECT * FROM mata_kuliah WHERE id_mata_kuliah = ?"
		rows, err := mK.DB.Query(sql, params["id"])

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
				var mataKuliahResponse = mataKuliahResponse{}

				rows.Scan(&mataKuliahResponse.IDMataKuliah, &mataKuliahResponse.Nama)

				response.MataKuliah.MataKuliahList = append(response.MataKuliah.MataKuliahList, mataKuliahResponse)
			}
		}

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// InsertMataKuliah function
func (mK MataKuliah) InsertMataKuliah(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == "POST" {
		var response response

		idMataKuliah := r.FormValue("id_mata_kuliah")
		nama := r.FormValue("nama")

		sql := "INSERT INTO mata_kuliah VALUES(?,?)"
		stmt, err := mK.DB.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()

			_, err := stmt.Exec(idMataKuliah, nama)

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

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// UpdateMataKuliah
func (mK MataKuliah) UpdateMataKuliah(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == "PUT" {
		var response response
		var params = mux.Vars(r)

		nama := r.FormValue("nama")

		sql := `UPDATE mata_kuliah SET
		nama = ?
		WHERE id_mata_kuliah = ?`

		stmt, err := mK.DB.Prepare(sql)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.StatusCode = http.StatusInternalServerError
			response.Message = "No connection into database"
		} else {
			defer stmt.Close()

			result, _ := stmt.Exec(nama, params["id"])
			rowAffected, _ := result.RowsAffected()

			if rowAffected == 0 {
				w.WriteHeader(http.StatusBadRequest)
				response.StatusCode = http.StatusBadRequest
				response.Message = "ID Mata Kuliah not found"
			} else {
				w.WriteHeader(http.StatusOK)
				response.StatusCode = http.StatusOK
				response.Message = "Update successful"
			}
		}

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// DeleteMataKuliah
func (mK MataKuliah) DeleteMataKuliah(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

	if r.Method == "DELETE" {
		var response response
		var params = mux.Vars(r)

		sql := "DELETE FROM mata_kuliah WHERE id_mata_kuliah = ?"

		stmt, err := mK.DB.Prepare(sql)

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
				response.Message = "ID Mata Kuliah not found"
			} else {
				w.WriteHeader(http.StatusOK)
				response.StatusCode = http.StatusOK
				response.Message = "Delete successful"
			}
		}

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}
