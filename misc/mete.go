package misc

import (
	"github.com/imdario/mergo"
)

// Mete ...
// Fungsi untuk menggabungkan template
func Mete(m1, m2 map[string]interface{}) map[string]interface{} {
	// Gabungkan
	mergo.Map(&m1, m2, mergo.WithOverride)
	return m1
}
