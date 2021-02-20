package dbquery

import (
	"nazwa/wrapper"
	"strings"
)

// LineShowAll select semua arah/line
func LineShowAll() ([]wrapper.LocationLine, error) {
	db := DB
	var lines []wrapper.LocationLine
	var ls []wrapper.LocationLineSelect

	query := `SELECT * FROM "zone_line" ORDER BY id`
	err := db.Select(&ls, query)
	if err != nil {
		log.Warn("dbquery.line.go LineShowAll() telah terjadi kesalahan saat memuat data")
		log.Error(err)
		return lines, err
	}

	for _, li := range ls {
		lines = append(lines, wrapper.LocationLine{
			ID:   li.ID,
			Code: li.Code,
			Name: li.Name,
		})
	}

	return lines, nil
}

// LineNew insert/buat arah baru
func LineNew(n wrapper.LocationLineNewForm) error {
	db := DB
	if _, err := db.Exec(`INSERT INTO "zone_line" (code, name) VALUES ($1, $2)`, strings.ToLower(n.Code), strings.ToLower(n.Name)); err != nil {
		log.Warn("dbquery.line.go LineNew() Gagal menambahkan arah")
		log.Error(err)
		return err
	}
	return nil
}

// LineCodeExist cek kode
func LineCodeExist(c string) bool {
	db := DB
	var id int
	query := `SELECT id FROM "zone_line" WHERE code=$1`
	err := db.Get(&id, query, c)

	if err == nil {
		return true
	}

	return false
}

// LineDelete hapus arah
func LineDelete(lid int) error {
	db := DB
	query := `DELETE FROM "zone_line"
	WHERE id=$1`
	_, err := db.Exec(query, lid)

	return err
}
