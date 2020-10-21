package mahasiswa

import (
	"database/sql"
	"encoding/xml"
	"net/http"

	"github.com/gorilla/mux"
)

// Mahasiswa struct to set DB
type Mahasiswa struct {
	DB *sql.DB
}

type response struct {
	XMLName       xml.Name `xml:"response"`
	StatusCode    int      `xml:"status_code"`
	Message       string   `xml:"message"`
	MahasiswaList struct {
		XMLName   xml.Name `xml:"mahasiswa_list"`
		Mahasiswa []mahasiswaResponse
	}
}

type mahasiswaResponse struct {
	XMLName     xml.Name `xml:"mahasiswa"`
	IDMahasiswa string   `xml:"id_mahasiswa,attr"`
	Nama        string   `xml:"nama"`
	Alamat      struct {
		XMLName   xml.Name `xml:"alamat"`
		Jalan     string   `xml:"jalan"`
		Kelurahan string   `xml:"kelurahan"`
		Kecamatan string   `xml:"kecamatan"`
		Kabupaten string   `xml:"kabupaten"`
		Provinsi  string   `xml:"provinsi"`
	}
	Fakultas string `xml:"fakultas"`
	Jurusan  string `xml:"jurusan"`
}

// GetAllMahasiswa function
func (m Mahasiswa) GetAllMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

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

				response.MahasiswaList.Mahasiswa = append(response.MahasiswaList.Mahasiswa, mahasiswaResponse)
			}
		}

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// GetMahasiswa by id
func (m Mahasiswa) GetMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

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

				response.MahasiswaList.Mahasiswa = append(response.MahasiswaList.Mahasiswa, mahasiswaResponse)
			}
		}

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// InsertMahasiswa function:
// to insert mahasiswa
func (m Mahasiswa) InsertMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

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

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// UpdateMahasiswa function
func (m Mahasiswa) UpdateMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

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

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)
	}
}

// DeleteMahasiswa function
func (m Mahasiswa) DeleteMahasiswa(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/xml")

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

		w.Write([]byte(xml.Header))
		encoder := xml.NewEncoder(w)
		encoder.Indent("", "\t")
		encoder.Encode(response)

	}
}
