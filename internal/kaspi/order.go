package kaspi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/nurtai325/kaspi-service/internal/models"
)

const (
	kaspiDeliveryState = "KASPI_DELIVERY"
)

type ordersResponse struct {
	Data []order `json:"data"`
}

type order struct {
	Id         string `json:"id"`
	Attributes struct {
		Code     string `json:"code"`
		Customer struct {
			CellPhone string `json:"cellPhone"`
			FirstName string `json:"firstName"`
			LastName  string `json:"lastName"`
		} `json:"customer"`
		Status string `json:"status"`
	} `json:"attributes"`
}

func AllOrders(ctx context.Context, token string, from, to time.Time) ([]models.Order, error) {
	var b strings.Builder
	b.WriteString("/orders")
	b.WriteString(fmt.Sprintf("?page[number]=%d&page[size]=%d", 0, 100))
	b.WriteString(fmt.Sprintf("&filter[orders][state]=%s", kaspiDeliveryState))
	b.WriteString(fmt.Sprintf("&filter[orders][creationDate][$ge]=%d&filter[orders][creationDate][$le]=%d", from.UnixMilli(), to.UnixMilli()))
	var parsedRes ordersResponse
	err := do(ctx, token, http.MethodGet, b.String(), &parsedRes)
	if err != nil {
		return nil, fmt.Errorf("error making get all orders request:\n %w", err)
	}
	orders := make([]models.Order, 0, len(parsedRes.Data))
	for _, order := range parsedRes.Data {
		entries, productCode, err := getOrderEntries(ctx, token, order.Id)
		if err != nil {
			return nil, fmt.Errorf("fetching entries for order %s\n %w", order.Id, err)
		}
		orders = append(orders, models.Order{
			Id:          order.Id,
			Code:        order.Attributes.Code,
			ProductCode: productCode,
			Customer:    fmt.Sprintf("%s %s", order.Attributes.Customer.FirstName, order.Attributes.Customer.LastName),
			Entries:     entries,
		})
	}
	return orders, nil
}

type entriesResponse struct {
	Data []struct {
		Id         string `json:"id"`
		Attributes struct {
			Offer struct {
				Code string `json:"code"`
				Name string `json:"Name"`
			} `json:"offer"`
			TotalPrice float32 `json:"totalPrice"`
			Quantity   int     `json:"quantity"`
		} `json:"attributes"`
	} `json:"data"`
}

func getOrderEntries(ctx context.Context, token, orderId string) ([]models.Entry, string, error) {
	link := fmt.Sprintf("/orders/%s/entries", orderId)
	var parsedEntries entriesResponse
	err := do(ctx, token, http.MethodGet, link, &parsedEntries)
	if err != nil {
		return nil, "", err
	}
	entries := make([]models.Entry, 0, len(parsedEntries.Data))
	for _, entry := range parsedEntries.Data {
		entries = append(entries, models.Entry{
			ProductName: entry.Attributes.Offer.Name,
			Quantity:    entry.Attributes.Quantity,
			Price:       entry.Attributes.TotalPrice,
		})
	}
	if len(parsedEntries.Data) < 1 {
		return nil, "", errors.New("parsed entries length is 0")
	}
	productCode := parsedEntries.Data[0].Attributes.Offer.Code
	splitProductCode := strings.Split(productCode, "_")
	if len(splitProductCode) != 2 {
		return nil, "", fmt.Errorf("invalid product code: %s", productCode)
	}
	productCode = splitProductCode[0]
	return entries, productCode, nil
}

type OrderStatus string

const (
	Completed             OrderStatus = "COMPLETED"
	CancelledStatus       OrderStatus = "CANCELLED"
	CancellingStatus      OrderStatus = "CANCELLING"
	ReturnRequestedStatus OrderStatus = "KASPI_DELIVERY_RETURN_REQUESTED"
	ReturnedStatus        OrderStatus = "RETURNED"
)

type orderResponse struct {
	Data order `json:"data"`
}

func GetOrderStatus(ctx context.Context, token, orderId string) (OrderStatus, error) {
	link := fmt.Sprintf("/orders/%s", orderId)
	var parsedRes orderResponse
	err := do(ctx, token, http.MethodGet, link, &parsedRes)
	if err != nil {
		return "", err
	}
	return OrderStatus(parsedRes.Data.Attributes.Status), nil
}
