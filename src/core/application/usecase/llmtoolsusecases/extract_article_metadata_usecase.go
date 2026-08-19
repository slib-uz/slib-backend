package llmtoolsusecases

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"unicode"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/domain/ports/pdf"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/config"
)

const (
	articleTextHeadRunes = 80000
	articleTextTailRunes = 20000
)

type ExtractArticleMetadataUseCase struct {
	extractor pdf.TextExtractor
	deepseek  gateway.DeepSeekGateway
	fields    repository.StudyFieldRepository
	maxSize   int64
}

// @inject
func NewExtractArticleMetadataUseCase(
	extractor pdf.TextExtractor,
	deepseek gateway.DeepSeekGateway,
	fields repository.StudyFieldRepository,
	env *config.Config,
) *ExtractArticleMetadataUseCase {
	return &ExtractArticleMetadataUseCase{
		extractor: extractor,
		deepseek:  deepseek,
		fields:    fields,
		maxSize:   int64(env.UploadedFileMaxSize),
	}
}

func (this *ExtractArticleMetadataUseCase) Execute(file *multipart.FileHeader, journalID uint, articleLanguage string) (*entity.ArticleMetadataExtraction, error) {
	articleLanguage = strings.ToLower(strings.TrimSpace(articleLanguage))
	if journalID == 0 || articleLanguage == "" {
		return nil, response.InvalidArgument
	}

	if file == nil {
		return nil, response.InvalidFileError
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".pdf" {
		return nil, response.InvalidFileError
	}

	uploadType, ok := enum.LookupUploadType(".pdf")
	if !ok {
		return nil, response.InvalidFileError
	}

	if this.maxSize > 0 && file.Size > this.maxSize {
		return nil, response.EntityToLargeError
	}

	data, err := readFileBytes(file, this.maxSize)
	if err != nil {
		return nil, response.InvalidFileError
	}
	if this.maxSize > 0 && int64(len(data)) > this.maxSize {
		return nil, response.EntityToLargeError
	}
	if !uploadType.MatchesMagic(data) {
		return nil, response.InvalidFileError
	}

	langs := entity.MetadataLangs(articleLanguage)

	text, err := this.extractor.Extract(data)
	if err != nil {
		return nil, response.InvalidFileError
	}
	if strings.TrimSpace(text) == "" {
		empty := entity.EmptyArticleMetadataExtraction(langs)
		empty.ArticleLanguage = articleLanguage
		return empty, nil
	}

	text = truncateArticleText(text, articleTextHeadRunes, articleTextTailRunes)

	tree, err := this.fields.GetByJournalID(journalID)
	if err != nil {
		return nil, err
	}
	catalog := entity.FlattenStudyFieldCatalog(tree)

	result, err := this.deepseek.ExtractArticleMetadata(text, catalog, langs)
	if err != nil {
		return nil, err
	}

	if result == nil {
		result = entity.EmptyArticleMetadataExtraction(langs)
	}

	result.StudyFieldIDs = entity.FilterStudyFieldIDs(result.StudyFieldIDs, entity.CatalogIDSet(catalog))
	result.Normalize(langs)
	result.ArticleLanguage = articleLanguage
	return result, nil
}

func readFileBytes(file *multipart.FileHeader, maxSize int64) ([]byte, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	limit := maxSize
	if limit <= 0 {
		limit = 10 << 20
	}
	return io.ReadAll(io.LimitReader(src, limit+1))
}

func truncateArticleText(s string, head, tail int) string {
	runes := []rune(s)
	if len(runes) <= head+tail {
		return s
	}
	var b strings.Builder
	b.WriteString(string(runes[:head]))
	b.WriteString("\n\n...\n\n")
	b.WriteString(string(runes[len(runes)-tail:]))
	return strings.TrimRightFunc(b.String(), unicode.IsSpace)
}
