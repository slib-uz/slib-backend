package enum

import (
	"database/sql/driver"
	"fmt"
)

type AcademicDegreeEnum struct {
	Code  string
	Label string
}

var (
	AcademicDegreeSecondarySpecialized  = AcademicDegreeEnum{"O", "Oʻrta yoki oʻrta maxsus"}
	AcademicDegreeProfessionalEducation = AcademicDegreeEnum{"P", "Professional taʼlim"}
	AcademicDegreeBachelor              = AcademicDegreeEnum{"B", "Bakalavr"}
	AcademicDegreeMaster                = AcademicDegreeEnum{"M", "Magistr"}
	AcademicDegreePhilosophyDoctor      = AcademicDegreeEnum{"F", "Falsafa doktori"}
	AcademicDegreeScienceDoctor         = AcademicDegreeEnum{"D", "Fan Doktori"}
)

var AllAcademicDegrees = []AcademicDegreeEnum{
	AcademicDegreeSecondarySpecialized,
	AcademicDegreeProfessionalEducation,
	AcademicDegreeBachelor,
	AcademicDegreeMaster,
	AcademicDegreePhilosophyDoctor,
	AcademicDegreeScienceDoctor,
}

func GetDegreeByCode(code string) (*AcademicDegreeEnum, bool) {
	for _, d := range AllAcademicDegrees {
		if d.Code == code {
			return &d, true
		}
	}
	return nil, false
}

type AcademicDegreeCode string

func (c AcademicDegreeCode) Value() (driver.Value, error) {
	return string(c), nil
}

func (c *AcademicDegreeCode) Scan(value interface{}) error {
	str, ok := value.(string)
	if value == nil {
		*c = AcademicDegreeCode("")
		return nil
	}
	if !ok {
		return fmt.Errorf("invalid academic degree code type")
	}
	if str == "" {
		*c = AcademicDegreeCode("")
		return nil
	}
	if _, exists := GetDegreeByCode(str); !exists {
		return fmt.Errorf("unknown academic degree code: %s", str)
	}

	*c = AcademicDegreeCode(str)
	return nil
}

func (c AcademicDegreeCode) MarshalJSON() ([]byte, error) {
	if c == "" {
		return []byte("null"), nil
	}
	degree, exists := GetDegreeByCode(string(c))
	if !exists {
		return []byte("null"), nil
	}
	return fmt.Appendf(nil, `"%s"`, degree.Label), nil
}
