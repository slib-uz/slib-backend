package asynctask

import (
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/async"
	"slib.uz/src/entrypoint/asynctask/handlers"
)

type App struct {
	server async.AsyncServer

	// handlers
	spellcheckHandler              *handlers.SpellCheckTaskHandler
	antiPlagStatusUpdateHandler    *handlers.AntiPlagStatusUpdateTaskHandler
	aiDetectStatusUpdateHandler    *handlers.AiDetectStatusUpdateTaskHandler
	sendNotificationHandler        *handlers.SendNotificationTaskHandler
	articleViewsCountUpdateHandler *handlers.ArticleViewsCountUpdateTaskHandler
	newsViewsCountUpdateHandler    *handlers.NewsViewsCountUpdateTaskHandler
	journalViewsCountUpdateHandler *handlers.JournalViewsCountUpdateTaskHandler
	deadlineReminderTaskHandler    *handlers.DeadlineReminderTaskHandler
	sendArticleKafkaTaskHandler    *handlers.SendArticleKafkaTaskHandler
	sendTelegramAlertTaskHandler   *handlers.SendTelegramAlertTaskHandler
	outboxEventTaskHandler         *handlers.OutboxEventTaskHandler
	roiSenderTaskHandler           *handlers.RoiSenderTaskHandler
	roiUpdateTaskHandler           *handlers.UpdateRoiTaskHandler
	roiPublishTaskHandler          *handlers.PublishRoiTaskHandler
}

// @inject
func NewApp(
	server async.AsyncServer,
	spellcheckHandler *handlers.SpellCheckTaskHandler,
	antiPlagStatusUpdateHandler *handlers.AntiPlagStatusUpdateTaskHandler,
	aiDetectStatusUpdateHandler *handlers.AiDetectStatusUpdateTaskHandler,
	sendNotificationHandler *handlers.SendNotificationTaskHandler,
	articleViewsCountUpdateHandler *handlers.ArticleViewsCountUpdateTaskHandler,
	newsViewsCountUpdateHandler *handlers.NewsViewsCountUpdateTaskHandler,
	journalViewsCountUpdateHandler *handlers.JournalViewsCountUpdateTaskHandler,
	deadlineReminderTaskHandler *handlers.DeadlineReminderTaskHandler,
	sendArticleKafkaTaskHandler *handlers.SendArticleKafkaTaskHandler,
	sendTelegramAlertTaskHandler *handlers.SendTelegramAlertTaskHandler,
	outboxEventTaskHandler *handlers.OutboxEventTaskHandler,
	roiSenderTaskHandler *handlers.RoiSenderTaskHandler,
	roiUpdateTaskHandler *handlers.UpdateRoiTaskHandler,
	roiPublishTaskHandler *handlers.PublishRoiTaskHandler,
) *App {
	return &App{
		server:                         server,
		spellcheckHandler:              spellcheckHandler,
		antiPlagStatusUpdateHandler:    antiPlagStatusUpdateHandler,
		aiDetectStatusUpdateHandler:    aiDetectStatusUpdateHandler,
		sendNotificationHandler:        sendNotificationHandler,
		articleViewsCountUpdateHandler: articleViewsCountUpdateHandler,
		newsViewsCountUpdateHandler:    newsViewsCountUpdateHandler,
		journalViewsCountUpdateHandler: journalViewsCountUpdateHandler,
		deadlineReminderTaskHandler:    deadlineReminderTaskHandler,
		sendArticleKafkaTaskHandler:    sendArticleKafkaTaskHandler,
		sendTelegramAlertTaskHandler:   sendTelegramAlertTaskHandler,
		outboxEventTaskHandler:         outboxEventTaskHandler,
		roiSenderTaskHandler:           roiSenderTaskHandler,
		roiUpdateTaskHandler:           roiUpdateTaskHandler,
		roiPublishTaskHandler:          roiPublishTaskHandler,
	}
}

func (this *App) Init() {
	this.server.Init()
	this.initHandlers()
}

func (this *App) initHandlers() {
	this.server.HandlerFunc(enum.TaskSpellcheck, this.spellcheckHandler)
	this.server.HandlerFunc(enum.TaskAntiPlagStatusUpdate, this.antiPlagStatusUpdateHandler)
	this.server.HandlerFunc(enum.TaskAiDetectStatusUpdate, this.aiDetectStatusUpdateHandler)
	this.server.HandlerFunc(enum.TaskSendNotification, this.sendNotificationHandler)
	this.server.HandlerFunc(enum.TaskUpdateArticleViewsCount, this.articleViewsCountUpdateHandler)
	this.server.HandlerFunc(enum.TaskUpdateNewsViewsCount, this.newsViewsCountUpdateHandler)
	this.server.HandlerFunc(enum.TaskUpdateJournalViewsCount, this.journalViewsCountUpdateHandler)
	this.server.HandlerFunc(enum.TaskDeadlineReminder, this.deadlineReminderTaskHandler)
	this.server.HandlerFunc(enum.TaskSendArticleKafka, this.sendArticleKafkaTaskHandler)
	this.server.HandlerFunc(enum.TaskSendTelegramAlert, this.sendTelegramAlertTaskHandler)
	this.server.HandlerFunc(enum.TaskOutboxEvent, this.outboxEventTaskHandler)
	this.server.HandlerFunc(enum.TaskSetRoi, this.roiSenderTaskHandler)
	this.server.HandlerFunc(enum.TaskUpdateRoi, this.roiUpdateTaskHandler)
	this.server.HandlerFunc(enum.TaskPublishRoi, this.roiPublishTaskHandler)
}

func (this *App) Start() {
	if err := this.server.Run(); err != nil {
		panic(err)
	}
}
