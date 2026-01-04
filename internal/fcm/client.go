package fcm

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"

	"fcm_microservice/internal/model"
)

type Client struct {
	MsgClient *messaging.Client
}

// NewClient menginisialisasi koneksi ke Firebase menggunakan file JSON
func NewClient(credentialPath string) (*Client, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile(credentialPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting messaging client: %v", err)
	}

	return &Client{MsgClient: client}, nil
}

// Send mengirim notifikasi ke FCM
func (c *Client) Send(job model.FcmPayload) error {
	ctx := context.Background()

	// Bangun pesan
	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: job.Title,
			Body:  job.Body,
		},
		Data: job.Data,
	}

	// Tentukan apakah kirim ke Topic atau Token spesifik
	if job.IsTopic {
		msg.Topic = job.Target
	} else {
		msg.Token = job.Target
	}

	// Eksekusi kirim
	response, err := c.MsgClient.Send(ctx, msg)
	if err != nil {
		return err
	}

	log.Printf("FCM Sent ID: %s", response)
	return nil
}
