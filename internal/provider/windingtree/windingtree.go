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

type Address struct {
	City        string `json:"city"`
	CountryCode string `json:"countryCode"`
	Road        string `json:"road"`
	HouseNumber string `json:"houseNumber"`
}

type Customer struct {
	Name    string   `json:"name"`
	Surname string   `json:"surname"`
	Address *Address `json:"address"`
	Email   string   `json:"email"`
}

type CancellationFee struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Amount uint16 `json:"amount"`
}

type Pricing struct {
	Currency         string             `json:"currency"`
	Total            uint16             `json:"total"`
	CancellationFees []*CancellationFee `json:"cancellationFees"`
}

type Room struct {
	Id           string   `json:"id"`
	GuestInfoIds []string `json:"guestInfoIds"`
}

type GuestInfo struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type Booking struct {
	Arrival   string       `json:"arrival"`
	Departure string       `json:"departure"`
	Rooms     []*Room      `json:"rooms"`
	GuestInfo []*GuestInfo `json:"guestInfo"`
}

type Request struct {
	OriginAddress string    `json:"originAddress"`
	HotelId       string    `json:"hotelId"`
	Customer      *Customer `json:"customer"`
	Pricing       *Pricing  `json:"pricing"`
	Booking       *Booking  `json:"booking"`
}

type Response struct {
	Id     string `json:"id"`
	Status string `json:"status"`
}

type windingTree struct {
	httpClient     HttpClient
	windingTreeUrl string
}

func NewWindingTree(httpClient HttpClient, windingTreeUrl string) WindingTreeProvider {
	return &windingTree{
		httpClient:     httpClient,
		windingTreeUrl: windingTreeUrl,
	}
}

func (p *windingTree) Book(r *Request) (*Response, error) {
	requestBodyBytes, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("provider unserialize request err (%v)", err)
	}
	requestBody := bytes.NewReader(requestBodyBytes)

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
