package mapper

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/infrastructure/persistence/models"
)

func EthicsPolicyModelListToEntityList(m []*models.EthicsPolicyModel) []*entity.EthicsPolicyEntity {
	var entities []*entity.EthicsPolicyEntity
	for _, model := range m {
		entities = append(entities, EthicsPolicyModelToEntity(model))
	}
	return entities
}

func EthicsPolicyModelToEntity(m *models.EthicsPolicyModel) *entity.EthicsPolicyEntity {
	return entity.NewEthicsPolicyEntity(
		m.ID,
		FromGormJson[map[string]string](m.Content),
	)
}
