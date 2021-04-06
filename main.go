package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	ip := os.Args[1]
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("http://ip-api.com/json/%s", ip), nil)
	if err != nil {
		fmt.Println("Malformed request: ", err)
		os.Exit(1)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error returned from API: ", err)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	var success DetermineSuccess
	err = json.Unmarshal(body, &success)
	if err != nil {
		panic(err)
	}
	if success.Status == "fail" {
		var apiErr IpFail
		err := json.Unmarshal(body, &apiErr)
		if err != nil {
			panic(err)
		}
		fmt.Println(apiErr)
		os.Exit(1)
	}
	var ipJson IpResponse
	err = json.Unmarshal(body, &ipJson)
	if err != nil {
		panic(err)
	}
	fmt.Println(ipJson)
}

type IpResponse struct {
	Query       string  `json:"query"`
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	Isp         string  `json:"isp"`
	Org         string  `json:"org"`
	As          string  `json:"as"`
}

type IpFail struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Query   string `json:"query"`
}

type DetermineSuccess struct {
	Status string `json:"status"`
}
