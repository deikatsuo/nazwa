package wrapper

import "fmt"

func (a Address) String() string {
	var ps string

	ps = a.One

	if a.Two != "" {
		ps = fmt.Sprintf("%s, %s", ps, a.Two)
	}

	ps = fmt.Sprintf("%s, Ds. %s, Kec. %s %s, Kab. %s, %s", ps, a.VillageName, a.DistrictName, a.Zip, a.CityName, a.ProvinceName)

	return ps
}
