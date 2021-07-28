package main

import (
	"log"
	"sync"
	"time"

	"gopkg.in/zorkian/go-datadog-api.v2"
)

// Submitter Datadog metric point submitter
type Submitter struct {
	client     *datadog.Client
	metricName string
	mu         sync.Mutex
	buffer     [][2]float64
}

// NewSubmitter creates a new Datadog metric point submitter
func NewSubmitter(apiKey string, metricName string) *Submitter {
	return &Submitter{
		client:     datadog.NewClient(apiKey, ""),
		metricName: metricName,
		buffer:     nil,
	}
}

// Submit collects the data from channel `in` for 1 minute, then submits to Datadog
func (s *Submitter) Submit(in <-chan [2]float64) <-chan struct{} {
	done := make(chan struct{})
	ticker := time.NewTicker(time.Minute)

	go func() {
		for range ticker.C {
			s.submit()
		}
	}()

	go func() {
		for point := range in {
			s.mu.Lock()
			s.buffer = append(s.buffer, point)
			s.mu.Unlock()
		}

		ticker.Stop()
		s.submit()
		close(done)
	}()

	return done
}

func (s *Submitter) submit() {
	var metric datadog.Metric

	metric.SetMetric(s.metricName)
	metric.SetType("gauge")

	s.mu.Lock()
	for _, point := range s.buffer {
		timestamp := point[0]
		value := point[1]
		metric.Points = append(metric.Points, datadog.DataPoint{
			&timestamp,
			&value,
		})
	}

	s.buffer = nil
	s.mu.Unlock()

	err := s.client.PostMetrics([]datadog.Metric{metric})
	if err != nil {
		log.Printf("could not submit metric: %v\n", err)

		// will add back the points into buffer, for the next submission
		s.mu.Lock()
		for _, point := range metric.Points {
			s.buffer = append(s.buffer, [2]float64{*point[0], *point[1]})
		}
		s.mu.Unlock()
	}
}
