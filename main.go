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
	FuryToken      = "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImQ0ZDlhMGQzLWM4YTItNDY0Yi1hMGE5LWU3MWM2OTA0MjExNiIsInR5cCI6IkpXVCJ9.eyJhZGRpdGlvbmFsX2luZm8iOnsiZW1haWwiOiJ0YWxlcy5sb3Blc0BtZXJjYWRvbGl2cmUuY29tIiwiZnVsbF9uYW1lIjoiVGFsZXMgQmFsdGFyIExvcGVzIERhIFNpbHZhIiwidXNlcm5hbWUiOiJ0YWxvcGVzIn0sImV4cCI6MTY2OTc1NzYwNCwiaWF0IjoxNjY5NzI1MjA0LCJpZGVudGl0eSI6Im1ybjpzZWdpbmY6YWQ6dXNlci90YWxvcGVzIiwiaXNzIjoiZnVyeV90aWdlciIsInN1YiI6InRhbG9wZXMifQ.fDVTJJLPJHsJU0LEXAV7a7QiJx9TWXBSbllB-dF4Vqx1F1KlM-BgCAUH0VwyycavfRlqCrH43yFk_9I3r5LaPWbW6czNofsdwpLjH_G4R-c-b5bJqytkuIxNLViuByNJrWirXCS5gjepQg6kAmuNIWhgC4o4M-CoicJT7MljOFkGz9od4bijryqoMNWC_6yilk8F_sPlJIfP5Ien8TXbEcyhPsxV_H2jWwIaHqaJRrd4S3De64PSFu_3k8MelNwqsuSa0pIR_WZrTT7WA_4_quCImF3nLt4NIGGGKLf-SdFekFqz4HG5xrKNJiEpmFxAyQtmnFkmi6SszdL42J4Txw"
	RequestHeader  = "X-Tiger-Token"
)

var OutputHeader = []string{"arn", "reconciliation_date", "settlement_date", "value_date", "merchant_date", "working_days", "calendar_days", "valid_to_utc"}

// var mappedResponses = make(map[string][]string)

type Response struct {
	ReconciliationDate string `json:"reconciliation_date"`
	SettlementDate     string `json:"settlement_date"`
	ValueDate          string `json:"value_date"`
	MerchantDate       string `json:"merchant_date"`
	WorkingDays        int    `json:"working_days"`
	CalendarDays       int    `json:"calendar_days"`
	ValidToUTC         string `json:"valid_to_utc"`
}

func (r Response) converser(arn string) []string {
	return []string{arn, r.ReconciliationDate, r.SettlementDate, r.ValueDate, r.MerchantDate, fmt.Sprintf("%d", r.WorkingDays), fmt.Sprintf("%d", r.CalendarDays), r.ValidToUTC}
}

func getData(data [][]string, begin int, end int, i int) {

	client := http.Client{}

	var response Response

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
			json.Unmarshal(body, &response)

			mappedResponses[currentArn] = append(mappedResponses[currentArn], response.converser(currentArn))
			defer resp.Body.Close()
		}

		time.Sleep(100 * time.Millisecond)
	}

	outputFile, err := os.Create(fmt.Sprintf("%s%d.csv", OutputFilePath, i))

	check(err)

	defer outputFile.Close()

	w := csv.NewWriter(outputFile)

	err = w.Write(OutputHeader)

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
