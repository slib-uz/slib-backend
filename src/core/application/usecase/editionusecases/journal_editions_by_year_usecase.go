package editionusecases

import (
	"context"
	"sort"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditionsListByYearUseCase struct {
	repository repository.EditionRepository
}

// @inject
func NewJournalEditionsListByYearUseCase(repository repository.EditionRepository) *JournalEditionsListByYearUseCase {
	return &JournalEditionsListByYearUseCase{repository: repository}
}

type EditionListResponse struct {
	Years   []YearWithEditions `json:"years"`
	Total   int                `json:"total"`
	Page    int                `json:"page"`
	PerPage int                `json:"per_page"`
}

type YearWithEditions struct {
	Year     int                    `json:"year"`
	Editions []entity.EditionEntity `json:"editions"`
}

func (this *JournalEditionsListByYearUseCase) Execute(ctx context.Context, journalID uint, page, pageSize int) (*EditionListResponse, error) {
	// Fetch all editions for the journal
	allEditions, err := this.repository.GetByJournalID(ctx, journalID, 1, 100000, "", 0)
	if err != nil {
		return nil, err
	}

	// Group editions by year
	yearMap := make(map[int][]entity.EditionEntity)
	for _, editionPtr := range allEditions.Items {
		if editionPtr != nil {
			year := editionPtr.PublishedAt.Year()
			yearMap[year] = append(yearMap[year], *editionPtr)
		}
	}

	// Convert map to slice and sort years in descending order
	var years []int
	for year := range yearMap {
		years = append(years, year)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(years)))

	// Calculate total years
	totalYears := len(years)

	// Apply pagination on years
	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize

	if startIdx >= totalYears {
		return &EditionListResponse{
			Years:   []YearWithEditions{},
			Total:   totalYears,
			Page:    page,
			PerPage: pageSize,
		}, nil
	}

	if endIdx > totalYears {
		endIdx = totalYears
	}

	// Build response for paginated years; editions within each year by published_at

	var responseYears []YearWithEditions
	for i := startIdx; i < endIdx; i++ {
		year := years[i]
		editions := yearMap[year]
		sort.Slice(editions, func(i, j int) bool {
			return editions[i].PublishedAt.After(editions[j].PublishedAt)
		})
		responseYears = append(responseYears, YearWithEditions{
			Year:     year,
			Editions: editions,
		})
	}

	return &EditionListResponse{
		Years:   responseYears,
		Total:   totalYears,
		Page:    page,
		PerPage: pageSize,
	}, nil
}
