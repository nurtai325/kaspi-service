package whatsapp

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nurtai325/kaspi-service/internal/repositories"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
)

type pairingData struct {
	ImagePath string
}

func GetQr(clientRepo repositories.Client, id int) (pairingData, error) {
	dev := container.NewDevice()
	cli := whatsmeow.NewClient(dev, nil)
	qrCh, _ := cli.GetQRChannel(context.Background())
	err := cli.Connect()
	if err != nil {
		return pairingData{}, fmt.Errorf("error connecting to whatsapp websocket: %w", err)
	}
	for evt := range qrCh {
		if evt.Event != "code" {
			continue
		}
		imagePath := ""
		go func(qrCh <-chan whatsmeow.QRChannelItem) {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute*8)
			defer cancel()
			for evt := range qrCh {
				if evt.Error != nil {
					err := fmt.Errorf("image path: %s event: %s client: %d: %w", imagePath, evt.Event, id, evt.Error)
					log.Println(err)
					return
				}
				if evt.Event == "timeout" {
					err := fmt.Errorf("qr channel timed out. image path: %s event: %s client: %d: %w", imagePath, evt.Event, id, evt.Error)
					log.Println(err)
					return
				}
				if evt.Event == "success" {
					err := clientRepo.ConnectWh(ctx, id, cli.Store.ID.String())
					if err != nil {
						log.Println(err)
					}
					return
				}
			}
			log.Println(fmt.Errorf("qr channel closed unexpectedly"))
		}(qrCh)
		imagePath = fmt.Sprintf("/assets/%d-%d.qr.png", id, time.Now().UnixMilli())
		err = qrcode.WriteFile(evt.Code, qrcode.Medium, 512, "."+imagePath)
		if err != nil {
			return pairingData{}, fmt.Errorf("error generating qr code image: %w", err)
		}
		return pairingData{ImagePath: imagePath}, nil
	}
	return pairingData{}, fmt.Errorf("unsuccesfull pairing")
}
