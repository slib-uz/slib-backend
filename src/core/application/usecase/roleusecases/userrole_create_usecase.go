package roleusecases

import (
	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/repository"
)

type UserRoleCreateUseCase struct {
	repository            repository.UserRoleRepository
	journalRepository     repository.JournalRepository
	institutionRepository repository.InstitutionRepository
}

// @inject
func NewUserRoleCreateUseCase(repository repository.UserRoleRepository, journalRepository repository.JournalRepository, institutionRepository repository.InstitutionRepository) *UserRoleCreateUseCase {
	return &UserRoleCreateUseCase{repository: repository, journalRepository: journalRepository, institutionRepository: institutionRepository}
}

func (this *UserRoleCreateUseCase) Execute(creator *entity2.UserBasicEntity, newRole *entity2.UserRoleEntity) error {

	if err := this.checkCreatorPermission(creator, newRole); err != nil {
		return err
	}

	if err := this.validateRole(newRole); err != nil {
		return err
	}

	exists, err := this.repository.Exists(newRole.UserID, newRole.PublisherID, newRole.JournalID, newRole.InstitutionID)
	if err != nil {
		return err
	}
	if exists {
		return response.NewOptionalResponse(400, response.CodeRoleAlreadyExists, nil, "")
	}

	roleEntity := newRole
	if err := this.repository.Create(roleEntity); err != nil {
		return err
	}
	return nil
}

func (this *UserRoleCreateUseCase) validateRole(newRole *entity2.UserRoleEntity) error {
	if newRole.Role == enum.RoleInstitutionAdmin && newRole.InstitutionID == nil {
		return response.InvalidArgument
	}

	if newRole.InstitutionID != nil {
		if _, err := this.institutionRepository.GetByID(*newRole.InstitutionID); err != nil {
			return err
		}
	}

	return nil
}

func (this *UserRoleCreateUseCase) checkCreatorPermission(creator *entity2.UserBasicEntity, newRole *entity2.UserRoleEntity) error {

	access := false
	for _, c := range creator.Roles {

		if c.Role == enum.RoleAdmin {
			access = true
			break
		} else if c.Role == enum.RoleChiefEditor && c.JournalID != nil && newRole.JournalID != nil && *c.JournalID == *newRole.JournalID && newRole.Role == enum.RoleSecretary {
			access = true
			break
		} else if c.Role == enum.RolePublisherAdmin && c.PublisherID != nil && newRole.JournalID != nil {
			journal, err := this.journalRepository.FindByID(*newRole.JournalID)
			if err != nil {
				return err
			}
			if *c.PublisherID == journal.PublisherID {
				access = true
				break
			}

		}
	}
	if !access {
		return response.PermissionDeniedError
	}
	return nil

}
