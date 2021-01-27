package dbquery

import (
	"errors"
	"fmt"
	"nazwa/misc"
	"nazwa/wrapper"
	"strconv"
)

// PlaceGetProvinces mengambil data provinsi dari database
// Secara default akan mengambil seluruh data provinsi
// @param: 	true | akan mengambil data provinsi original
//			false | data yang ditambahkan manual
func PlaceGetProvinces(conf ...bool) ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	var query string

	if conf == nil {
		query = `SELECT id, INITCAP(name) AS name FROM "province" ORDER BY name`
	} else {
		if conf[0] {
			query = `SELECT id, INITCAP(name) AS name FROM "province" WHERE original=true ORDER BY name`
		} else {
			query = `SELECT id, INITCAP(name) AS name FROM "province" WHERE original=false ORDER BY name`
		}
	}
	err := db.Select(&p, query)
	if err != nil {
		log.Warn("dbquery.place.go PlaceGetProvinces() Kesalahan saat memuat data")
		log.Error(err)
		return []wrapper.Place{}, err
	}
	return p, nil
}

// PlaceGetCities mengambil data kota/kabupaten dari database
func PlaceGetCities(pid int, conf ...bool) ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	var query string

	if conf == nil {
		query = `SELECT id, INITCAP(name) AS name FROM "city" WHERE parent=$1 ORDER BY name`
	} else {
		if conf[0] {
			query = `SELECT id, INITCAP(name) AS name FROM "city" WHERE parent=$1 AND original=true ORDER BY name`
		} else {
			query = `SELECT id, INITCAP(name) AS name FROM "city" WHERE parent=$1 AND original=false ORDER BY name`
		}
	}
	err := db.Select(&p, query, pid)
	if err != nil {
		log.Warn("dbquery.place.go PlaceGetCities() Kesalahan saat memuat data")
		log.Error(err)
		return []wrapper.Place{}, err
	}
	return p, nil
}

// PlaceGetDistricts mengambil data distrik/kecamatan dari database
func PlaceGetDistricts(pid int, conf ...bool) ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	var query string

	if conf == nil {
		query = `SELECT id, INITCAP(name) AS name FROM "district" WHERE parent=$1 ORDER BY name`
	} else {
		if conf[0] {
			query = `SELECT id, INITCAP(name) AS name FROM "district" WHERE parent=$1 AND original=true ORDER BY name`
		} else {
			query = `SELECT id, INITCAP(name) AS name FROM "district" WHERE parent=$1 AND original=false ORDER BY name`
		}
	}
	err := db.Select(&p, query, pid)
	if err != nil {
		log.Warn("dbquery.place.go PlaceGetDistricts() Kesalahan saat memuat data")
		log.Error(err)
		return []wrapper.Place{}, err
	}
	return p, nil
}

// PlaceGetVillages mengambil data kelurahan/desa dari database
func PlaceGetVillages(pid int, conf ...bool) ([]wrapper.Place, error) {
	db := DB
	var p []wrapper.Place
	var query string

	if conf == nil {
		query = `SELECT id, INITCAP(name) AS name FROM "village" WHERE parent=$1 ORDER BY name`
	} else {
		if conf[0] {
			query = `SELECT id, INITCAP(name) AS name FROM "village" WHERE parent=$1 AND original=true ORDER BY name`
		} else {
			query = `SELECT id, INITCAP(name) AS name FROM "village" WHERE parent=$1 AND original=false ORDER BY name`
		}
	}

	err := db.Select(&p, query, pid)
	if err != nil {
		log.Warn("dbquery.place.go PlaceGetVillages() Kesalahan saat memuat data")
		log.Error(err)
		return []wrapper.Place{}, err
	}
	return p, nil
}

////////////////////////////////// ADD //

// PlaceNewProvince tambah provinsi baru
func PlaceNewProvince(countryID int, name string, uid int) error {
	db := DB
	if _, err := db.Exec(`INSERT INTO "province" (parent, name, original, created_by) VALUES ($1, $2, $3, $4)`, countryID, name, false, uid); err != nil {
		log.Warn("dbquery.place.go PlaceNewProvince() Gagal menambahkan provinsi")
		log.Error(err)
		return err
	}
	return nil
}

// PlaceNewCity insert kota ke provinsi
func PlaceNewCity(provinceID int, name string, uid int) error {
	db := DB
	var id int
	if max, err := CityMaxID(provinceID); err == nil {
		// Convert ke string agar dapat di slice
		fid := strconv.Itoa(max)

		if misc.CountDigits(max) > misc.CountDigits(provinceID) {
			fid = fid[misc.CountDigits(provinceID):]

			// Kembalikan ke integer
			nfid, err := strconv.Atoi(fid)
			if err != nil {
				log.Warn("dbquery.place.go PlaceNewCity() gagal konversi string number ke integer")
				log.Error(err)
				return err
			}

			// Harus ditambah satu agar tidak konflik
			id = nfid + 1
		} else {
			log.Warn("dbquery.place.go PlaceNewCity() digit ID kota tidak valid")
			return errors.New("dbquery.place.go PlaceNewCity() digit ID kota tidak valid")
		}
	} else {
		id = +1
	}

	sid := fmt.Sprintf("%d%02d", provinceID, id)

	if _, err := db.Exec(`INSERT INTO "city" (id, parent, name, original, created_by) VALUES ($1, $2, $3, $4, $5)`, sid, provinceID, name, false, uid); err != nil {
		log.Warn("dbquery.place.go PlaceNewCity() Gagal menambahkan kota")
		log.Error(err)
		return err
	}
	return nil
}

// PlaceNewDistrict insert distrik/kecamatan
func PlaceNewDistrict(cityID int, name string, uid int) error {
	db := DB
	var id int
	if max, err := DistrictMaxID(cityID); err == nil {
		// Convert ke string agar dapat di slice
		fid := strconv.Itoa(max)

		if misc.CountDigits(max) > misc.CountDigits(cityID) {
			fid = fid[misc.CountDigits(cityID):]

			// Kembalikan ke integer
			nfid, err := strconv.Atoi(fid)
			if err != nil {
				log.Warn("dbquery.place.go PlaceNewDistrict() gagal konversi string number ke integer")
				log.Error(err)
				return err
			}

			// Harus ditambah satu agar tidak konflik
			id = nfid + 1
		} else {
			log.Warn("dbquery.place.go PlaceNewDistrict() digit ID kota tidak valid")
			return errors.New("dbquery.place.go PlaceNewDistrict() digit ID kota tidak valid")
		}
	} else {
		id = +1
	}

	sid := fmt.Sprintf("%d%03d", cityID, id)

	if _, err := db.Exec(`INSERT INTO "district" (id, parent, name, original, created_by) VALUES ($1, $2, $3, $4, $5)`, sid, cityID, name, false, uid); err != nil {
		log.Warn("dbquery.place.go PlaceNewDistrict() Gagal menambahkan distrik/kecamatan")
		log.Error(err)
		return err
	}
	return nil
}

// PlaceNewVillage insert kelurahan/desa
func PlaceNewVillage(districtID int, name string, uid int) error {
	db := DB

	if _, err := db.Exec(`INSERT INTO "village" (parent, name, original, created_by) VALUES ($1, $2, $3, $4)`, districtID, name, false, uid); err != nil {
		log.Warn("dbquery.place.go PlaceNewDistrict() Gagal menambahkan distrik/kecamatan")
		log.Error(err)
		return err
	}
	return nil
}

///////////////////////////////////// MAXID //

// CityMaxID ID kota terbesar berdasarkan id provinsi
func CityMaxID(provinceID int) (int, error) {
	db := DB

	// Check bila sku sudah ada di database
	var indb int
	query := `SELECT MAX(id) FROM "city" WHERE parent=$1`
	err := db.Get(&indb, query, provinceID)

	return indb, err
}

// DistrictMaxID ID distrik terbesar berdasarkan id kota
func DistrictMaxID(cityID int) (int, error) {
	db := DB

	// Check bila sku sudah ada di database
	var indb int
	query := `SELECT MAX(id) FROM "district" WHERE parent=$1`
	err := db.Get(&indb, query, cityID)

	return indb, err
}

// VillageMaxID ID kelurahan terbesar berdasarkan id distrik
//func VillageMaxID(districtID int) (int, error) {
//	db := DB
//
//	// Check bila sku sudah ada di database
//	var indb int
//	query := `SELECT MAX(id) FROM "village" WHERE parent=$1`
//	err := db.Get(&indb, query, districtID)
//
//	return indb, err
//}
