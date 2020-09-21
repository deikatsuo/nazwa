package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type Data struct {
	Code      string  `csv:"Code"`
	Parent    int     `csv:"Parent"`
	Name      string  `csv:"Name"`
	Latitude  float32 `csv:"Latitude"`
	Longitude float32 `csv:"Longitude"`
	Postal    string  `csv:"Postal"`
}

// Buat ngetest csv hhh
func main() {
	openData("provinces.csv")
}

func openData(d string) {
	file, err := os.Open(fmt.Sprintf("./setup/%s", d))
	//file, err := os.OpenFile("./setup/villages.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	data := []*Data{}

	if err := gocsv.UnmarshalFile(file, &data); err != nil {
		fmt.Println("Telah terjadi error")
		fmt.Println(err)
	}

	for no, da := range data {
		fmt.Println("Nomor Urut", no)
		fmt.Println("Kode", da.Code)
		fmt.Println("Turunan", da.Parent)
		fmt.Println("Nama", da.Name)
		fmt.Println("Latitude", da.Latitude)
		fmt.Println("Longitude", da.Longitude)
		fmt.Println("Kode Pos", da.Postal)
	}

}
