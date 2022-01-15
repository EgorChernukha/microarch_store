package billing

import (
	"net/http"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"

	"store/pkg/common/infrastructure/httpclient"
	"store/pkg/order/app"
)

const processPaymentURL = "/internal/api/v1/payment"

func NewClient(client http.Client, serviceHost string) app.BillingClient {
	return &billingClient{httpClient: httpclient.NewClient(client, serviceHost)}
}

type billingClient struct {
	httpClient httpclient.Client
}

func (c *billingClient) ProcessOrderPayment(userID uuid.UUID, price float64) (succeeded bool, err error) {
	request := processPaymentRequest{
		UserID: uuid.UUID(userID).String(),
		Amount: price,
	}
	err = c.httpClient.MakeJSONRequest(request, nil, http.MethodPost, processPaymentURL)
	if err == nil {
		return true, nil
	}
	if e, ok := errors.Cause(err).(*httpclient.HTTPError); ok {
		if e.StatusCode == http.StatusBadRequest {
			return false, nil
		}
	}

	return false, err
}

type processPaymentRequest struct {
	UserID string  `json:"userID"`
	Amount float64 `json:"amount"`
}
