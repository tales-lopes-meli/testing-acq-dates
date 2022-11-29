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

var OutputHeader = []string{"arn", "installment_1", "business_date_1", "from_1", "to_1", "reconciliation_date_1", "settlement_date_1", "value_date_1", "merchant_date_1", "working_days_1", "calendar_days_1", "valid_to_utc_1", "installment_2", "business_date_2", "from_2", "to_2", "reconciliation_date_2", "settlement_date_2", "value_date_2", "merchant_date_2", "working_days_2", "calendar_days_2", "valid_to_utc_2"}

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

func converser(arn string, r []Response) []string {
	return []string{arn, fmt.Sprintf("%d", r[0].Installment), r[0].BusinessDate, r[0].From, r[0].To, r[0].Schedule.ReconciliationDate, r[0].Schedule.SettlementDate, r[0].Schedule.ValueDate, r[0].Schedule.MerchantDate, fmt.Sprintf("%d", r[0].Schedule.WorkingDays), fmt.Sprintf("%d", r[0].Schedule.CalendarDays), r[0].Schedule.ValidToUTC, fmt.Sprintf("%d", r[1].Installment), r[1].BusinessDate, r[1].From, r[1].To, r[1].Schedule.ReconciliationDate, r[1].Schedule.SettlementDate, r[1].Schedule.ValueDate, r[1].Schedule.MerchantDate, fmt.Sprintf("%d", r[1].Schedule.WorkingDays), fmt.Sprintf("%d", r[1].Schedule.CalendarDays), r[1].Schedule.ValidToUTC}
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

			mappedResponses[currentArn] = append(mappedResponses[currentArn], converser(currentArn, responses))
		}

		defer resp.Body.Close()

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
