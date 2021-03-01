package wrapper

// InstalmentShowByDate verifikasi tanggal
type InstalmentShowByDate struct {
	Date string `uri:"date" binding:"required,date"`
}

// InstalmentOrderReceipt wrapper list tagihan harian
type InstalmentOrderReceipt struct {
	OrderID   int
	OrderInfo OrderSimple
	Monthly   []OrderMonthlyCredit
}
