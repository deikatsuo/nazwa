package dbquery

import "nazwa/wrapper"

// AddressGetByID ambil alamat berdasarkan ID alamat
func AddressGetByID(aid int) (wrapper.Address, error) {
	db := DB
	var addressesSelect wrapper.AddressSelect
	var addresses wrapper.Address
	query := `SELECT a.id, a.user_id, a.name, a.description, a.one, a.two, a.zip, a.village_id, a.district_id, a.city_id, a.province_id, INITCAP(p.name) AS province_name, INITCAP(c.name) AS city_name, INITCAP(d.name) AS district_name, INITCAP(v.name) AS village_name
	FROM "address" a
	JOIN "province" p ON p.id=a.province_id
	JOIN "city" c ON c.id=a.city_id
	JOIN "district" d ON d.id=a.district_id
	JOIN "village" v ON v.id=a.village_id
	WHERE a.id=$1`
	err := db.Get(&addressesSelect, query, aid)
	if err != nil {
		log.Warn("dbquery.address.go AddressGetByID() Tidak ada alamat ditemukan")
		log.Error(err)
		return addresses, err
	}

	addresses = wrapper.Address{
		ID:           addressesSelect.ID,
		UserID:       addressesSelect.UserID,
		Name:         addressesSelect.Name,
		Description:  addressesSelect.Description.String,
		One:          addressesSelect.One,
		Two:          addressesSelect.Two.String,
		Zip:          addressesSelect.Zip.String,
		Province:     addressesSelect.Province,
		City:         addressesSelect.City,
		District:     addressesSelect.District,
		Village:      addressesSelect.Village,
		ProvinceName: addressesSelect.ProvinceName,
		CityName:     addressesSelect.CityName,
		DistrictName: addressesSelect.DistrictName,
		VillageName:  addressesSelect.VillageName,
	}

	return addresses, err
}

// AddressAdd menambahkan alamat baru
func AddressAdd(address wrapper.AddressInsert) error {
	db := DB
	query := `INSERT INTO "address" (user_id, name, description, one, two, zip, province_id, city_id, district_id, village_id)
	VALUES (:user_id, :name, :description, :one, :two, :zip, :province_id, :city_id, :district_id, :village_id)`
	_, err := db.NamedExec(query, address)
	if err != nil {
		log.Warn("dbquery.address.go AddressAdd() Gagal menambahkan alamat")
		return err
	}
	return nil
}

// AddressUpdate ubah data alamat
func AddressUpdate(address wrapper.AddressInsert) error {
	db := DB
	query := `UPDATE "address"
	SET name=:name, description=:description, one=:one, two=:two, zip=:zip, province_id=:province_id, city_id=:city_id, district_id=:district_id, village_id=:village_id
	WHERE id=:edit`
	_, err := db.NamedExec(query, address)
	if err != nil {
		log.Warn("dbquery.address.go AddressUpdate() Gagal mengubah alamat")
		return err
	}
	return nil
}

// AddressGetByUserID mengambil alamat berdasarkan id user
func AddressGetByUserID(userid int) ([]wrapper.Address, error) {
	db := DB
	var addressesSelect []wrapper.AddressSelect
	var addresses []wrapper.Address
	query := `SELECT a.id, a.user_id, a.name, a.description, a.one, a.two, a.zip, a.village_id, a.district_id, a.city_id, a.province_id, INITCAP(p.name) AS province_name, INITCAP(c.name) AS city_name, INITCAP(d.name) AS district_name, INITCAP(v.name) AS village_name
	FROM "address" a
	JOIN "province" p ON p.id=a.province_id
	JOIN "city" c ON c.id=a.city_id
	JOIN "district" d ON d.id=a.district_id
	JOIN "village" v ON v.id=a.village_id
	WHERE user_id=$1`
	err := db.Select(&addressesSelect, query, userid)
	if err != nil {
		log.Warn("dbquery.address.go AddressGetByUserID() Tidak ada alamat ditemukan")
		return addresses, err
	}

	for _, tmp := range addressesSelect {
		addresses = append(addresses, wrapper.Address{
			ID:           tmp.ID,
			UserID:       tmp.UserID,
			Name:         tmp.Name,
			Description:  tmp.Description.String,
			One:          tmp.One,
			Two:          tmp.Two.String,
			Zip:          tmp.Zip.String,
			Province:     tmp.Province,
			City:         tmp.City,
			District:     tmp.District,
			Village:      tmp.Village,
			ProvinceName: tmp.ProvinceName,
			CityName:     tmp.CityName,
			DistrictName: tmp.DistrictName,
			VillageName:  tmp.VillageName,
		})
	}

	return addresses, err
}
