package llmtoolsusecases_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/llmtoolsusecases"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/infrastructure/config"
)

type fakeExtractor struct {
	text string
	err  error
}

func (f fakeExtractor) Extract(_ []byte) (string, error) {
	return f.text, f.err
}

type fakeDeepSeek struct {
	lastText   string
	lastFields []entity.StudyFieldCatalogItem
	lastLangs  []string
	result     *entity.ArticleMetadataExtraction
	err        error
	calls      int
}

func (f *fakeDeepSeek) ExtractArticleMetadata(articleText string, studyFields []entity.StudyFieldCatalogItem, langs []string) (*entity.ArticleMetadataExtraction, error) {
	f.calls++
	f.lastText = articleText
	f.lastFields = studyFields
	f.lastLangs = langs
	return f.result, f.err
}

type fakeStudyFields struct {
	repository.StudyFieldRepository
	fields    []*entity.StudyFieldEntity
	journalID uint
}

func (f *fakeStudyFields) GetByJournalID(journalID uint) ([]*entity.StudyFieldEntity, error) {
	f.journalID = journalID
	return f.fields, nil
}

func makeFileHeader(t *testing.T, name string, content []byte) *multipart.FileHeader {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", name)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := fw.Write(content); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	r := multipart.NewReader(&buf, w.Boundary())
	form, err := r.ReadForm(int64(len(content)) + 1024)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}
	return form.File["file"][0]
}

func newUseCase(extractor fakeExtractor, gw *fakeDeepSeek, fields *fakeStudyFields) *llmtoolsusecases.ExtractArticleMetadataUseCase {
	return llmtoolsusecases.NewExtractArticleMetadataUseCase(
		extractor,
		gw,
		fields,
		&config.Config{UploadedFileMaxSize: 10 << 20},
	)
}

func TestRejectsMissingJournalOrLanguage(t *testing.T) {
	gw := &fakeDeepSeek{}
	uc := newUseCase(fakeExtractor{text: "x"}, gw, &fakeStudyFields{})
	fh := makeFileHeader(t, "maqola.pdf", []byte("%PDF-1.7\ncontent"))
	if _, err := uc.Execute(fh, 0, "uz"); !errors.Is(err, response.InvalidArgument) {
		t.Fatalf("journal_id=0: %v", err)
	}
	if _, err := uc.Execute(fh, 1, "  "); !errors.Is(err, response.InvalidArgument) {
		t.Fatalf("empty lang: %v", err)
	}
	if gw.calls != 0 {
		t.Fatalf("gateway chaqilmasligi kerak")
	}
}

func TestRejectsNonPdfExtension(t *testing.T) {
	gw := &fakeDeepSeek{}
	uc := newUseCase(fakeExtractor{text: "x"}, gw, &fakeStudyFields{})
	fh := makeFileHeader(t, "evil.html", []byte("<html></html>"))
	_, err := uc.Execute(fh, 1, "uz")
	if !errors.Is(err, response.InvalidFileError) {
		t.Fatalf("got %v", err)
	}
	if gw.calls != 0 {
		t.Fatalf("gateway chaqilmasligi kerak")
	}
}

func TestRejectsHtmlNamedPdf(t *testing.T) {
	gw := &fakeDeepSeek{}
	uc := newUseCase(fakeExtractor{text: "x"}, gw, &fakeStudyFields{})
	fh := makeFileHeader(t, "maqola.pdf", []byte("<html></html>"))
	_, err := uc.Execute(fh, 1, "uz")
	if !errors.Is(err, response.InvalidFileError) {
		t.Fatalf("got %v", err)
	}
}

func TestEmptyExtractedTextSkipsGateway(t *testing.T) {
	gw := &fakeDeepSeek{}
	uc := newUseCase(fakeExtractor{text: "   "}, gw, &fakeStudyFields{})
	fh := makeFileHeader(t, "maqola.pdf", []byte("%PDF-1.7\nempty"))
	got, err := uc.Execute(fh, 7, "de")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if gw.calls != 0 {
		t.Fatalf("gateway chaqilmasligi kerak")
	}
	if got.ArticleLanguage != "de" {
		t.Fatalf("lang=%s", got.ArticleLanguage)
	}
	if _, ok := got.ArticleName["de"]; !ok {
		t.Fatalf("de kaliti kerak: %#v", got.ArticleName)
	}
}

func TestHappyPathFiltersUnknownStudyFields(t *testing.T) {
	gw := &fakeDeepSeek{
		result: &entity.ArticleMetadataExtraction{
			ArticleName:   map[string]string{"uz": "Nom"},
			StudyFieldIDs: []uint{1, 99},
		},
	}
	fields := &fakeStudyFields{
		fields: []*entity.StudyFieldEntity{
			{
				ID:   1,
				Name: map[string]string{"uz": "IT"},
				Children: []*entity.StudyFieldEntity{
					{ID: 3, Name: map[string]string{"uz": "AI"}},
				},
			},
		},
	}
	uc := newUseCase(fakeExtractor{text: "maqola matni"}, gw, fields)
	fh := makeFileHeader(t, "maqola.pdf", []byte("%PDF-1.7\ncontent"))
	got, err := uc.Execute(fh, 42, "uz")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if fields.journalID != 42 {
		t.Fatalf("journalID=%d", fields.journalID)
	}
	if gw.calls != 1 {
		t.Fatalf("calls=%d", gw.calls)
	}
	if len(gw.lastFields) != 2 || gw.lastFields[1].ID != 3 {
		t.Fatalf("catalog=%v", gw.lastFields)
	}
	if len(got.StudyFieldIDs) != 1 || got.StudyFieldIDs[0] != 1 {
		t.Fatalf("ids=%v", got.StudyFieldIDs)
	}
	if got.ArticleLanguage != "uz" {
		t.Fatalf("lang=%s", got.ArticleLanguage)
	}
	if _, ok := got.ArticleName["de"]; ok {
		t.Fatalf("uz uchun de bo'lmasligi kerak: %#v", got.ArticleName)
	}
	if got.ArticleName["en"] != "" {
		t.Fatalf("normalize qilinishi kerak: %#v", got.ArticleName)
	}
}

func TestGermanAddsDeKey(t *testing.T) {
	gw := &fakeDeepSeek{
		result: &entity.ArticleMetadataExtraction{
			ArticleName: map[string]string{"uz": "Nom", "de": "Name"},
		},
	}
	uc := newUseCase(fakeExtractor{text: "text"}, gw, &fakeStudyFields{})
	fh := makeFileHeader(t, "maqola.pdf", []byte("%PDF-1.7\ncontent"))
	got, err := uc.Execute(fh, 1, "DE")
	if err != nil {
		t.Fatalf("err=%v", err)
	}
	if got.ArticleLanguage != "de" {
		t.Fatalf("lang=%s", got.ArticleLanguage)
	}
	if got.ArticleName["de"] != "Name" {
		t.Fatalf("name=%v", got.ArticleName)
	}
	if len(gw.lastLangs) != 4 || gw.lastLangs[3] != "de" {
		t.Fatalf("langs=%v", gw.lastLangs)
	}
}

func TestExtractErrorIsInvalidFile(t *testing.T) {
	gw := &fakeDeepSeek{}
	uc := newUseCase(fakeExtractor{err: errors.New("broken")}, gw, &fakeStudyFields{})
	fh := makeFileHeader(t, "maqola.pdf", []byte("%PDF-1.7\ncontent"))
	_, err := uc.Execute(fh, 1, "uz")
	if !errors.Is(err, response.InvalidFileError) {
		t.Fatalf("got %v", err)
	}
}
