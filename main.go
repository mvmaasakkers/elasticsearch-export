package main

import (
	"errors"
	"flag"
	// "fmt"
	"log"
	"os"
)

var sourceHostname = flag.String("sh", "localhost", "Hostname for source")
var sourcePort = flag.String("sp", "9200", "Port for source")
var sourceIndex = flag.String("si", "", "Index for source")

var destinationHostname = flag.String("dh", "localhost", "Hostname for destination")
var destinationPort = flag.String("dp", "9200", "Port for destination")
var destinationIndex = flag.String("di", "", "Index for destination")

var bulkAmount = flag.Int64("ba", 50, "Bulk amount")

var (
	errorSourceIndexEmpty      = errors.New("Source index cannot be empty")
	errorDestinationIndexEmpty = errors.New("Destination index cannot be empty")
)

func init() {
	flag.Parse()

	if *sourceIndex == "" {
		log.Println(errorSourceIndexEmpty)
		os.Exit(1)
	}

	if *destinationIndex == "" {
		log.Println(errorDestinationIndexEmpty)
		os.Exit(1)
	}
}

func main() {
	log.Println("Elasticsearch Import/Export tool starting")

	// First, check index
	CheckIndex()

	// Second, check mapping (and update)
	em, errMapping := GetMapping()
	if errMapping != nil {
		log.Println("Error while getting mapping:", errMapping)
		os.Exit(1)
	}
	for _, index := range em {
		for key, mapping := range index.Mappings {
			PutMapping(key, mapping)
		}
	}

	// Do a query with amount = 1 to get total amount of docs to export
	statsHit, err := Get(1, 0)
	if err != nil {
		log.Println("Error while fetching data:", err)
		os.Exit(1)
	}
	log.Printf("Index %s has %d documents", *sourceIndex, statsHit.Hits.Total)

	var x int64

	// Loop through structure
	for x = 0; x <= statsHit.Hits.Total; x = x + *bulkAmount {
		log.Printf("Putting documents. %d%% Done (%d/%d)\n", int((float64(x)/float64(statsHit.Hits.Total))*100), x, statsHit.Hits.Total)

		results, errGet := Get(*bulkAmount, x)
		if errGet != nil {
			log.Println("Error while fetching data:", err)
			os.Exit(1)
		}
		bulk, errBulk := results.Bulk()
		if errBulk != nil {
			log.Println("Error while creating bulk data:", errBulk)
			os.Exit(1)
		}
		ebp, errPut := Put(bulk)
		if errPut != nil {
			log.Println("Something went wrong posting data:", errPut)
			os.Exit(1)
		}

		for _, item := range ebp.Items {
			if item.Status != 200 && item.Status != 201 {
				log.Println("Something went wrong posting data:", item)
				os.Exit(1)
			}
		}

	}

	log.Printf("Putting documents. %d%% Done (%d/%d)\n", 100, statsHit.Hits.Total, statsHit.Hits.Total)
}
