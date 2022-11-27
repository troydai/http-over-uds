package summary

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

type Series struct {
	name        string
	latencies   []float64
	errors      []error
	statusCodes []int
}

func NewSeries(name string) *Series {
	return &Series{
		name:        name,
		latencies:   make([]float64, 0),
		errors:      make([]error, 0),
		statusCodes: make([]int, 0),
	}
}

// Merge creates a new Series and merge the given ones
func Merge(name string, others ...*Series) *Series {
	merged := &Series{name: name}

	for _, o := range others {
		merged.latencies = append(merged.latencies, o.latencies...)
		merged.errors = append(merged.errors, o.errors...)
		merged.statusCodes = append(merged.statusCodes, o.statusCodes...)
	}

	return merged
}

func (s *Series) Append(resp *http.Response, err error, beginning time.Time) {
	latency := float64(time.Since(beginning).Nanoseconds())
	s.latencies = append(s.latencies, latency)

	if resp != nil {
		s.statusCodes = append(s.statusCodes, resp.StatusCode)
	} else {
		s.statusCodes = append(s.statusCodes, 0)
	}

	s.errors = append(s.errors, err)
}

func (s *Series) Name() string {
	return s.name
}

func (s *Series) Len() int {
	return len(s.latencies)
}

// LatencyDistribution returns the p99, p95, and p50 of the latency
func (s *Series) LatencyDistribution() (p99 float64, p95 float64, p50 float64) {
	sample := make([]float64, 0, len(s.latencies))
	sample = append(sample, s.latencies...)
	sort.Float64s(sample)
	count := len(sample)

	return sample[(count*99)/100] / 1_000_000, sample[(count*95)/100] / 1_000_000, sample[(count*50)/100] / 1_000_000
}

func (s *Series) ErrorRate() float64 {
	if len(s.errors) == 0 {
		return -1
	}

	var errCount float64
	for _, e := range s.errors {
		if e != nil {
			errCount++
		}
	}
	return errCount / float64(len(s.errors))
}

func (s *Series) StatusCodeMap() map[int]float64 {
	r := make(map[int]int)
	for _, sc := range s.statusCodes {
		v, found := r[sc]
		if found {
			r[sc] = v + 1
		} else {
			r[sc] = 1
		}
	}

	rate := make(map[int]float64)
	for code, total := range r {
		rate[code] = float64(total) / float64(s.Len())
	}

	return rate
}

func (s *Series) PresentData() []string {
	p99, p95, p50 := s.LatencyDistribution()

	statusCodeDistribution := ""
	for code, rate := range s.StatusCodeMap() {
		if statusCodeDistribution != "" {
			statusCodeDistribution += ", "
		}
		statusCodeDistribution += fmt.Sprintf("%d=%2.f%%", code, rate*100)
	}

	return []string{
		s.name,
		fmt.Sprintf("%d", s.Len()),
		fmt.Sprintf("%2.f", s.ErrorRate()*100),
		fmt.Sprintf("%.1f", p99),
		fmt.Sprintf("%.1f", p95),
		fmt.Sprintf("%.1f", p50),
		statusCodeDistribution,
	}
}
