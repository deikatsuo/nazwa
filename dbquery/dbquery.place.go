package dbquery

import "nazwa/wrapper"

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

// PlaceNewProvince tambah provinsi baru
func PlaceNewProvince(name string, uid int) error {
	db := DB
	if _, err := db.Exec(`INSERT INTO "province" (parent, name, original, created_by) VALUES ($1, $2, $3, $4)`, 62, name, false, uid); err != nil {
		log.Warn("dbquery.zone.go ZoneNew() Gagal menambahkan provinsi")
		log.Error(err)
		return err
	}
	return nil
}
