package stock

import (
	"log"
	"net/http"

	"github.com/pkg/errors"

	"store/pkg/common/infrastructure/httpclient"
	"store/pkg/order/app"
)

const reservePositionURL = "/internal/api/v1/position/reserve"

func NewClient(client http.Client, serviceHost string) app.StockClient {
	return &stockClient{httpClient: httpclient.NewClient(client, serviceHost)}
}

type stockClient struct {
	httpClient httpclient.Client
}

func (d *stockClient) ReserveOrderPositions(input app.ReserveOrderPositionInput) (succeeded bool, err error) {
	positions := make([]reservePositionRequestPosition, 0, len(input.Positions))
	for _, positionInput := range input.Positions {
		positions = append(positions, reservePositionRequestPosition{
			PositionID: positionInput.PositionID.String(),
			OrderID:    positionInput.OrderID.String(),
			Count:      positionInput.Count,
		})
	}

	request := reservePositionRequest{
		Positions: positions,
	}
	err = d.httpClient.MakeJSONRequest(request, nil, http.MethodPost, reservePositionURL)
	if err == nil {
		return true, nil
	}
	log.Println(err.Error())
	if e, ok := errors.Cause(err).(*httpclient.HTTPError); ok {
		if e.StatusCode == http.StatusBadRequest {
			return false, nil
		}
	}

	return false, err
}

type reservePositionRequestPosition struct {
	PositionID string `json:"position_id"`
	OrderID    string `json:"order_id"`
	Count      int    `json:"count"`
}

type reservePositionRequest struct {
	Positions []reservePositionRequestPosition `json:"positions"`
}
