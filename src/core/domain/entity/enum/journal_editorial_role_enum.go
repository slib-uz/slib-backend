package enum

type JournalEditorialRole int

const (
	JournalEditorialRoleChiefEditor JournalEditorialRole = 1 // Bosh muharrir
	JournalEditorialRoleSecretary   JournalEditorialRole = 2 // Mas'ul kotib
	JournalEditorialRoleEditor      JournalEditorialRole = 3 // Muharrir
	JournalEditorialRoleReviewer    JournalEditorialRole = 4 // Taqrizchilar
	JournalEditorialRoleDesigner    JournalEditorialRole = 5 // Dizayner-sahifalovchi
	JournalEditorialRoleWebEditor   JournalEditorialRole = 6 // Sayt muharriri
)

type JournalEditorialRoleEnum struct {
	Code  int
	Label string
}

var AllJournalEditorialRoles = []JournalEditorialRoleEnum{
	{Code: int(JournalEditorialRoleChiefEditor), Label: "Bosh muharrir"},
	{Code: int(JournalEditorialRoleSecretary), Label: "Mas'ul kotib"},
	{Code: int(JournalEditorialRoleEditor), Label: "Muharrir"},
	{Code: int(JournalEditorialRoleReviewer), Label: "Taqrizchilar"},
	{Code: int(JournalEditorialRoleDesigner), Label: "Dizayner-sahifalovchi"},
	{Code: int(JournalEditorialRoleWebEditor), Label: "Sayt muharriri"},
}

func GetJournalEditorialRoleByCode(code int) (*JournalEditorialRoleEnum, bool) {
	for _, role := range AllJournalEditorialRoles {
		if role.Code == code {
			return &role, true
		}
	}
	return nil, false
}

func (r JournalEditorialRole) IsValid() bool {
	_, ok := GetJournalEditorialRoleByCode(int(r))
	return ok
}
