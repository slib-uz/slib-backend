package gateway

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

// @inject
func NewFirebaseMessagingClient() *messaging.Client {
	ctx := context.Background()
	opt := option.WithCredentialsFile("assets/firebase_sdk.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Println("error initializing firebase app")
		// panic(fmt.Sprintf("error initializing app: %v", err))
	}
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Println("error getting Messaging client")
		// panic(fmt.Sprintf("error getting Messaging client: %v", err))
	}
	return client
}
