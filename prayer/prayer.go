package prayer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/anas639/prayertime/location"
	"github.com/anas639/prayertime/network"
)

type NextPrayer struct {
	TimeLeft time.Duration
	Time     time.Time
	Name     string
}

type PrayerMeta struct {
	Latitude             float64 `json:"latitude"`
	Longitude            float64 `json:"longitude"`
	Timezone             string  `json:"timezone"`
	LatitudeAdjustMethod string  `json:"latitudeAdjustmentMethod"`
	MidnightMode         string  `json:"midnightMode"`
	School               string  `json:"school"`
}

type PrayerTiming struct {
	Fajr       string `json:"Fajr"`
	Sunrise    string `json:"Sunrise"`
	Dohr       string `json:"Dhuhr"`
	Asr        string `json:"Asr"`
	Sunset     string `json:"Sunset"`
	Maghrib    string `json:"Maghrib"`
	Isha       string `json:"Isha"`
	Imsak      string `json:"Imsak"`
	Midnight   string `json:"Midnight"`
	Firstthird string `json:"Firstthird"`
	Lastthird  string `json:"Lastthird"`
}

type PrayerData struct {
	Timings  PrayerTiming `json:"timings"`
	Meta     PrayerMeta   `json:"meta"`
	Location *time.Location
}

type PrayerResponse struct {
	Data PrayerData `json:"data"`
}

func GetPrayerTime(client *network.HttpClient, city *location.LatLng, date string) (*PrayerData, error) {
	query := map[string]string{
		"latitude":  city.Lat,
		"longitude": city.Lon,
	}

	res, err := client.GET(fmt.Sprintf("/timings/%s", date), query)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unsuccessful response %d", res.StatusCode)
	}

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var prayerResponse PrayerResponse
	json.Unmarshal(bytes, &prayerResponse)
	prayerData := prayerResponse.Data
	prayerData.Location, _ = time.LoadLocation(prayerData.Meta.Timezone)

	return &prayerData, nil

}

func (p *PrayerData) GetNextPrayer(fromTime time.Time) (*NextPrayer, error) {
	if p.Location != nil {
		fromTime = fromTime.In(p.Location)
	}
	timings := p.Timings

	fivePrayers := []string{timings.Fajr, timings.Dohr, timings.Asr, timings.Maghrib, timings.Isha}
	fivePrayerNames := []string{"Fajr", "Dohr", "Asr", "Maghrib", "Isha"}
	// check if fromTime is greater than last prayer time
	timeString := strings.Split(fivePrayers[len(fivePrayers)-1], ":")
	hour, _ := strconv.Atoi(timeString[0])
	minute, _ := strconv.Atoi(timeString[1])
	lastPrayerTime := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), hour, minute, 0, 0, fromTime.Location())
	if fromTime.Equal(lastPrayerTime) || fromTime.After(lastPrayerTime) {
		return nil, fmt.Errorf("EOD")
	}

	for i := range fivePrayers {
		prayerTime := strings.Split(fivePrayers[i], ":")
		hour, _ := strconv.Atoi(prayerTime[0])
		minute, _ := strconv.Atoi(prayerTime[1])
		currentPrayer := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), hour, minute, 0, 0, fromTime.Location())
		if currentPrayer.After(fromTime) {
			return &NextPrayer{
				TimeLeft: currentPrayer.Sub(fromTime),
				Time:     currentPrayer,
				Name:     fivePrayerNames[i],
			}, nil

		}
	}

	return nil, fmt.Errorf("No Prayer Found")
}
