package misc

import (
	"fmt"
	"time"
)

// IsLastMonth tanggal adalah bulan lalu
func IsLastMonth(d string) bool {
	var result bool
	if d == "" {
		return result
	}

	// hapus tanggal/hari
	d = d[3:]

	// bulan sekarang kurangi 1 = bulan kemarin
	dt := time.Now().AddDate(0, -1, 0)
	comp := fmt.Sprintf("%02d-%d", dt.Month(), dt.Year())

	// Cocokan
	if d == comp {
		result = true
	} else {
		result = false
	}

	return result
}
