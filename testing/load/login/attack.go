package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	vegeta "github.com/tsenart/vegeta/v12/lib"
)

type Credential struct {
	Username string
	Password string
}

func loadCredentials(path string) ([]Credential, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	creds := make([]Credential, 0, len(records))
	for _, rec := range records {
		if len(rec) != 2 {
			continue
		}
		creds = append(creds, Credential{
			Username: rec[0],
			Password: rec[1],
		})
	}

	return creds, nil
}

func main() {

	creds, err := loadCredentials("users.csv")
	if err != nil {
		panic(err)
	}
	if len(creds) == 0 {
		panic(errors.New("no credentials loaded"))
	}

	rate := vegeta.Rate{Freq: 100, Per: time.Second}
	duration := 1 * time.Minute

	targeter := func(t *vegeta.Target) error {
		// Pick random credentials
		c := creds[rand.Intn(len(creds))]

		payload := map[string]string{
			"username": c.Username,
			"password": c.Password,
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		t.Method = "POST"
		t.URL = "http://localhost:8080/login"
		t.Body = body
		t.Header = map[string][]string{
			"Content-Type": {"application/json"},
		}

		return nil
	}

	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Login stress test") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success: %.2f%%\n", metrics.Success*100)
	fmt.Printf("Latency mean: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Latency p95: %s\n", metrics.Latencies.P95)
}
