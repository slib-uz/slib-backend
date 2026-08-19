package handlers

import (
	"context"

	"slib.uz/src/core/application/usecase/spellcheckusecases"
	"slib.uz/src/core/domain/entity"
)

type SpellCheckTaskHandler struct {
	uc *spellcheckusecases.SpellCheckProcessUseCase
}

// @inject
func NewSpellCheckTaskHandler(uc *spellcheckusecases.SpellCheckProcessUseCase) *SpellCheckTaskHandler {
	return &SpellCheckTaskHandler{uc: uc}
}

func (h *SpellCheckTaskHandler) ProcessTask(_ context.Context, task *entity.AsyncTask) error {
	payload, err := UnmarshalPayload[entity.SpellCheckResultEntity](task)
	if err != nil {
		return err
	}
	return h.uc.Execute(&payload)
}
