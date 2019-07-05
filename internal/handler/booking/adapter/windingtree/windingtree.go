package windingtree

import (
	"encoding/json"
	"fmt"
	"github.com/giuseppemaniscalco/piratetree/internal/handler/booking/request"
	"github.com/giuseppemaniscalco/piratetree/internal/handler/booking/response"
	provider "github.com/giuseppemaniscalco/piratetree/internal/provider/windingtree"
)

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

type WindingTreeAdapter interface {
	Book(request *request.Request) (*response.Response, error)
}

type windingTree struct {
	provider provider.WindingTreeProvider
}

func NewWindingTree(provider provider.WindingTreeProvider) WindingTreeAdapter {
	return &windingTree{
		provider: provider,
	}
}

func (a *windingTree) Book(aReq *request.Request) (*response.Response, error) {
	req := new(provider.Request)

	pReq, err := parseRequest(aReq)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(pReq, req); err != nil {
		return nil, fmt.Errorf("provider unserialize request err (%v)", err)
	}

	resp, err := a.provider.Book(req)
	if err != nil {
		return nil, err
	}

	aResp, err := parseResponse(resp)
	if err != nil {
		return nil, err
	}

	return aResp, nil
}

func parseRequest(aReq *request.Request) ([]byte, error) {
	//TODO serialize request

	return []byte(test), nil
}

func parseResponse(resp *provider.Response) (*response.Response, error) {
	aResp := new(response.Response)

	aResp.BookingId = resp.Id

	return aResp, nil
}
