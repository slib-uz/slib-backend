package journalstatisticsusecase

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
)

type JournalStatisticListExcelUseCase struct {
	listUseCase *JournalStatisticListUseCase
}

// @inject
func NewJournalStatisticListExcelUseCase(
	listUseCase *JournalStatisticListUseCase,
) *JournalStatisticListExcelUseCase {
	return &JournalStatisticListExcelUseCase{
		listUseCase: listUseCase,
	}
}

func (this *JournalStatisticListExcelUseCase) Execute(institutionID uint) (*excelize.File, error) {
	statistics, err := this.fetchAllStatistics(institutionID, 0, nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	return this.generateExcel(statistics)
}

func (this *JournalStatisticListExcelUseCase) fetchAllStatistics(
	institutionID, publisherID uint,
	name, description, issn, publisherName *string,
) ([]*entity.JournalStatisticV2Entity, error) {
	page := 1
	pageSize := 1000

	var allStatistics []*entity.JournalStatisticV2Entity

	for {
		result, err := this.listUseCase.Execute(page, pageSize, institutionID, publisherID, name, description, issn, publisherName)
		if err != nil {
			return nil, err
		}

		allStatistics = append(allStatistics, result.Items...)

		if len(result.Items) < pageSize {
			break
		}

		page++
	}

	return allStatistics, nil
}

func (this *JournalStatisticListExcelUseCase) generateExcel(statistics []*entity.JournalStatisticV2Entity) (*excelize.File, error) {
	f := excelize.NewFile()
	sheetName := fmt.Sprintf("Statistics %s", time.Now().Format("2006-01-02"))

	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	headers := []string{
		"Jurnal nomi",
		"Nashriyot nomi",
		"Muassasa nomi",
		"ISSN onlayn",
		"ISSN paper",
		"Hudud",
		"Maqola qabul qilish",
		"Adminlar soni",
		"Nashrlar",
		"Maqolalar",
		"Maqola arizalari",
		"Hammualliflar",
		"To'ldirilganlik foizi",
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	dataStyle, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{
			Horizontal: "left",
			Vertical:   "center",
		},
	})
	if err != nil {
		return nil, err
	}

	lastCol, _ := excelize.ColumnNumberToName(len(headers))

	for i, header := range headers {
		col, _ := excelize.ColumnNumberToName(i + 1)
		_ = f.SetCellValue(sheetName, col+"1", header)
	}

	for i, stat := range statistics {
		row := i + 2
		values := []interface{}{
			formatJournalName(stat.JournalName),
			stat.PublisherName,
			stat.InstitutionName,
			stat.IssnOnline,
			stat.IssnPaper,
			formatRegionName(stat.JournalRegionName),
			formatAccessType(stat.JournalSubmissionAccess),
			stat.AdminCount,
			stat.EditionCount,
			stat.ArticleCount,
			stat.ArticleApplicationCount,
			stat.CoAuthorCount,
			stat.CompletionPercent,
		}

		for colIdx, value := range values {
			col, _ := excelize.ColumnNumberToName(colIdx + 1)
			_ = f.SetCellValue(sheetName, col+fmt.Sprint(row), value)
		}
	}

	lastRow := len(statistics) + 1
	_ = f.SetCellStyle(sheetName, "A1", lastCol+"1", headerStyle)
	if lastRow > 1 {
		_ = f.SetCellStyle(sheetName, "A2", lastCol+fmt.Sprint(lastRow), dataStyle)
	}

	if err := f.AutoFitColWidth(sheetName, "A:"+lastCol); err != nil {
		return nil, err
	}

	return f, nil
}

func formatJournalName(name map[string]string) string {
	if value, ok := name["uz"]; ok {
		return value
	}
	return ""
}

func formatRegionName(regionName *string) string {
	if regionName == nil {
		return ""
	}
	return *regionName
}

func formatAccessType(accessType enum.AccessType) string {
	switch accessType {
	case enum.PublicAccessType:
		return "Ochiq"
	case enum.PrivateAccessType:
		return "Yopiq"
	default:
		return fmt.Sprint(int(accessType))
	}
}
