package delivery

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/httpclient"
	"store/pkg/order/app"
)

const createDeliveryURL = "/internal/api/v1/order_delivery"

func NewClient(client http.Client, serviceHost string) app.DeliveryClient {
	return &deliveryClient{httpClient: httpclient.NewClient(client, serviceHost)}
}

type deliveryClient struct {
	httpClient httpclient.Client
}

func (d *deliveryClient) ReserveDelivery(userID uuid.UUID, orderID uuid.UUID) (succeeded bool, err error) {
	request := createDeliveryRequest{
		OrderID: orderID.String(),
		UserID:  userID.String(),
	}
	err = d.httpClient.MakeJSONRequest(request, nil, http.MethodPost, createDeliveryURL)
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

type createDeliveryRequest struct {
	OrderID string `json:"orderID"`
	UserID  string `json:"userID"`
}
