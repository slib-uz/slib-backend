package mapper

import (
	entity2 "slib.uz/src/core/domain/entity"
)

func JournalEntityToOutput(journal *entity2.JournalEntity) *entity2.ROIJournalEntity {
	var issnPaper, issnOnline string
	if journal.ISSNPaper != nil {
		issnPaper = *journal.ISSNPaper
	}
	if journal.ISSNOnline != nil {
		issnOnline = *journal.ISSNOnline
	}
	shortName := journal.ShortName
	if shortName == nil {
		shortName = make(map[string]string)
	}

	return entity2.NewROIJournalEntity(
		journal.Name,
		shortName,
		journal.DateOfEstablished.Format("2006-01-02"),
		issnPaper,
		issnOnline,
		journal.PublisherID,
	)
}
