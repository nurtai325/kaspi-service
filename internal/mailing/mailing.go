package mailing

import (
	"context"
	"errors"
	"fmt"
	"log"
	"maps"
	"time"

	"github.com/nurtai325/kaspi-service/internal/kaspi"
	"github.com/nurtai325/kaspi-service/internal/models"
	"github.com/nurtai325/kaspi-service/internal/repositories"
	"github.com/nurtai325/kaspi-service/internal/whatsapp"
)

const (
	mailingIntervalMinutes = 5
)

var (
	ErrMailingRun = errors.New("error running mailing")
)

func Run(clientRepo repositories.Client, messenger *whatsapp.Messenger) {
	var cycleRunningTime time.Duration
	for {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*mailingIntervalMinutes)
		defer cancel()
		timeOffset := cycleRunningTime + time.Minute*mailingIntervalMinutes
		from := time.Now().UTC().Add(-timeOffset)
		err := newOrders(ctx, clientRepo, messenger, from)
		if err != nil {
			log.Println(err)
		}
		err = completedOrders(ctx, messenger)
		if err != nil {
			log.Println(err)
		}
		cycleRunningTime = time.Now().Sub(start)
		time.Sleep(time.Minute * mailingIntervalMinutes)
	}
}

func newOrders(ctx context.Context, clientRepo repositories.Client, messenger *whatsapp.Messenger, from time.Time) error {
	clients, err := clientRepo.AllNotExpired(ctx)
	if err != nil {
		err = errors.Join(ErrMailingRun, err)
		return err
	}
	for _, client := range clients {
		if !client.Connected {
			log.Printf("warning client with id: %d, name: %s, phone: %s is not connected to whatsapp\n", client.Id, client.Name, client.Phone)
			continue
		}
		to := time.Now().UTC()
		orders, err := kaspi.AllOrders(ctx, client.Token, from, to)
		if err != nil {
			err = errors.Join(ErrMailingRun, fmt.Errorf("error getting all kaspi orders for client %s %s: %w", client.Name, client.Phone, err))
			return err
		}
		for _, order := range orders {
			log.Println(client, order)
			ordersQ.add(order.Id, models.QueuedOrder{
				ClientName:  client.Name,
				ClientPhone: client.Phone,
				ClientJid:   client.Jid,
				Order:       order,
				Token:       client.Token,
			})
			textMessage := newOrderMessage(client, order)
			err = messenger.Message(ctx, client.Jid, order.CustomerPhone, textMessage)
			if err != nil {
				err = errors.Join(ErrMailingRun, fmt.Errorf("error getting all kaspi orders for client %s %s: %w", client.Name, client.Phone, err))
				return err
			}
		}
	}
	return nil
}

func completedOrders(ctx context.Context, messenger *whatsapp.Messenger) error {
	ordersQ.mu.Lock()
	copiedOrdersQueue := make(map[string]models.QueuedOrder, len(ordersQ.queue))
	maps.Copy(copiedOrdersQueue, ordersQ.queue)
	ordersQ.mu.Unlock()
	for orderId, orderQ := range copiedOrdersQueue {
		status, err := kaspi.GetOrderStatus(ctx, orderQ.Token, orderId)
		if err != nil {
			return fmt.Errorf("getting order %s status: %w", orderId, err)
		}
		if failedOrder(status) {
			ordersQ.mu.Lock()
			delete(ordersQ.queue, orderId)
			ordersQ.mu.Unlock()
			continue
		} else if status != kaspi.Completed {
			continue
		}
		textMessage := completedOrderMessage(orderQ.ClientName, orderQ.Order)
		err = messenger.Message(ctx, orderQ.ClientJid, orderQ.Order.CustomerPhone, textMessage)
		if err != nil {
			return fmt.Errorf("sending message: %w", err)
		}
		ordersQ.mu.Lock()
		delete(ordersQ.queue, orderId)
		ordersQ.mu.Unlock()
	}
	return nil
}

func failedOrder(status kaspi.OrderStatus) bool {
	switch status {
	case kaspi.CancelledStatus:
		return true
	case kaspi.CancellingStatus:
		return true
	case kaspi.ReturnedStatus:
		return true
	case kaspi.ReturnRequestedStatus:
		return true
	default:
		return false
	}
}
