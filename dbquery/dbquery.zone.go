package dbquery

import "nazwa/wrapper"

// ZoneShowAll Ambil data zona wilayah dari database
func ZoneShowAll() ([]wrapper.Zone, error) {
	db := DB
	var zones []wrapper.Zone
	var zs []wrapper.ZoneSelect

	query := `SELECT * FROM "zone"`
	err := db.Select(&zs, query)
	if err != nil {
		log.Warn("dbquery.zone.go ZoneShowAll() telah terjadi kesalahan saat memuat data")
		log.Error(err)
		return zones, err
	}

	for _, z := range zs {
		var list []wrapper.ZoneListSelect
		if zl, err := ZoneShowZoneList(z.ID); err == nil {
			list = zl
		} else {
			log.Warn("dbquery.zone.go ZoneShowAll() ambil data list zona")
			log.Error(err)
		}

		zones = append(zones, wrapper.Zone{
			ID:     z.ID,
			CityID: z.CityID,
			Name:   z.Name,
			List:   list,
		})
	}

	return zones, nil
}

// ZoneShowZoneList wilayah
func ZoneShowZoneList(zid int) ([]wrapper.ZoneListSelect, error) {
	db := DB
	var zl []wrapper.ZoneListSelect
	query := `SELECT zl.id, d.name 
	FROM "zone_list" zl
	LEFT JOIN "district" d ON d.id=zl.district_id
	WHERE zl.zone_id=$1`
	err := db.Select(&zl, query, zid)
	if err != nil {
		log.Warn("dbquery.zone.go ZoneShowZoneList() telah terjadi kesalahan saat memuat data")
		log.Error(err)
		return zl, err
	}

	return zl, nil
}
