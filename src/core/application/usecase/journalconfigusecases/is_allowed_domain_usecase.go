package journalconfigusecases

import (
	"strings"

	"slib.uz/src/core/domain/ports/repository"
)

type IsAllowedDomainUseCase struct {
	repository repository.JournalConfigRepository
}

// @inject
func NewIsAllowedDomainUseCase(repository repository.JournalConfigRepository) *IsAllowedDomainUseCase {
	return &IsAllowedDomainUseCase{repository: repository}
}

func (this *IsAllowedDomainUseCase) Execute(domain string) (bool, error) {
	domainVariants := this.variants(domain)
	return this.repository.ExistsByDomain(domainVariants)

}

func (this *IsAllowedDomainUseCase) variants(domain string) []string {
	host := this.normalizeAskDomain(domain)
	if host == "" {
		return nil
	}

	return []string{
		host,
		"https://" + host,
		"https://" + host + "/",
		"http://" + host,
		"http://" + host + "/",
	}
}

func (*IsAllowedDomainUseCase) normalizeAskDomain(d string) string {
	d = strings.TrimSpace(strings.ToLower(d))
	d = strings.TrimSuffix(d, ".") // trailing dot bo'lsa olib tashlaymiz
	return d
}
