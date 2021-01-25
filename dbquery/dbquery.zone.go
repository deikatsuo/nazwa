package dbquery

import (
	"errors"
	"fmt"
	"nazwa/wrapper"
)

///////////////////////////////////////// SHOW //

// ZoneShowAll Ambil data zona wilayah dari database
func ZoneShowAll() ([]wrapper.Zone, error) {
	db := DB
	var zones []wrapper.Zone
	var zs []wrapper.ZoneSelect

	query := `SELECT * FROM "zone" ORDER BY id`
	err := db.Select(&zs, query)
	if err != nil {
		log.Warn("dbquery.zone.go ZoneShowAll() telah terjadi kesalahan saat memuat data")
		log.Error(err)
		return zones, err
	}

	for _, z := range zs {
		var list []wrapper.ZoneListsSelect
		if zl, err := ZoneShowZoneList(z.ID); err == nil {
			list = zl
		} else {
			log.Warn("dbquery.zone.go ZoneShowAll() ambil data list zona")
			log.Error(err)
		}

		merge := wrapper.Zone{
			ID:   z.ID,
			Name: z.Name,
			List: list,
		}

		var user wrapper.User
		isNE := false
		if z.CollectorID.Valid {
			if zl, err := UserGetByID(int(z.CollectorID.Int32)); err == nil {
				user = zl
				isNE = true
			} else {
				log.Warn("dbquery.zone.go ZoneShowAll() Gagal mengambil data id kolektor")
				log.Error(err)
			}
		}

		if isNE {
			merge.Collector = wrapper.NameIDCode{
				ID:        user.ID,
				Name:      fmt.Sprintf("%s %s", user.Firstname, user.Lastname),
				Thumbnail: user.Avatar,
				Code:      user.Username,
			}
		}

		zones = append(zones, merge)
	}

	return zones, nil
}

// ZoneShowZoneList wilayah
func ZoneShowZoneList(zid int) ([]wrapper.ZoneListsSelect, error) {
	db := DB
	var zl []wrapper.ZoneListsSelect
	query := `SELECT zl.id, zl.district_id, d.name 
	FROM "zone_list" zl
	LEFT JOIN "district" d ON d.id=zl.district_id
	WHERE zl.zone_id=$1 ORDER BY d.name`
	err := db.Select(&zl, query, zid)
	if err != nil {
		log.Warn("dbquery.zone.go ZoneShowZoneList() telah terjadi kesalahan saat memuat data")
		log.Error(err)
		return zl, err
	}

	return zl, nil
}

////////////////////////////////////// UPDATE //

// ZoneUpdateCollector mengubah kolektor pada zona
func ZoneUpdateCollector(zid, uid int) error {
	db := DB
	query := `UPDATE "zone"
	SET collector_id=$1
	WHERE id=$2`
	_, err := db.Exec(query, uid, zid)

	return err
}

////////////////////////////////////// DELETE //

// ZoneDeleteCollector mengosongkan kolektor pada zona
func ZoneDeleteCollector(zid int) error {
	db := DB
	query := `UPDATE "zone"
	SET collector_id=null
	WHERE id=$1`
	_, err := db.Exec(query, zid)

	return err
}

// ZoneDeleteList menghapus list dari zona
func ZoneDeleteList(zid, lid int) error {
	db := DB
	query := `DELETE FROM "zone_list"
	WHERE id=$1 AND zone_id=$2`
	_, err := db.Exec(query, lid, zid)

	return err
}

//////////////////////////////////////////// ADD //

// ZoneAddList menambahkan list wilayah ke zona
func ZoneAddList(zid int, lists wrapper.ZoneAddListForm) error {
	db := DB
	if len(lists.Lists) > 0 {
		for _, lid := range lists.Lists {
			// Insert id wilayah
			if _, err := db.Exec(`INSERT INTO "zone_list" (zone_id, district_id) VALUES ($1, $2)`, zid, lid); err != nil {
				log.Warn("dbquery.zone.go ZoneAddList() gagal insert id wilayah")
				log.Error(err)
				return err
			}
		}
	} else {
		return errors.New("Tidak ada wilayah untuk dimasukan ke zona")
	}

	return nil
}

/////////////////////////////////////// CHECK //

// ZoneListExistsAndRet check nomor telepon
func ZoneListExistsAndRet(lid int, ret *wrapper.NameIDNameID) bool {
	db := DB
	var list wrapper.NameIDNameID
	query := `SELECT INITCAP(d.name) as name_one, d.id as name_one_id, INITCAP(z.name) as name_two, z.id as name_two_id
	FROM "zone_list" zl
	LEFT JOIN "district" d ON d.id=zl.district_id
	LEFT JOIN "zone" z ON z.id=zl.zone_id
	WHERE zl.district_id=$1`
	err := db.Get(&list, query, lid)

	if err == nil {
		*ret = list
		return true
	}

	return false
}

/////////////////////////////////////// NEW //

// ZoneNew buat zona baru
func ZoneNew(name string, uid int) error {
	db := DB
	if _, err := db.Exec(`INSERT INTO "zone" (name, created_by) VALUES ($1, $2)`, name, uid); err != nil {
		log.Warn("dbquery.zone.go ZoneNew() Gagal menambahkan zona")
		log.Error(err)
		return err
	}
	return nil
}
