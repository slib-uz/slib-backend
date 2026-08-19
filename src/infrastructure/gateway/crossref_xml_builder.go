package gateway

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"slib.uz/src/core/domain/entity"
	gw "slib.uz/src/core/domain/ports/gateway"
)

// ─── CrossRef 4.4.2 XML structures ─────────────────────────────────────────

type doiBatch struct {
	XMLName        xml.Name  `xml:"doi_batch"`
	Xmlns          string    `xml:"xmlns,attr"`
	XmlnsXsi       string    `xml:"xmlns:xsi,attr"`
	Version        string    `xml:"version,attr"`
	SchemaLocation string    `xml:"xsi:schemaLocation,attr"`
	Head           head      `xml:"head"`
	Body           batchBody `xml:"body"`
}

type head struct {
	DoiBatchID string    `xml:"doi_batch_id"`
	Timestamp  string    `xml:"timestamp"`
	Depositor  depositor `xml:"depositor"`
	Registrant string    `xml:"registrant"`
}

type depositor struct {
	DepositorName string `xml:"depositor_name"`
	EmailAddress  string `xml:"email_address"`
}

type batchBody struct {
	Journal journalXML `xml:"journal"`
}

type journalXML struct {
	Metadata journalMetadata  `xml:"journal_metadata"`
	Articles []journalArticle `xml:"journal_article"`
}

type journalMetadata struct {
	Language  string `xml:"language,attr"`
	FullTitle string `xml:"full_title"`
	AbbrTitle string `xml:"abbrev_title"`
	ISSNs     []issn `xml:"issn"`
}

type issn struct {
	MediaType string `xml:"media_type,attr"`
	Value     string `xml:",chardata"`
}

type publicationDate struct {
	MediaType string `xml:"media_type,attr"`
	Year      string `xml:"year"`
}

type journalArticle struct {
	PublicationType string          `xml:"publication_type,attr"`
	Titles          titles          `xml:"titles"`
	Contributors    *contributors   `xml:"contributors,omitempty"`
	PubDate         publicationDate `xml:"publication_date"`
	DOIData         doiData         `xml:"doi_data"`
}

type titles struct {
	Title string `xml:"title"`
}

type contributors struct {
	Persons []personName `xml:"person_name"`
}

type personName struct {
	ContributorRole string `xml:"contributor_role,attr"`
	Sequence        string `xml:"sequence,attr"`
	GivenName       string `xml:"given_name"`
	Surname         string `xml:"surname"`
	Affiliation     string `xml:"affiliation,omitempty"`
	ORCID           string `xml:"ORCID,omitempty"`
}

type doiData struct {
	DOI      string `xml:"doi"`
	Resource string `xml:"resource"`
}

// ─── Builder ────────────────────────────────────────────────────────────────

func buildCrossRefXML(params *gw.CrossRefDepositParams, depositorEmail string, articleResourceURL string) ([]byte, error) {
	journal := params.Journal

	registrantName := ""
	if journal.Publisher != nil && journal.Publisher.Name != nil {
		registrantName = *journal.Publisher.Name
	}

	batch := doiBatch{
		Xmlns:          "http://www.crossref.org/schema/4.4.2",
		XmlnsXsi:       "http://www.w3.org/2001/XMLSchema-instance",
		Version:        "4.4.2",
		SchemaLocation: "http://www.crossref.org/schema/4.4.2 http://www.crossref.org/schema/deposit/crossref4.4.2.xsd",
		Head: head{
			DoiBatchID: params.BatchID,
			Timestamp:  time.Now().Format("20060102150405"),
			Depositor: depositor{
				DepositorName: registrantName,
				EmailAddress:  depositorEmail,
			},
			Registrant: registrantName,
		},
		Body: batchBody{
			Journal: buildJournal(journal, params.JournalName, params.Articles, articleResourceURL),
		},
	}

	output, err := xml.MarshalIndent(batch, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("[CrossRefXML] marshal error: %w", err)
	}

	return append([]byte(xml.Header), output...), nil
}

func buildJournal(journal *entity.JournalEntity, journalName string, articles []*entity.ArticleEntity, articleResourceURL string) journalXML {
	j := journalXML{
		Metadata: buildJournalMetadata(journal, journalName),
	}

	for _, article := range articles {
		j.Articles = append(j.Articles, buildArticle(article, articleResourceURL))
	}

	return j
}

func buildJournalMetadata(journal *entity.JournalEntity, journalName string) journalMetadata {
	fullTitle := strings.TrimSpace(journalName)
	if fullTitle == "" {
		fullTitle = getLocaleValue(journal.Name, "en")
	}

	meta := journalMetadata{
		Language:  "en",
		FullTitle: fullTitle,
		AbbrTitle: getLocaleValue(journal.ShortName, "en"),
	}

	if journal.ISSNPaper != nil && *journal.ISSNPaper != "" {
		meta.ISSNs = append(meta.ISSNs, issn{MediaType: "print", Value: *journal.ISSNPaper})
	}
	if journal.ISSNOnline != nil && *journal.ISSNOnline != "" {
		meta.ISSNs = append(meta.ISSNs, issn{MediaType: "electronic", Value: *journal.ISSNOnline})
	}

	return meta
}

func buildArticle(article *entity.ArticleEntity, articleResourceURL string) journalArticle {
	a := journalArticle{
		PublicationType: "full_text",
		Titles: titles{
			Title: getLocaleValue(article.Name, "en"),
		},
		PubDate: publicationDate{
			MediaType: "online",
			Year:      fmt.Sprintf("%d", article.PublicationDate.Year()),
		},
	}

	if article.DOI != nil && *article.DOI != "" {
		a.DOIData = doiData{
			DOI:      *article.DOI,
			Resource: fmt.Sprintf(articleResourceURL, article.ID),
		}
	}

	if c := buildContributors(article); len(c.Persons) > 0 {
		a.Contributors = &c
	}

	return a
}

func buildContributors(article *entity.ArticleEntity) contributors {
	var persons []personName

	for i, aff := range article.ArticleAuthorAffiliation {
		if aff.Author == nil {
			continue
		}

		seq := "additional"
		if i == 0 {
			seq = "first"
		}

		given, surname := splitName(aff.Author.FullName)

		p := personName{
			ContributorRole: "author",
			Sequence:        seq,
			GivenName:       given,
			Surname:         surname,
			Affiliation:     aff.OrganizationName,
		}

		if aff.Author.ORCIDID != nil && *aff.Author.ORCIDID != "" {
			p.ORCID = "https://orcid.org/" + *aff.Author.ORCIDID
		}

		persons = append(persons, p)
	}

	return contributors{Persons: persons}
}

func splitName(fullName string) (string, string) {
	parts := strings.Fields(fullName)
	if len(parts) == 0 {
		return "", ""
	}
	if len(parts) == 1 {
		return "", parts[0]
	}
	return strings.Join(parts[:len(parts)-1], " "), parts[len(parts)-1]
}

func getLocaleValue(m map[string]string, lang string) string {
	if v, ok := m[lang]; ok && v != "" {
		return v
	}
	for _, v := range m {
		if v != "" {
			return v
		}
	}
	return ""
}
