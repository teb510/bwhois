package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

func main() {
	var batch BatchRequest
	flag.Var(&batch, "batch", "Pass a comma separated list of IPs to process multiple IPs simultaneously")
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Println("No IP specified, please add an IP an try again")
		os.Exit(1)
	}
	if len(flag.Args()) > 0 && len(batch) > 0 {
		fmt.Println("Cannot combine individual IP requests and batch requests simultaneously\nIf you meant to do only a batch request please try again with IPs being comma separated.\nIf you meant to pass a single IP then rerun the comman without the batch flag")
		os.Exit(1)
	}

	if len(batch) > 0 {
		for _, ip := range batch {
			doRequest(ip)
		}

	} else {
		ip := os.Args[1]
		doRequest(ip)
	}
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

func PrintSuccess(ip IpResponse) {
	fmt.Fprintf(os.Stdout, "\nQueried IP: %s\nCountry of Origin: %s\nISP: %s\nASN: %s\n\n", ip.Query, ip.Country, ip.Isp, ip.As)
}

func PrintError(ipErr IpFail) {
	fmt.Fprintf(os.Stdout, "\nFailed Request!\nStatus: %s\nMessage: %s\nQuery: %s\n\n", ipErr.Status, ipErr.Message, ipErr.Query)
}

type BatchRequest []string

func (br *BatchRequest) Set(s string) error {
	*br = strings.Split(s, ",")
	return nil
}

func (br *BatchRequest) String() string {
	return fmt.Sprintln(*br)
}

func doRequest(ip string) {
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
		PrintError(apiErr)
		os.Exit(1)
	}
	var ipData IpResponse
	err = json.Unmarshal(body, &ipData)
	if err != nil {
		panic(err)
	}
	PrintSuccess(ipData)
}
