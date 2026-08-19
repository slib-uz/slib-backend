package enum

import (
	"database/sql/driver"
	"fmt"
)

type CityEnum struct {
	Code  string
	Label string
}

var (
	CityAndijon         = CityEnum{"AN", "Andijon viloyati"}
	CityBuxoro          = CityEnum{"BX", "Buxoro viloyati"}
	CityFargona         = CityEnum{"FR", "Fargʻona viloyati"}
	CityJizzax          = CityEnum{"JZ", "Jizzax viloyati"}
	CityXorazm          = CityEnum{"XR", "Xorazm viloyati"}
	CityNamangan        = CityEnum{"NM", "Namangan viloyati"}
	CityNavoiy          = CityEnum{"NV", "Navoiy viloyati"}
	CityQashqadaryo     = CityEnum{"QD", "Qashqadaryo viloyati"}
	CityQoraqalpogiston = CityEnum{"QR", "Qoraqalpogʻiston Respublikasi"}
	CitySamarqand       = CityEnum{"SN", "Samarqand viloyati"}
	CitySirdaryo        = CityEnum{"SR", "Sirdaryo viloyati"}
	CitySurxondaryo     = CityEnum{"SD", "Surxondaryo viloyati"}
	CityToshkent        = CityEnum{"TV", "Toshkent viloyati"}
	CityToshkentShahri  = CityEnum{"TN", "Toshkent shahri"}
	CityForeignCitizen  = CityEnum{"FC", "Chet el fuqarosi"}
)

var AllCities = []CityEnum{
	CityAndijon,
	CityBuxoro,
	CityFargona,
	CityJizzax,
	CityXorazm,
	CityNamangan,
	CityNavoiy,
	CityQashqadaryo,
	CityQoraqalpogiston,
	CitySamarqand,
	CitySirdaryo,
	CitySurxondaryo,
	CityToshkent,
	CityToshkentShahri,
	CityForeignCitizen,
}

func GetCityByCode(code string) (*CityEnum, bool) {
	for _, c := range AllCities {
		if c.Code == code {
			return &c, true
		}
	}
	return nil, false
}

// CityCode represents a city code in Uzbekistan
// @Description City code enum type
// swagger:type string
type CityCode string

func (c CityCode) Value() (driver.Value, error) {
	return string(c), nil
}

func (c *CityCode) Scan(value interface{}) error {
	if value == nil {
		*c = CityCode("")
		return nil
	}
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid city code type")
	}

	if str == "" {
		*c = CityCode("")
		return nil
	}
	if _, exists := GetCityByCode(str); !exists {
		return fmt.Errorf("unknown city code: %s", str)
	}

	*c = CityCode(str)
	return nil
}

func (c CityCode) MarshalJSON() ([]byte, error) {
	if c == "" {
		return []byte("null"), nil
	}
	city, exists := GetCityByCode(string(c))
	if !exists {
		return []byte("null"), nil
	}
	return fmt.Appendf(nil, `"%s"`, city.Label), nil
}
