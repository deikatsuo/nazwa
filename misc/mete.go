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
	// Map untuk menyimpan data hasil merge,
	// Tidak langsung di merge ke m1/m2 untuk menghindari
	// error map concurrent read/write (data race)
	nm := map[string]interface{}{}

	// Lock read
	Mut.RLock()
	defer Mut.RUnlock()

	// Masukan data dari m1 (pertama yang akan di tindih/replace)
	for i, v := range m1 {
		nm[i] = v
	}

	// Masukan data dari m2 (dan mereplace data m1 jika index nya sama)
	for i, v := range m2 {
		for range m1 {
			nm[i] = v
		}
	}

	return nm
}
