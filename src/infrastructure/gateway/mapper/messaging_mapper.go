package mapper

import (
	"time"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/config"
)

func MessagingMapper(article map[string][]byte, env *config.Config) entity.MessageEntity {
	return entity.NewMessageEntity(
		article["key"],
		article["value"],
		time.Now(),
		article,
		env.KafkaTopic,
	)
}
