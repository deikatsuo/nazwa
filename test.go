package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

// Buat ngetest csv hhh
func main() {
	openCountries()
}

func openCountries() {
	file, err := os.Open("./setup/sub-districts.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	read := csv.NewReader(file)
	records, err := read.ReadAll()

	fmt.Println("Rekor", records)

}
