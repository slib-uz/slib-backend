package journaleditorialusecases

import (
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type JournalEditorialListByRoleUseCase struct {
	repository repository.JournalEditorialRepository
}

// @inject
func NewJournalEditorialListByRoleUseCase(repository repository.JournalEditorialRepository) *JournalEditorialListByRoleUseCase {
	return &JournalEditorialListByRoleUseCase{repository: repository}
}

type JournalEditorialListByRoleResponse struct {
	Roles []RoleWithMembers `json:"roles"`
}

type RoleWithMembers struct {
	RoleCode  int                           `json:"role_code"`
	RoleTitle string                        `json:"role_title"`
	Members   []entity.JournalEditorialEntity `json:"members"`
}

func (this *JournalEditorialListByRoleUseCase) Execute(journalID uint) (*JournalEditorialListByRoleResponse, error) {
	editorials, err := this.repository.GetAllByJournalID(journalID)
	if err != nil {
		return nil, err
	}

	roleMap := make(map[int][]entity.JournalEditorialEntity)
	for _, editorial := range editorials {
		if editorial == nil {
			continue
		}
		roleMap[editorial.RoleCode] = append(roleMap[editorial.RoleCode], *editorial)
	}

	roles := make([]RoleWithMembers, 0, len(enum.AllJournalEditorialRoles))
	for _, role := range enum.AllJournalEditorialRoles {
		members := roleMap[role.Code]
		if members == nil {
			members = []entity.JournalEditorialEntity{}
		}
		roles = append(roles, RoleWithMembers{
			RoleCode:  role.Code,
			RoleTitle: role.Label,
			Members:   members,
		})
	}

	return &JournalEditorialListByRoleResponse{Roles: roles}, nil
}
