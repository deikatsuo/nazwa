package dbquery

// InstalmentsPrintedStatus update nama arah
func InstalmentsPrintedStatus(rid int, printed bool) error {
	db := DB
	query := `UPDATE "order_monthly_credit"
	SET printed=$2
	WHERE id=$1`
	_, err := db.Exec(query, rid, printed)

	return err
}
