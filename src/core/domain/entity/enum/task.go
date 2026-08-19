package enum

type Task string

var (
	TaskSpellcheck              Task = "SPELLCHECK"
	TaskAntiPlagStatusUpdate    Task = "ANTIPLAG.STATUS.UPDATE"
	TaskSendNotification        Task = "SEND.NOTIFICATION"
	TaskUpdateArticleViewsCount Task = "ARTICLE.VIEWS.COUNT.UPDATE"
	TaskUpdateNewsViewsCount    Task = "NEWS.VIEWS.COUNT.UPDATE"
	TaskUpdateJournalViewsCount Task = "JOURNAL.VIEWS.COUNT.UPDATE"
	TaskDeadlineReminder        Task = "DEADLINE.REMIND"
	TaskSetRoi                  Task = "SET.ROI"
	TaskUpdateRoi               Task = "UPDATE.ROI"
	TaskPublishRoi              Task = "PUBLISH.ROI"
	TaskSendArticleKafka        Task = "SEND.ARTICLE.KAFKA"
	TaskSendArticleToKafka      Task = "SEND.ARTICLE.TO.KAFKA"
	TaskSendTelegramAlert       Task = "SEND.TELEGRAM.ALERT"
	TaskOutboxEvent             Task = "OUTBOX.EVENT"
	TaskAiDetectStatusUpdate    Task = "AI.DETECT.STATUS.UPDATE"
)
