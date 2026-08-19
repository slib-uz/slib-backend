package response

import (
	"encoding/xml"
	"fmt"
	"strings"

	"slib.uz/src/core/domain/entity/enum"
)

type CrossRefDiagnostic struct {
	SubmissionID string
	Status       enum.DoiDepositStatus
	DOI          string
	Message      string
}

func ParseCrossRefDiagnostic(body string) (*CrossRefDiagnostic, error) {
	xmlStart := strings.Index(body, "<?xml")
	if xmlStart == -1 {
		xmlStart = strings.Index(body, "<doi_batch_diagnostic")
	}
	if xmlStart == -1 {
		return nil, fmt.Errorf("no XML content found in body")
	}

	var d doiBatchDiagnostic
	if err := xml.Unmarshal([]byte(body[xmlStart:]), &d); err != nil {
		return nil, err
	}

	status := enum.DoiDepositStatusSuccess
	if d.BatchData.FailureCount > 0 {
		status = enum.DoiDepositStatusFailure
	}

	var doi string
	var messages []string
	for _, r := range d.Records {
		if r.DOI != "" && doi == "" {
			doi = r.DOI
		}
		if r.Msg != "" {
			messages = append(messages, r.Msg)
		}
	}

	return &CrossRefDiagnostic{
		SubmissionID: d.SubmissionID,
		Status:       status,
		DOI:          doi,
		Message:      strings.Join(messages, "; "),
	}, nil
}

type doiBatchDiagnostic struct {
	XMLName      xml.Name           `xml:"doi_batch_diagnostic"`
	Status       string             `xml:"status,attr"`
	SubmissionID string             `xml:"submission_id"`
	BatchID      string             `xml:"batch_id"`
	Records      []recordDiagnostic `xml:"record_diagnostic"`
	BatchData    batchData          `xml:"batch_data"`
}

type recordDiagnostic struct {
	Status string `xml:"status,attr"`
	DOI    string `xml:"doi"`
	Msg    string `xml:"msg"`
}

type batchData struct {
	RecordCount  int `xml:"record_count"`
	SuccessCount int `xml:"success_count"`
	WarningCount int `xml:"warning_count"`
	FailureCount int `xml:"failure_count"`
}
