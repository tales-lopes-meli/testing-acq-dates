package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

const (
	FilePath       = "paths.csv"
	OutputFilePath = "output"
	BlockSize      = 30
	RoutinesAmount = 6
	FuryToken      = "Bearer <Fury_Token>"
	RequestHeader  = "X-Tiger-Token"
)

// var OutputHeader = []string{"arn", "installment", "business_date", "from", "to", "reconciliation_date", "settlement_date", "value_date", "merchant_date", "working_days", "calendar_days", "valid_to_utc"}

// var mappedResponses = make(map[string][]string)

type Schedule struct {
	ReconciliationDate string `json:"reconciliation_date"`
	SettlementDate     string `json:"settlement_date"`
	ValueDate          string `json:"value_date"`
	MerchantDate       string `json:"merchant_date"`
	WorkingDays        int    `json:"working_days"`
	CalendarDays       int    `json:"calendar_days"`
	ValidToUTC         string `json:"valid_to_utc"`
}

type Response struct {
	Installment  int      `json:"installment"`
	BusinessDate string   `json:"business_date"`
	From         string   `json:"from"`
	To           string   `json:"to"`
	Schedule     Schedule `json:"schedule"`
}

func (r Response) converser(arn string) []string {
	return []string{arn, fmt.Sprintf("%d", r.Installment), r.BusinessDate, r.From, r.To, r.Schedule.ReconciliationDate, r.Schedule.SettlementDate, r.Schedule.ValueDate, r.Schedule.MerchantDate, fmt.Sprintf("%d", r.Schedule.WorkingDays), fmt.Sprintf("%d", r.Schedule.CalendarDays), r.Schedule.ValidToUTC}
}

func getData(data [][]string, begin int, end int, i int) {

	client := http.Client{}

	var responses []Response

	mappedResponses := make(map[string][][]string)

	for i := begin; i <= end; i++ {
		fmt.Printf("%d is getting processed\n", i)
		currentURL := data[i][1]
		currentArn := data[i][0]

		// Using net/http for a HTTP GET request
		// resp, err := http.Get(currentURL)
		req, err := http.NewRequest("GET", currentURL, nil)

		check(err)

		req.Header.Set(RequestHeader, FuryToken)
		resp, err := client.Do(req)

		if err == nil {
			body, err := ioutil.ReadAll(resp.Body)

			check(err)

			// Unmarshaling data
			json.Unmarshal(body, &responses)

			for idx, response := range responses {
				mappedResponses[fmt.Sprintf("%s_%d", currentArn, idx+1)] = append(mappedResponses[fmt.Sprintf("%s_%d", currentArn, idx+1)], response.converser(fmt.Sprintf("%s_%d", currentArn, idx+1)))
			}

			defer resp.Body.Close()
		}

		time.Sleep(100 * time.Millisecond)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s%d.csv", OutputFilePath, i))

	check(err)

	defer outputFile.Close()

	w := csv.NewWriter(outputFile)

	check(err)

	fmt.Println("Writing on output.csv")
	for _, value := range mappedResponses {
		for _, subValue := range value {
			err = w.Write(subValue)
			check(err)
		}
	}

	w.Flush()

	check(w.Error())

	fmt.Println("Writing finisehd!")

}

func main() {

	f, err := os.Open(FilePath)

	check(err)

	defer f.Close()

	r := csv.NewReader(f)

	// data: matrix of strings, therefore the n-th row is the n-th record on paths.csv
	// Each row has 2 fields:
	// 0: ARN
	// 1: URL
	// It is possible to use index to get data, so data[2][0] would return second record's ARN.
	data, err := r.ReadAll()

	check(err)

	fmt.Println("Execution started!")

	getData(data, 1, len(data)-1, 1)

	// Writing

	fmt.Println("Execution ended!")
}
