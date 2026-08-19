package gateway

import (
	"fmt"
	"time"

	core "slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
	"slib.uz/src/infrastructure/gateway/response"
)

type ETaqrizV2Gateway struct {
	client  *network.CHTTpClient
	Token   string
	Secret  string
	BaseURL string
}

// @inject
func NewETaqrizV2Gateway(client *network.CHTTpClient, env *config.Config) gateway.ETaqrizGateway {
	return &ETaqrizV2Gateway{client: client, Token: env.ETaqrizToken, Secret: env.ETaqrizSecret, BaseURL: env.ETaqrizBaseURL}
}

func (this *ETaqrizV2Gateway) GetSubmissionDetail(submissionExternalID uint) (*entity.PeerReviewSubmissionEntity, error) {
	url := fmt.Sprintf("%s/api/submissions/detail/%d", this.BaseURL, submissionExternalID)

	auth := network.NewBasicAuth(this.Token, this.Secret)

	resp, err := this.client.Get(url, auth, nil)
	if err != nil {
		return nil, fmt.Errorf("ETaqrizV2Gateway GetSubmissionDetail error while making GET request: %v", err)
	}
	if resp.StatusCode != 200 {
		body, err := network.GetBody[any](resp)
		if err != nil {
			return nil, fmt.Errorf("ETaqrizV2Gateway GetSubmissionDetail error while making GET request: %v", err)
		}
		return nil, fmt.Errorf("[ETaqrizV2Gateway] GetSubmissionDetail error: unexpected status code %d with body %v", resp.StatusCode, body)
	}

	body, err := network.GetBody[response.SubmissionDetailResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("ETaqrizV2Gateway GetSubmissionDetail error unexpected response body %v", err)
	}

	return body.ToEntity(), nil
}

func (this *ETaqrizV2Gateway) Submit(idempotencyKey string, applicationID uint, authorFullName, authorScienceID string, senderID uint, title, senderTitle string, reviewers []*entity.ReviewerEntity, reviewMethod enum.PeerReviewMethod, file string) ([]*entity.PeerReviewSubmissionEntity, error) {
	url := fmt.Sprintf("%s/api/submissions/external/create", this.BaseURL)

	headers := map[string]string{"X-Idempotency-Key": idempotencyKey}

	auth := network.NewBasicAuth(this.Token, this.Secret)

	reviewerExternalIds := make([]uint, len(reviewers))
	for i, r := range reviewers {
		reviewerExternalIds[i] = r.ExternalID
	}

	payload := map[string]any{
		"content_type":  "ARTICLE",
		"extra_data":    map[string]any{"application_id": applicationID, "author_full_name": authorFullName, "author_science_id": authorScienceID},
		"file":          file,
		"review_method": reviewMethod,
		"title":         title,
		"reviewer_ids":  reviewerExternalIds,
		"sender_title":  senderTitle,
	}

	resp, err := this.client.Post(url, payload, auth, headers)
	if err != nil {
		return nil, fmt.Errorf("ETaqrizV2Gateway Submit error while making POST request: %w", err)
	}

	if resp.StatusCode != 201 {
		body, err := network.GetBody[any](resp)
		if err != nil {
			b, _ := network.GetBody[any](resp)
			return nil, fmt.Errorf("ETaqrizV2Gateway Submit error unexpected response body %v %v", err, b)
		}
		return nil, fmt.Errorf("[ETaqrizV2Gateway] Submit error: unexpected status code %d with body %v", resp.StatusCode, body)
	}

	body, err := network.GetBody[response.SubmitResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("ETaqrizV2Gateway Submit error unexpected response body %v", err)
	}

	return body.ToEntityList(idempotencyKey, applicationID, senderID, reviewers), nil

}

func (this *ETaqrizV2Gateway) FindReviewerByScienceID(id string) (*entity.ReviewerEntity, error) {
	url := fmt.Sprintf("%s/api/users/reviewer/find?science_id=%s", this.BaseURL, id)

	resp, err := this.client.Get(url, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("ETaqrizV2Gateway FindReviewerByScienceID error while making GET request: %w", err)
	}

	if resp.StatusCode == 200 {
		reviewer, err := network.GetBody[response.ReviewerResponse](resp)
		if err != nil {
			body, _ := network.GetBody[any](resp)
			return nil, fmt.Errorf("ETaqrizV2Gateway FindReviewerByScienceID error while parsing response body: %w %v", err, body)
		}
		return reviewer.ToEntity(), nil
	}

	if resp.StatusCode == 404 {
		return nil, core.NotFoundError
	}
	errorBody, _ := network.GetBody[any](resp)
	return nil, fmt.Errorf("ETaqrizV2Gateway FindReviewerByScienceID unexpected status code %d: %v", resp.StatusCode, errorBody)
}

func (this *ETaqrizV2Gateway) ExtendDeadline(submissionExternalID uint, deadline time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (this *ETaqrizV2Gateway) Resubmit(submissionExternalID uint, fileBase64 string) (deadline time.Time, _ error) {
	url := fmt.Sprintf("%s/api/submissions/external/resubmit", this.BaseURL)

	payload := map[string]any{
		"submission_id": submissionExternalID,
		"file":          fileBase64,
	}

	resp, err := this.client.Post(url, payload, this.auth(), nil)
	if err != nil {
		return time.Time{}, fmt.Errorf("ETaqrizV2Gateway Resubmit error while making POST request: %w", err)
	}

	if resp.StatusCode != 200 {
		return time.Time{}, fmt.Errorf("ETaqrizV2Gateway Resubmit unexpected status code %d", resp.StatusCode)
	}

	body, err := network.GetBody[response.SubmissionResubmitResponse](resp)
	if err != nil {
		return time.Time{}, fmt.Errorf("ETaqrizV2Gateway Resubmit error unexpected response body %v", err)
	}

	return time.Parse(time.RFC3339, body.Deadline)

}

func (this *ETaqrizV2Gateway) auth() *network.BasicAuth {
	return network.NewBasicAuth(this.Token, this.Secret)
}
