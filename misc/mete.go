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

// Mete mete v2
func Mete(m1, m2 map[string]interface{}) map[string]interface{} {

	for i, v := range m2 {
		/*
			nm := map[string]interface{}{
				i: v,
			}
			fmt.Println("NEW: ", nm)
		*/
		for range m1 {

			m1[i] = v

		}
	}
	// fmt.Println("ALL: ", mret)
	return m1
}
