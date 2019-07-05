package windingtree

import (
	"github.com/giuseppemaniscalco/piratetree/internal/handler/booking/request"
	"github.com/giuseppemaniscalco/piratetree/internal/handler/booking/response"
	provider "github.com/giuseppemaniscalco/piratetree/internal/provider/windingtree"
)

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

func parseResponse(resp *provider.Response) (*response.Response, error) {
	aResp := new(response.Response)

	//TODO parse provider response
	aResp.BookingId = resp.Id

	return aResp, nil
}
