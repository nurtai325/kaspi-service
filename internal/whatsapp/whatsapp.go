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

type Messenger struct {
	container *sqlstore.Container
}

func NewMessenger(dbConn *sql.DB) *Messenger {
	container := sqlstore.NewWithDB(dbConn, "postgres", nil)
	err := container.Upgrade()
	if err != nil {
		log.Panic(fmt.Errorf("error making new sql whatsapp container: %w", err))
	}
	return &Messenger{container: container}
}

func (m *Messenger) Message(ctx context.Context, jid, to, text string) error {
	parsedJid, err := types.ParseJID(jid)
	if err != nil {
		return fmt.Errorf("invalid whatsapp device jid: %s: %w", jid, err)
	}
	device, err := m.container.GetDevice(parsedJid)
	if err != nil {
		return fmt.Errorf("can't find whatsapp device jid: %s: %w", jid, err)
	}
	client := whatsmeow.NewClient(device, nil)
	err = client.Connect()
	if err != nil {
		return fmt.Errorf("whatsapp client connect failure jid: %s: %w", jid, err)
	}
	_, err = client.SendMessage(ctx, types.NewJID(to, types.DefaultUserServer), &waE2E.Message{
		Conversation: proto.String(text),
	})
	if err != nil {
		return fmt.Errorf("whatsapp message sending error jid: %s: %w", jid, err)
	}
	return nil
}
