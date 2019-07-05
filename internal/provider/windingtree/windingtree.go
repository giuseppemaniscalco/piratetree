package windingtree

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	bookingPath = "/booking"
)

type WindingTreeProvider interface {
	Book(*Request) (*Response, error)
}

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Request struct {
	OriginAddress string `json:"originAddress"`
	HotelId       string `json:"hotelId"`
	//TODO wt request property
}

type Response struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type windingTree struct {
	httpClient     HttpClient
	windingTreeUrl string
}

var test = `{
		  "originAddress": "0x0275e1A76B1C3B67575e66074CdF4fD19D43983A",
		  "hotelId": "0xcca04822Ad9c178bdf9da9091218e241f4C28042",
		  "customer": {
			"name": "Sherlock",
			"surname": "Holmes",
			"address": {
			  "city": "London",
			  "countryCode": "GB",
			  "road": "cool street",
			  "houseNumber": "420"
			},
			"email": "sherlock.holmes@houndofthebaskervilles.net"
		  },
		  "pricing": {
			"currency": "RON",
			"total": 24,
			"cancellationFees": [
			  {"from":"2019-07-04","to":"2019-08-08","amount":24}
			]
		  },
		  "booking": {
			"arrival": "2019-08-29",
			"departure": "2019-08-30",
			"rooms": [
			  {
				"id": "single-room-economy",
				"guestInfoIds": ["1"]
			  }
			],
			"guestInfo": [
			  {
				"id": "1",
				"name": "Sherlock",
				"surname": "Holmes"
			  }
			]
		  }
		}`

func NewWindingTree(httpClient HttpClient, windingTreeUrl string) WindingTreeProvider {
	return &windingTree{
		httpClient:     httpClient,
		windingTreeUrl: windingTreeUrl,
	}
}

func (p *windingTree) Book(r *Request) (*Response, error) {
	//TODO replace test with r
	//requestBodyBytes, err := json.Marshal(r)
	//if err != nil {
	//	return nil, fmt.Errorf("provider unserialize request err (%v)", err)
	//}
	requestBody := bytes.NewReader([]byte(test))

	req, err := http.NewRequest(http.MethodPost, p.windingTreeUrl, requestBody)
	if err != nil {
		return nil, fmt.Errorf("provider create request err (%v)", err)
	}
	req.URL.Path = bookingPath
	req.Header.Add("content-type", "application/json")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("provider send request err (%v)", err)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("provider read response err (%v)", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("provider status code response no ok (%s)", responseBody)
	}

	response := new(Response)
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("provider serialize response err (%v)", err)
	}

	return response, nil
}
