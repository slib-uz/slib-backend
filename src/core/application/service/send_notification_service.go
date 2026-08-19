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

type SendNotificationService struct {
	notificationRepo repository.NotificationRepository
	sendNotification *tasks.SendNotificationTask
	templateRepo     repository.NotificationTemplateRepository
}

// @inject
func NewSendNotificationService(notificationRepo repository.NotificationRepository, sendNotification *tasks.SendNotificationTask, templateRepo repository.NotificationTemplateRepository) *SendNotificationService {
	return &SendNotificationService{notificationRepo: notificationRepo, sendNotification: sendNotification, templateRepo: templateRepo}
}

// --- Public

func (this *SendNotificationService) PaymentSuccessful(userID []uint, appID uint, appNum string) error {
	return this.send(userID, enum.TemplateApplicationPaymentSuccessful,
		map[string]string{"Number": appNum}, map[string]string{"Number": appNum},
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.PaymentStage.ToString(), "status": enum.StatusAccepted.ToString()},
	)
}

func (this *SendNotificationService) PayToPublish(userID []uint, appID uint, appNum string, amount int) error {
	d := map[string]string{"Number": appNum, "Amount": strconv.Itoa(amount)}
	return this.send(userID, enum.TemplatePayToPublishArticle, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.PaymentStage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) RejectedAtPeerReview(userIds []uint, appID uint, appNum, reason string) error {
	d := map[string]string{"Number": appNum, "Reason": reason}
	return this.send(userIds, enum.TemplateRejectedAtPeerReview, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.PeerReviewStage.ToString(), "status": enum.StatusRejected.ToString()},
	)
}

func (this *SendNotificationService) AcceptedAtPeerReview(userIds []uint, appID uint, appNum, reason string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userIds, enum.TemplateAcceptedAtPeerReview, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.PeerReviewStage.ToString(), "status": enum.StatusAccepted.ToString()},
	)
}

func (this *SendNotificationService) ApplicationPublished(userID []uint, appID uint, appNum string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateApplicationPublished, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.PublishStage.ToString(), "status": enum.StatusAccepted.ToString()},
	)
}

func (this *SendNotificationService) NewApplication(userID []uint, appID uint, appNum string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateNewApplication, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.TechnicalReviewStage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) AntiPlagCheckFinished(userID []uint, appID uint, appNum string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateAntiPlagCheckFinished, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.TechnicalReviewStage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) AiDetectCheckFinished(userID []uint, appID uint, appNum string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateAiDetectCheckFinished, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.TechnicalReviewStage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) SpellCheckFinished(userID []uint, appID uint, appNum string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateSpellCheckFinished, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.TechnicalReviewStage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) NewPeerReview(userID []uint, appID uint, appNum string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateNewPeerReview, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.PeerReviewStage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) ReviewDeadlineRemind(stage enum.Stage, userID []uint, appID uint, appNum string, deadline time.Time) error {

	duration := time.Until(deadline)

	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24

	data := map[string]string{
		"Number":   appNum,
		"Days":     strconv.Itoa(days),
		"Hours":    strconv.Itoa(hours),
		"Deadline": deadline.Format("02.01.2006 15:04"),
	}

	var _template enum.NotificationTemplate
	switch stage {
	case enum.TechnicalReviewStage:
		_template = enum.TemplateDeadlineDateAtTechnicalReview
	case enum.PeerReviewStage:
		_template = enum.TemplateDeadlineDateAtPeerReview
	case enum.PaymentStage:
		_template = enum.TemplateDeadlineDateAtPaymentStage
	case enum.PublishStage:
		_template = enum.TemplateDeadlineDateAtPublishStage
	}

	return this.send(userID, _template, data, data,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": stage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) ReSent(userID []uint, appID uint, appNum string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateApplicationReSent, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.TechnicalReviewStage.ToString(), "status": enum.StatusPending.ToString()},
	)
}

func (this *SendNotificationService) RejectedAtTechnicalReview(userID []uint, appID uint, appNum, reason string) error {
	td := map[string]string{"Number": appNum}
	bd := map[string]string{"Number": appNum, "Reason": reason}
	return this.send(userID, enum.TemplateRejectedAtTechnicalReview, td, bd,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.TechnicalReviewStage.ToString(), "status": enum.StatusRejected.ToString()},
	)
}

func (this *SendNotificationService) AcceptedAtTechnicalReview(userID []uint, appID uint, appNum, reason string) error {
	d := map[string]string{"Number": appNum}
	return this.send(userID, enum.TemplateAcceptedAtTechnicalReview, d, d,
		map[string]string{"application_id": strconv.Itoa(int(appID)), "stage": enum.TechnicalReviewStage.ToString(), "status": enum.StatusAccepted.ToString()},
	)
}

func (this *SendNotificationService) NewSupportQuestion(userID []uint, chatID uint, message string) error {
	td := map[string]string{"ChatID": strconv.Itoa(int(chatID))}
	bd := map[string]string{"Message": message}
	return this.send(userID, enum.TemplateNewSupportQuestion, td, bd,
		map[string]string{"chat_id": strconv.Itoa(int(chatID))},
	)
}

func (this *SendNotificationService) NewSupportAnswer(userID []uint, chatID uint, message string) error {
	td := map[string]string{"ChatID": strconv.Itoa(int(chatID))}
	bd := map[string]string{"Message": message}
	return this.send(userID, enum.TemplateNewSupportAnswer, td, bd,
		map[string]string{"chat_id": strconv.Itoa(int(chatID))},
	)
}

// --- Universal notification sender
func (this *SendNotificationService) send(
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

// --- notification builder
// build and create notifications for users

func (this *SendNotificationService) buildNotification(
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

	notifications := make([]*entity.NotificationEntity, len(userIds))

	for i, id := range userIds {
		notifications[i] = entity.NewNotificationEntity(0, &id, title, body, nil, extraData, false, false, false, false, time.Now())
	}

	ids, err := this.notificationRepo.BulkCreate(notifications)
	if err != nil {
		return nil, err
	}
	for i, id := range ids {
		extraData["id"] = strconv.Itoa(int(id))
		notifications[i].ExtraData = extraData
		notifications[i].ID = id
	}

	return notifications, nil
}

func (this *SendNotificationService) parseLocale(locale map[string]string, data any) (map[string]string, error) {
	uz, err := this._parseTemplate(locale["uz"], data)
	ru, _ := this._parseTemplate(locale["ru"], data)
	en, _ := this._parseTemplate(locale["en"], data)
	if err != nil {
		return nil, err
	}
	return map[string]string{"uz": uz, "ru": ru, "en": en}, nil
}

func (this *SendNotificationService) _parseTemplate(pattern string, data any) (string, error) {
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
