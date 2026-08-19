package tasks_test

import (
	"errors"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"slib.uz/src/core/application/tasks"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/logger"

	portpublisher "slib.uz/src/core/domain/ports/tasks/publisher"
)

type publisherPort = portpublisher.TaskPublisher

// failingPublisher — Publish har doim xato qaytaradigan soxta nashr etuvchi.
type failingPublisher struct {
	publisherPort
	err error
}

func (this *failingPublisher) Publish(task *entity.TaskEntity[any], maxRetryCount int) error {
	return this.err
}

func observedLogger() (*logger.AsyncLogger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.ErrorLevel)
	return &logger.AsyncLogger{Logger: zap.New(core)}, logs
}

// TestRoiSendersSwallowPublishErrorAndLogIt — uchala ROI sender'i ham
// navbatga qo'yish xatosini chaqiruvchiga qaytarmaydi (ROI yuborish
// yordamchi qadam), lekin uni jimgina yo'qotmaydi ham: xato zap orqali
// ERROR darajasida yoziladi.
func TestRoiSendersSwallowPublishErrorAndLogIt(t *testing.T) {
	publishErr := errors.New("navbat mavjud emas")

	tests := []struct {
		name        string
		run         func(pub *failingPublisher, log *logger.AsyncLogger) error
		idFieldKey  string
		idFieldWant int64
	}{
		{
			name: "UpdateRoiSenderTask",
			run: func(pub *failingPublisher, log *logger.AsyncLogger) error {
				return tasks.NewUpdateRoiSenderTask(pub, log).
					Run(tasks.UpdateRoiSenderPayload{ArticleID: 42})
			},
			idFieldKey:  "article_id",
			idFieldWant: 42,
		},
		{
			name: "PublishRoiSenderTask",
			run: func(pub *failingPublisher, log *logger.AsyncLogger) error {
				return tasks.NewPublishRoiSenderTask(pub, log).
					Run(tasks.PublishRoiSenderPayload{ArticleID: 42})
			},
			idFieldKey:  "article_id",
			idFieldWant: 42,
		},
		{
			name: "SetRoiSenderTask",
			run: func(pub *failingPublisher, log *logger.AsyncLogger) error {
				return tasks.NewSetRoiSenderTask(pub, log).
					Run(tasks.SetRoiSenderPayload{ApplicationID: 42})
			},
			idFieldKey:  "application_id",
			idFieldWant: 42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			log, logs := observedLogger()
			pub := &failingPublisher{err: publishErr}

			if err := tt.run(pub, log); err != nil {
				t.Fatalf("%s.Run xato qaytardi: %v (nil kutilgan edi — "+
					"ROI navbati asosiy oqimni yiqitmasligi kerak)", tt.name, err)
			}

			entries := logs.All()
			if len(entries) != 1 {
				t.Fatalf("%s: %d ta ERROR yozuvi yozildi, 1 ta kutilgan edi",
					tt.name, len(entries))
			}
			if !containsErrorField(entries[0], publishErr) {
				t.Errorf("%s: log yozuvida asl xato yo'q: %+v",
					tt.name, entries[0].ContextMap())
			}
			if !containsInt64Field(entries[0], tt.idFieldKey, tt.idFieldWant) {
				t.Errorf("%s: log yozuvida %q maydoni (qiymati %d) yo'q: %+v",
					tt.name, tt.idFieldKey, tt.idFieldWant, entries[0].ContextMap())
			}
		})
	}
}

func containsErrorField(entry observer.LoggedEntry, want error) bool {
	for _, field := range entry.Context {
		if field.Key == "error" && field.Interface == want {
			return true
		}
	}
	return false
}

// containsInt64Field entry.Context'da berilgan kalit va qiymatga ega
// zap.Uint maydoni borligini tekshiradi. zap.Uint qiymatni Integer
// (int64) sifatida saqlaydi, tur esa zapcore.Uint... Type bo'ladi.
func containsInt64Field(entry observer.LoggedEntry, key string, want int64) bool {
	for _, field := range entry.Context {
		if field.Key == key && field.Integer == want {
			return true
		}
	}
	return false
}
