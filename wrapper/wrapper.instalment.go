package wrapper

// InstalmentShowByDate verifikasi tanggal
type InstalmentShowByDate struct {
	Date   string `uri:"date" binding:"required,date"`
	ZoneID int    `uri:"zid" binding:"required,numeric"`
}

// InstalmentOrderReceipt wrapper list tagihan harian
type InstalmentOrderReceipt struct {
	OrderID      int
	SuggestPrint bool
	// Total kwitansi yang belum selesai yang sudah masuk jatuh tempo
	Undone    int
	OrderInfo OrderSimple
	Monthly   []OrderMonthlyCredit
}

// InstalmentPrintReceipt data kwitansi yang akan di print
type InstalmentPrintReceipt struct {
	ID             int
	Nth            int
	DueDate        string
	Promise        string
	PrintDate      string
	Code           string
	CreditCode     string
	Customer       string
	BillingAddress string
	Deposit        int
	Monthly        int
	Items          string
	Total          int
	Collector      string
}

// InstalmentUpdateReceitPrintStatus status print
type InstalmentUpdateReceiptPrintStatus struct {
	Printed *bool `form:"set" binding:"required"`
}

type InstalmentMoneyIn struct {
	MoneyIn  string `json:"moneyin" binding:"numeric"`
	Receiver int    `json:"receiver" binding:"numeric"`
	Notes    string `json:"notes" binding:"omitempty"`
	Mode     string `json:"mode" binding:"required"`
}
