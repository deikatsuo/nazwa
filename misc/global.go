package misc

import (
	"fmt"
	"strings"
	"sync"
)

// Mut mutable shared data
var Mut = &sync.RWMutex{}

// Mutex global
//var Mutex = &sync.Mutex{}

// GenerateSimpleInsertValues membuat insert dinamis
func GenerateSimpleInsertValues(sm map[string]string) string {
	var result string

	if len(sm) > 0 {
		var kk []string
		var kv []string
		for k, v := range sm {
			kk = append(kk, k)
			kv = append(kv, v)
		}

		result = fmt.Sprintf("(%s) VALUES (%s)", strings.Join(kk, ","), strings.Join(kv, ","))
	}

	return result
}
