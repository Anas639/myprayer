package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/anas639/prayertime/location"
	"github.com/anas639/prayertime/network"
	"github.com/anas639/prayertime/prayer"
)

const ddMMYYYY = "02-01-2006"

var geoLocalisation = network.NewClient("https://nominatim.openstreetmap.org")
var aladhan = network.NewClient("http://api.aladhan.com/v1")

func main() {
	var city string
	var date string

	flag.StringVar(&city, "city", "", "Your location")
	flag.StringVar(&city, "c", "", "Your location")
	flag.StringVar(&date, "date", "", "Set prayer date")
	flag.StringVar(&date, "d", "", "Set prayer date")

	flag.Parse()

	command := flag.Arg(0)

	if city == "" {
		log.Fatalf("Please provide a city 🏙️\n")
	}
	latLng, err := location.GetLocationFromCity(geoLocalisation, city)
	if err != nil {
		log.Fatalf("GeoLocalisation API failed! 💥, %s\n", err.Error())
	}
	switch command {
	case "next":
		nextPrayer, err := GetNextPrayer(aladhan, latLng, time.Now())
		if err != nil {
			log.Fatalf("💁 An error occurred while calculating the next prayer: %s", err)
		}
		fmt.Printf("%s\n", latLng.DisplayName)
		fmt.Println("************************")
		fmt.Printf("%s until salat Al-%s📿🕌\n", nextPrayer.TimeLeft.Round(time.Second), nextPrayer.Name)
		break
	default:
		timings, err := GetPrayerTimings(date, latLng)

		if err != nil {
			log.Fatalf("💁 An error occurred while getting prayer timings: %s", err)
		}
		fmt.Printf("%s\n", latLng.DisplayName)
		fmt.Println("************************")
		fmt.Printf("Fajr:\t\t%s\n", timings.Fajr)
		fmt.Printf("Dohr:\t\t%s\n", timings.Dohr)
		fmt.Printf("Asr:\t\t%s\n", timings.Asr)
		fmt.Printf("Maghrib:\t%s\n", timings.Maghrib)
		fmt.Printf("Isha:\t\t%s\n", timings.Isha)
	}
}

func GetPrayerTimings(date string, latLng *location.LatLng) (prayer.PrayerTiming, error) {
	var prayerDate string
	if date == "" {
		prayerDate = time.Now().Format(ddMMYYYY)
	} else {
		_, err := time.Parse(ddMMYYYY, date)
		if err != nil {
			return prayer.PrayerTiming{}, fmt.Errorf("Unable to parse the date: %s. Make sure you are using the dd-MM-yyyy format\n", date)
		}
		prayerDate = date
	}

	prayerData, err := prayer.GetPrayerTime(aladhan, latLng, prayerDate)
	if err != nil {
		return prayer.PrayerTiming{}, err
	}

	return prayerData.Timings, nil
}

func GetNextPrayer(client *network.HttpClient, latLng *location.LatLng, fromTime time.Time) (*prayer.NextPrayer, error) {
	prayerDate, err := prayer.GetPrayerTime(client, latLng, fromTime.Format(ddMMYYYY))

	if err != nil {
		return nil, err
	}

	nextPrayer, err := prayerDate.GetNextPrayer(fromTime)

	if err != nil && err.Error() == "EOD" {
		nextDay := time.Now().Add(24 * time.Hour)

		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		return GetNextPrayer(client, latLng, nextDay)
	}

	return nextPrayer, err
}
