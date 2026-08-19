package enum

type ServiceCode string

const (
	ServiceSavodxon       ServiceCode = "SAVODXON"
	ServiceAntiPlag       ServiceCode = "ANTI_PLAG"
	ServiceAIDetection    ServiceCode = "AI_DETECTION"
	ServiceAIETaqriz      ServiceCode = "AI_ETAQRIZ"
	ServiceJournalWebsite ServiceCode = "JOURNAL_WEBSITE"
	ServiceDOI            ServiceCode = "DOI"
)

func (this ServiceCode) String() string {
	return string(this)
}

func (this ServiceCode) IsValid() bool {
	switch this {
	case ServiceSavodxon, ServiceAntiPlag, ServiceAIDetection, ServiceAIETaqriz, ServiceJournalWebsite, ServiceDOI:
		return true
	}
	return false
}
