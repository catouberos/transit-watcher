package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"example.com/go-hello/internal/crawler"
	"example.com/go-hello/internal/models"
)

const (
	OpTimeLayout = "15:04 -0700"
)

func main() {
	done := make(chan bool)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		done <- true
	}()

	transitDataCrawler := crawler.NewPeriodicCrawler(1*time.Hour, 0*time.Second)

	crawler := crawler.NewPeriodicCrawler(30*time.Second, 100*time.Millisecond)

	go func() {
		for {
			<-done
			transitDataCrawler.Done() <- true
			crawler.Done() <- true
		}
	}()

	transitDataUrls := []string{"https://api.gobus.vn/transit/data/getAllData"}

	transitDataCrawler.SetURLs(transitDataUrls)

	go transitDataCrawler.Start()

	go func() {
		for data := range transitDataCrawler.Result() {
			// unmarshal
			routes := &[]models.Route{}
			err := json.Unmarshal(data, routes)

			if err != nil {
				continue
			}

			urls := []string{}

			for _, route := range *routes {
				opTimes := strings.Split(route.Info.OperationTime, " - ")

				if len(opTimes) == 3 {
					now := time.Now()

					startTime, err := time.Parse(OpTimeLayout, opTimes[0]+" +0700")
					startTime = time.Date(now.Year(), now.Month(), now.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), 0, now.Location())

					if err != nil {
						continue
					}

					if time.Now().Before(startTime) {
						continue
					}

					endTime, err := time.Parse(OpTimeLayout, opTimes[1]+" +0700")
					endTime = time.Date(now.Year(), now.Month(), now.Day(), endTime.Hour(), endTime.Minute(), endTime.Second(), 0, now.Location())

					fmt.Println(route.Id, endTime)

					if err != nil {
						continue
					}

					if endTime.After(time.Now().Add(-2 * time.Hour)) {
						continue
					}

				}

				id := route.Id

				urls = append(urls, fmt.Sprintf("https://multipass-api.golabs.vn/v2/public/busmap/route_bus_gps?regionCode=hcm&routeId=%s&direction=0", id))

				if len(route.Variants) > 1 {
					urls = append(urls, fmt.Sprintf("https://multipass-api.golabs.vn/v2/public/busmap/route_bus_gps?regionCode=hcm&routeId=%s&direction=1", id))
				}
			}

			log.Printf("[INFO] Added %d routes/variants", len(urls))

			crawler.SetURLs(urls)
		}
	}()

	go crawler.Start()

	go func() {
		for event := range crawler.Result() {
			geolocation := &[]models.Geolocation{}
			err := json.Unmarshal(event, geolocation)

			if err != nil {
				continue
			}

			log.Printf("[GEO] %#v\n", geolocation)
		}
	}()

	<-done
}
