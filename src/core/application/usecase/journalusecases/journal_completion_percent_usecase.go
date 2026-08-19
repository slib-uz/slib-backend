package journalusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalCompletionPercentUseCase struct {
	journalRepository repository.JournalRepository
}

// @inject
func NewJournalCompletionPercentUseCase(journalRepository repository.JournalRepository) *JournalCompletionPercentUseCase {
	return &JournalCompletionPercentUseCase{
		journalRepository: journalRepository,
	}
}

func (this *JournalCompletionPercentUseCase) Execute(journalID uint) (*entity.JournalCompletionResultEntity, error) {
	completionStats, err := this.journalRepository.GetJournalsCompletionStats([]uint{journalID})
	if err != nil {
		return nil, err
	}

	for _, cs := range completionStats {
		if cs.JournalID == journalID {
			result := calculateCompletionResult(cs)
			return &result, nil
		}
	}

	return &entity.JournalCompletionResultEntity{
		CompletionPercent: 0,
		IncompleteFields:  allCompletionFieldNames(),
	}, nil
}

func (this *JournalCompletionPercentUseCase) ExecuteBatch(journalIDs []uint) (map[uint]int, error) {
	result := make(map[uint]int)
	if len(journalIDs) == 0 {
		return result, nil
	}

	completionStats, err := this.journalRepository.GetJournalsCompletionStats(journalIDs)
	if err != nil {
		return nil, err
	}

	for _, cs := range completionStats {
		result[cs.JournalID] = calculateCompletionResult(cs).CompletionPercent
	}

	return result, nil
}

func calculateCompletionResult(cs entity.JournalCompletionStatsEntity) entity.JournalCompletionResultEntity {
	incompleteFields := make([]string, 0, 21)

	if !cs.HasName {
		incompleteFields = append(incompleteFields, "name")
	}
	if !cs.HasShortName {
		incompleteFields = append(incompleteFields, "short_name")
	}
	if !cs.HasJournalLanguages {
		incompleteFields = append(incompleteFields, "journal_languages")
	}
	if !cs.HasEstablishedDate {
		incompleteFields = append(incompleteFields, "established_date")
	}
	if !cs.HasPhoneNumber {
		incompleteFields = append(incompleteFields, "phone_number")
	}
	if !cs.HasEmail {
		incompleteFields = append(incompleteFields, "email")
	}
	if !cs.HasWebsite {
		incompleteFields = append(incompleteFields, "website")
	}
	if !cs.HasTelegramSupport {
		incompleteFields = append(incompleteFields, "telegram_support")
	}
	if !cs.HasIssnOnline {
		incompleteFields = append(incompleteFields, "issn_online")
	}
	if !cs.HasIssnPaper {
		incompleteFields = append(incompleteFields, "issn_paper")
	}
	if !cs.HasSocialNetworks {
		incompleteFields = append(incompleteFields, "social_networks")
	}
	if !cs.HasStudyFields {
		incompleteFields = append(incompleteFields, "study_fields")
	}
	if !cs.HasImageFile {
		incompleteFields = append(incompleteFields, "image_file")
	}
	if !cs.HasCertificateFile {
		incompleteFields = append(incompleteFields, "certificate_file")
	}
	if !cs.HasOakCertificateFile {
		incompleteFields = append(incompleteFields, "oak_certificate_file")
	}
	if !cs.HasAddress {
		incompleteFields = append(incompleteFields, "address")
	}
	if !cs.HasJournalIndexing {
		incompleteFields = append(incompleteFields, "journal_indexing")
	}
	if !cs.HasDescription {
		incompleteFields = append(incompleteFields, "description")
	}
	if !cs.HasRule {
		incompleteFields = append(incompleteFields, "rule")
	}
	if !cs.HasArticleConditions {
		incompleteFields = append(incompleteFields, "article_conditions")
	}
	if !cs.HasRegionID {
		incompleteFields = append(incompleteFields, "region_id")
	}

	totalFields := 21
	completedFields := totalFields - len(incompleteFields)

	return entity.JournalCompletionResultEntity{
		CompletionPercent: int(float64(completedFields) / float64(totalFields) * 100),
		IncompleteFields:  incompleteFields,
	}
}

func allCompletionFieldNames() []string {
	return []string{
		"name",
		"short_name",
		"journal_languages",
		"established_date",
		"phone_number",
		"email",
		"website",
		"telegram_support",
		"issn_online",
		"issn_paper",
		"social_networks",
		"study_fields",
		"image_file",
		"certificate_file",
		"oak_certificate_file",
		"address",
		"journal_indexing",
		"description",
		"rule",
		"article_conditions",
		"region_id",
	}
}
