package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

type Framer struct {
	Ctx           context.Context
	Message       string
	dbLocker      sync.Mutex
	Db            []Flight
	SearchResults chan []Flight
}
type Flight struct {
	To        string
	From      string
	Price     float32
	Departure time.Time
	Arrival   time.Time
	Captain   string
	Duration  time.Duration
}

func New() *Framer {
	log.SetPrefix("Framer: ")
	log.Println("Started .........")
	parentctx := context.Background()
	ctx, cancel := context.WithTimeout(parentctx, 2*time.Second)
	defer cancel()
	return &Framer{
		Ctx:           ctx,
		Message:       "Flamer spinned",
		SearchResults: make(chan []Flight, 10),
	}
}
func (f *Framer) AddFlight(to, from, captain string, departure, arrival time.Time) error {
	flight, err := NewFlight(to, from, captain, departure, arrival)
	if err != nil {
		return err
	}
	f.dbLocker.Lock()
	f.Db = append(f.Db, *flight)
	f.dbLocker.Unlock()
	log.Println("Flight added succesifuly")
	return nil
}
func NewFlight(to, from, captain string, departure, arrival time.Time) (*Flight, error) {
	if to == "" {
		return nil, fmt.Errorf("destination must no be empty")
	}
	if from == "" {
		return nil, fmt.Errorf("destination must no be empty")
	}
	if captain == "" {
		return nil, fmt.Errorf("every flight must have a captain")
	}
	duration := arrival.Sub(departure)
	return &Flight{
		To:        to,
		From:      from,
		Captain:   captain,
		Departure: departure,
		Arrival:   arrival,
		Duration:  duration,
	}, nil
}
func (f *Framer) AllFlights() {
	for key, flight := range f.Db {
		fmt.Printf("%d. Flight from %s to %s departs at %v and arrives at %v with captain %s it takes %.2f hrs \n", key, flight.From, flight.To, flight.Departure, flight.Arrival, flight.Captain, flight.Duration.Hours())
	}
}
func (f *Framer) Search(to, from string, status bool) {
	log.Printf("Search for %s to %s is %v \n", to, from, status)
	if status {
		time.Sleep(5 * time.Second)
	}

	res := []Flight{}
	wg.Add(1)
	go func() {
		for _, flight := range f.Db {
			if to == flight.To && from == flight.From {
				res = append(res, flight)
				// fmt.Printf("%d. Flight from %s to %s departs at %v and arrives at %v with captain %s it takes %.0f hrs \n", key, flight.From, flight.To, flight.Departure, flight.Arrival, flight.Captain, flight.Duration.Hours())
			}
			f.SearchResults <- res
		}
		wg.Done()
	}()

	wg.Done()
}
func main() {
	framer := New()

	depacture := time.Date(2023, 3, 24, 10, 00, 0, 0, time.Local)
	arrival := time.Date(2023, 3, 24, 18, 15, 0, 0, time.Local)
	framer.AddFlight("Nairobi", "London", "Captain David", depacture, arrival)
	depacture1 := time.Date(2023, 3, 24, 12, 00, 0, 0, time.Local)
	arrival1 := time.Date(2023, 3, 24, 15, 20, 0, 0, time.Local)
	framer.AddFlight("Nairobi", "Dubai", "Captain Larson", depacture1, arrival1)
	depacture2 := time.Date(2023, 3, 24, 12, 00, 0, 0, time.Local)
	arrival2 := time.Date(2023, 3, 25, 06, 45, 0, 0, time.Local)
	framer.AddFlight("London", "NYC", "Captain James", depacture2, arrival2)
	depacture3 := time.Date(2023, 3, 24, 10, 00, 0, 0, time.Local)
	arrival3 := time.Date(2023, 3, 24, 18, 00, 0, 0, time.Local)
	framer.AddFlight("London", "Nairobi", "Captain Ann", depacture3, arrival3)
	fmt.Printf("%s\n", strings.Repeat("-", 50))
	framer.AllFlights()
	wg.Add(1)
	go framer.Search("Nairobi", "London", true)
	fmt.Printf("%s \n", strings.Repeat("-", 50))
	wg.Add(1)
	go framer.Search("London", "NYC", false)
	wg.Wait()
	// fmt.Printf("%d number of goroutines running\n", runtime.NumGoroutine())
	// select {
	// case <-framer.Ctx.Done():
	// 	fmt.Println(framer.Ctx.Err())
	// 	os.Exit(0)
	// }
}
