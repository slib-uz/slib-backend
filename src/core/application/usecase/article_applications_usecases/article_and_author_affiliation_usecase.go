package usecase

import (
	"slib.uz/src/core/application/response"
	entity2 "slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/repository"
	"slib.uz/src/core/utils"
)

type ArticleSubmissionUseCase struct {
	authorRepository                   repository.AuthorRepository
	articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository
}

// @inject
func NewArticleSubmissionUseCase(authorRepository repository.AuthorRepository, articleAuthorAffiliationRepository repository.ArticleAuthorAffiliationRepository) *ArticleSubmissionUseCase {
	return &ArticleSubmissionUseCase{authorRepository: authorRepository, articleAuthorAffiliationRepository: articleAuthorAffiliationRepository}
}

func (this *ArticleSubmissionUseCase) Execute(coAuthorIDs, updateAuthorAffiliationIDs []uint) ([]*entity2.ArticleAuthorAffiliationEntity, error) {
	// remove duplicate co-author IDs
	uniqueCoAuthorIDs := this.removeDuplicateIDs(coAuthorIDs)

	// get the author IDs from the existing article author affiliations
	existingAuthorIDs, err := this.articleAuthorAffiliationRepository.GetAuthorIdsByArticleAuthorAffiliationIDs(updateAuthorAffiliationIDs)
	if err != nil {
		return nil, err
	}

	// add the existing author IDs to the unique co-author IDs
	for _, authorID := range existingAuthorIDs {
		if !utils.In(authorID, uniqueCoAuthorIDs) {
			uniqueCoAuthorIDs = append(uniqueCoAuthorIDs, authorID)
		}
	}

	// create a map for quick search of existing authors
	existingAuthorIDMap := this.createIDMap(existingAuthorIDs)

	// define the co-authors that need a job (not in the existing affiliations)
	coAuthorIDsNeedingJobs := this.findIDsNotInMap(uniqueCoAuthorIDs, existingAuthorIDMap)

	// get the jobs for the co-authors
	jobs, err := this.authorRepository.GetJobs(coAuthorIDsNeedingJobs)
	if err != nil {
		return nil, err
	}

	// get the unique jobs (the last job for each author)
	uniqueJobsByAuthor := this.getUniqueJobsByAuthor(jobs)

	// check if all the co-authors have jobs
	if len(uniqueJobsByAuthor) != len(coAuthorIDsNeedingJobs) {
		return nil, response.NewFailResponse(400, "Some co-authors do not have jobs")
	}

	// create a map of jobs by author ID
	jobByAuthorMap := this.createJobByAuthorMap(uniqueJobsByAuthor)

	var articleAuthorAffiliations []*entity2.ArticleAuthorAffiliationEntity

	for _, coAuthorID := range uniqueCoAuthorIDs {
		authorJob := jobByAuthorMap[coAuthorID]
		if authorJob == nil {
			continue
		}
		articleAuthorAffiliation := entity2.NewArticleAuthorAffiliationEntity(
			0,
			nil,
			coAuthorID,
			authorJob.OrganizationID,
			authorJob.OrganizationName,
			authorJob.OrganizationTin,
			authorJob.PositionName,
		)
		articleAuthorAffiliations = append(articleAuthorAffiliations, articleAuthorAffiliation)
	}

	return articleAuthorAffiliations, nil
}

// removeDuplicateIDs removes duplicates from the slice of IDs
func (this *ArticleSubmissionUseCase) removeDuplicateIDs(ids []uint) []uint {
	seen := make(map[uint]bool)
	unique := make([]uint, 0, len(ids))

	for _, id := range ids {
		if !seen[id] {
			seen[id] = true
			unique = append(unique, id)
		}
	}

	return unique
}

// createIDMap creates a map for quick search of IDs
func (this *ArticleSubmissionUseCase) createIDMap(ids []uint) map[uint]bool {
	idMap := make(map[uint]bool, len(ids))
	for _, id := range ids {
		idMap[id] = true
	}
	return idMap
}

// findIDsNotInMap finds the IDs that are not in the map
func (this *ArticleSubmissionUseCase) findIDsNotInMap(ids []uint, idMap map[uint]bool) []uint {
	result := make([]uint, 0, len(ids))
	for _, id := range ids {
		if !idMap[id] {
			result = append(result, id)
		}
	}
	return result
}

// getUniqueJobsByAuthor returns the unique jobs for each author (the last job)
func (this *ArticleSubmissionUseCase) getUniqueJobsByAuthor(jobs []*entity2.JobWithAuthorIDEntity) []*entity2.JobWithAuthorIDEntity {
	if len(jobs) == 0 {
		return []*entity2.JobWithAuthorIDEntity{}
	}

	// iterate from the end to get the last job for each author
	seen := make(map[uint]bool)
	result := make([]*entity2.JobWithAuthorIDEntity, 0, len(jobs))

	for i := len(jobs) - 1; i >= 0; i-- {
		job := jobs[i]
		if !seen[job.AuthorID] {
			seen[job.AuthorID] = true
			result = append(result, job)
		}
	}

	// reverse the result to preserve the order
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

// createJobByAuthorMap creates a map of jobs by author ID
func (this *ArticleSubmissionUseCase) createJobByAuthorMap(jobs []*entity2.JobWithAuthorIDEntity) map[uint]*entity2.JobWithAuthorIDEntity {
	jobMap := make(map[uint]*entity2.JobWithAuthorIDEntity, len(jobs))
	for _, job := range jobs {
		jobMap[job.AuthorID] = job
	}
	return jobMap
}
