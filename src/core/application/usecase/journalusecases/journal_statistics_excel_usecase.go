package journalusecases

import (
	"github.com/xuri/excelize/v2"
	"slib.uz/src/core/domain/entity"
)

type JournalStatisticsExcelUseCase struct {
	listUseCase              *JournalStatisticsListUseCase
	completionPercentUseCase *JournalCompletionPercentUseCase
}

// @inject
func NewJournalStatisticsExcelUseCase(
	listUseCase *JournalStatisticsListUseCase,
	completionPercentUseCase *JournalCompletionPercentUseCase,
) *JournalStatisticsExcelUseCase {
	return &JournalStatisticsExcelUseCase{
		listUseCase:              listUseCase,
		completionPercentUseCase: completionPercentUseCase,
	}
}

func (this *JournalStatisticsExcelUseCase) Execute() ([]byte, error) {
	page := 1
	pageSize := 1000

	var allStatistics []*entity.JournalStatisticEntity

	for {
		result, err := this.listUseCase.Execute(page, pageSize)
		if err != nil {
			return nil, err
		}

		allStatistics = append(allStatistics, result.Items...)

		if len(result.Items) < pageSize {
			break
		}

		page++
	}

	journalIDs := make([]uint, len(allStatistics))
	for i, stat := range allStatistics {
		journalIDs[i] = stat.JournalID
	}

	completionPercentMap, err := this.completionPercentUseCase.ExecuteBatch(journalIDs)
	if err != nil {
		return nil, err
	}

	for i, stat := range allStatistics {
		if percent, exists := completionPercentMap[stat.JournalID]; exists {
			allStatistics[i].CompletionPercent = percent
		}
	}

	return this.generateExcel(allStatistics)
}

func (this *JournalStatisticsExcelUseCase) generateExcel(statistics []*entity.JournalStatisticEntity) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Statistics"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	headers := []string{
		"Journal Name",
		"Publisher Name",
		"ISSN Online",
		"ISSN Paper",
		"Completion %",
		"Has Study Fields",
		"Article Count",
		"Phone Number",
		"Telegram Contact",
		"Chief Editor ID",
		"Chief Editor Name",
		"Chief Editor Science ID",
		"Chief Editor Phone",
		"Secretaries Count",
	}

	for i, header := range headers {
		cell := string(rune('A'+i)) + "1"
		f.SetCellValue(sheetName, cell, header)
	}

	for i, stat := range statistics {
		row := i + 2
		cells := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P"}
		rowStr := string(rune('0'+row/10)) + string(rune('0'+row%10))
		f.SetCellValue(sheetName, cells[1]+rowStr, stat.JournalName)
		f.SetCellValue(sheetName, cells[3]+rowStr, stat.PublisherName)
		f.SetCellValue(sheetName, cells[4]+rowStr, stat.IssnOnline)
		f.SetCellValue(sheetName, cells[5]+rowStr, stat.IssnPaper)
		f.SetCellValue(sheetName, cells[6]+rowStr, stat.CompletionPercent)
		f.SetCellValue(sheetName, cells[7]+rowStr, stat.HasStudyFields)
		f.SetCellValue(sheetName, cells[8]+rowStr, stat.ArticleCount)
		f.SetCellValue(sheetName, cells[9]+rowStr, stat.JournalPhoneNumber)
		f.SetCellValue(sheetName, cells[10]+rowStr, stat.JournalTelegramContact)

		if stat.JournalChiefEditor != nil {
			f.SetCellValue(sheetName, cells[11]+rowStr, stat.JournalChiefEditor.UserID)
			f.SetCellValue(sheetName, cells[12]+rowStr, stat.JournalChiefEditor.FullName)
			f.SetCellValue(sheetName, cells[13]+rowStr, stat.JournalChiefEditor.ScienceID)
			f.SetCellValue(sheetName, cells[14]+rowStr, stat.JournalChiefEditor.PhoneNumber)
		}

		f.SetCellValue(sheetName, cells[15]+rowStr, len(stat.JournalSecretaries))
	}

	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
