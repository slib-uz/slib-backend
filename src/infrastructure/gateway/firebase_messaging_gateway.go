package gateway

import (
	"context"
	"fmt"

	"firebase.google.com/go/v4/messaging"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/utils"
)

type FirebaseMessagingGateway struct {
	client *messaging.Client
}

// @inject
func NewFirebaseMessagingGateway(client *messaging.Client) gateway.NotificationGateway {
	return &FirebaseMessagingGateway{client: client}
}

func (this *FirebaseMessagingGateway) SendToTokens(id uint, title, body string, extraData map[string]string, tokens []string) (*entity.NotificationSendResultEntity, error) {

	message := &messaging.MulticastMessage{
		Notification: &messaging.Notification{
			Title: utils.TruncateUTF8(title, 256),
			Body:  utils.TruncateUTF8(body, 1024),
		},
		Tokens: tokens,
		Data:   extraData,
	}
	resp, err := this.client.SendEachForMulticast(context.Background(), message)
	if err != nil {
		return nil, err
	}
	var failedTokens, errorMessages []string

	for i, result := range resp.Responses {
		if !result.Success && result.Error != nil {
			failedTokens = append(failedTokens, tokens[i])
			errorMessages = append(errorMessages, result.Error.Error())
		}
	}

	return entity.NewNotificationSendResultEntity(
		0,
		id,
		resp.SuccessCount,
		resp.FailureCount,
		failedTokens,
		errorMessages,
	), nil
}

func (this *FirebaseMessagingGateway) SendToTopic(title, body string, extraData map[string]string, topic enum.NotificationTopic) error {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: utils.TruncateUTF8(title, 256),
			Body:  utils.TruncateUTF8(body, 1024),
		},
		Data:  extraData,
		Topic: string(topic),
	}

	_, err := this.client.Send(context.Background(), message)

	return err
}

func (this *FirebaseMessagingGateway) SubscribeToTopic(topic enum.NotificationTopic, token string) error {
	info, err := this.client.SubscribeToTopic(context.Background(), []string{token}, string(topic))
	if err != nil {
		return err
	}

	if len(info.Errors) > 0 {
		return fmt.Errorf("%v", info.Errors[0].Reason)
	}

	return nil
}

func (this *FirebaseMessagingGateway) UnsubscribeFromTopic(topic enum.NotificationTopic, token string) error {
	info, err := this.client.UnsubscribeFromTopic(context.Background(), []string{token}, string(topic))
	if err != nil {
		return err
	}

	if len(info.Errors) > 0 {
		return fmt.Errorf("%v", info.Errors[0].Reason)
	}

	return nil
}
