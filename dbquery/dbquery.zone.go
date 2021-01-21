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

		zones = append(zones, wrapper.Zone{
			ID:     z.ID,
			CityID: z.CityID,
			Name:   z.Name,
		})
	}

	return zones, nil
}
