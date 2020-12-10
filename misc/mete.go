package misc

/*
// Mete2 menggabungkan dua map
// digunakan untuk menggabungkan template
func Mete2(m1, m2 map[string]interface{}) map[string]interface{} {
	// Gabungkan
	mergo.Merge(&m1, m2, mergo.WithOverride)
	return m1
}
*/

// Mete menggabungkan dua map secara concurrent
func Mete(m1, m2 map[string]interface{}) map[string]interface{} {
	nm := map[string]interface{}{}

	Mut.RLock()
	defer Mut.RUnlock()

	for i, v := range m1 {
		nm[i] = v
	}

	for i, v := range m2 {
		for range m1 {

			nm[i] = v

		}
	}

	return nm
}
