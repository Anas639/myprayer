package location

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/anas639/prayertime/network"
)

type LatLng struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

func GetLocationFromCity(client *network.HttpClient, city string) (*LatLng, error) {
	query := map[string]string{
		"format": "json",
		"q":      city,
	}

	res, err := client.GET("/search", query)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unsuccessful response %d\n", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var list []*LatLng
	json.Unmarshal(bytes, &list)
	if len(list) == 0 {
		return nil, fmt.Errorf("Location not found: %s\n", city)
	}

	return list[0], nil

}
