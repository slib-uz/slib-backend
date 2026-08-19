package service

import (
	"bytes"
	"strconv"
	"text/template"
	"time"

	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type SendEmailService struct {
	notificationRepo repository.NotificationRepository
	sendNotification *tasks.SendNotificationTask
	templateRepo     repository.NotificationTemplateRepository
}

// @inject
func NewSendEmailService(
	notificationRepo repository.NotificationRepository,
	sendNotification *tasks.SendNotificationTask,
	templateRepo repository.NotificationTemplateRepository,
) *SendEmailService {
	return &SendEmailService{
		notificationRepo: notificationRepo,
		sendNotification: sendNotification,
		templateRepo:     templateRepo,
	}
}

// --- Public

func (this *SendEmailService) NewSupportQuestion(userID []uint, chatID uint, message string) error {
	d := map[string]string{"ChatID": strconv.Itoa(int(chatID)), "Message": message}
	return this.send(userID, enum.TemplateNewSupportQuestion, d, d,
		map[string]string{"chat_id": strconv.Itoa(int(chatID))},
	)
}

func (this *SendEmailService) NewSupportAnswer(userID []uint, chatID uint, message string) error {
	d := map[string]string{"ChatID": strconv.Itoa(int(chatID)), "Message": message}
	return this.send(userID, enum.TemplateNewSupportAnswer, d, d,
		map[string]string{"chat_id": strconv.Itoa(int(chatID))},
	)
}

func (this *SendEmailService) send(
	userIds []uint,
	tplKey enum.NotificationTemplate,
	titleData, bodyData map[string]string,
	extraData map[string]string,
) error {
	_template, err := this.templateRepo.GetByKey(tplKey)
	if err != nil {
		return err
	}
	notifications, err := this.buildNotification(userIds, _template.Title, _template.Body, extraData, titleData, bodyData)
	if err != nil {
		return err
	}
	for _, n := range notifications {
		_ = this.sendNotification.Run(n)
	}
	return nil
}

func (this *SendEmailService) buildNotification(
	userIds []uint,
	_title map[string]string,
	_body map[string]string,
	extraData map[string]string,
	titleData map[string]string,
	bodyData map[string]string,
) ([]*entity.NotificationEntity, error) {
	title, err := this.parseLocale(_title, titleData)
	if err != nil {
		return nil, err
	}
	body, err := this.parseLocale(_body, bodyData)
	if err != nil {
		return nil, err
	}

	list := make([]*entity.NotificationEntity, len(userIds))

	for i, id := range userIds {
		list[i] = entity.NewNotificationEntity(
			0,
			&id,
			title,
			body,
			nil,
			extraData,
			true, // IsEmail = true
			false,
			false,
			false,
			time.Now(),
		)
	}

	return list, nil
}

func (this *SendEmailService) parseLocale(locale map[string]string, data any) (map[string]string, error) {
	uz, err := this._parseTemplate(locale["uz"], data)
	ru, _ := this._parseTemplate(locale["ru"], data)
	en, _ := this._parseTemplate(locale["en"], data)
	if err != nil {
		return nil, err
	}
	return map[string]string{"uz": uz, "ru": ru, "en": en}, nil
}

func (this *SendEmailService) _parseTemplate(pattern string, data any) (string, error) {
	tmpl, err := template.New("notification").Parse(pattern)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
