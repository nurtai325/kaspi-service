package whatsapp

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"google.golang.org/protobuf/proto"
)

var (
	container *sqlstore.Container
)

type Messenger struct {
	container *sqlstore.Container
}

func NewMessenger(dbConn *sql.DB) *Messenger {
	if container == nil {
		newContainer := sqlstore.NewWithDB(dbConn, "postgres", nil)
		err := newContainer.Upgrade()
		if err != nil {
			log.Panic(fmt.Errorf("error making new sql whatsapp container: %w", err))
		}
		container = newContainer
	}
	return &Messenger{container: container}
}

func (m *Messenger) Message(ctx context.Context, jid, to, text string) error {
	parsedJid, err := types.ParseJID(jid)
	if err != nil {
		return fmt.Errorf("invalid whatsapp device jid: %s: %w", jid, err)
	}
	device, err := m.container.GetDevice(parsedJid)
	if err != nil || device == nil {
		return fmt.Errorf("can't find whatsapp device jid: %s: %w", jid, err)
	}
	client := whatsmeow.NewClient(device, nil)
	err = client.Connect()
	if err != nil {
		return fmt.Errorf("whatsapp client connect failure jid: %s: %w", jid, err)
	}
	defer client.Disconnect()
	_, err = client.SendMessage(ctx, types.NewJID(to, types.DefaultUserServer), &waE2E.Message{
		Conversation: proto.String(text),
	})
	if err != nil {
		return fmt.Errorf("whatsapp message sending error jid: %s to: %s: %w", jid, to, err)
	}
	return nil
}
