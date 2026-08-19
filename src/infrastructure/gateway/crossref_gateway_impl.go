package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	appResponse "slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/entity/enum"
	gw "slib.uz/src/core/domain/ports/gateway"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
	crossrefResp "slib.uz/src/infrastructure/gateway/response"
)

type CrossRefGatewayImpl struct {
	depositURL         string
	loginURL           string
	depositorEmail     string
	articleResourceURL string
	client             *http.Client
}

// @inject
func NewCrossRefGatewayImpl(env *config.Config, client *http.Client) gw.CrossRefGateway {
	return &CrossRefGatewayImpl{
		depositURL:         env.CrossRefDepositURL,
		loginURL:           env.CrossRefLoginURL,
		depositorEmail:     env.CrossRefDepositorEmail,
		articleResourceURL: env.CrossRefArticleResourceURL,
		client:             client,
	}
}

func (this *CrossRefGatewayImpl) CheckAuth(username, password string) (bool, error) {
	req, err := http.NewRequest("GET", this.loginURL+"?usr="+username+"&pwd="+password, nil)
	if err != nil {
		return false, fmt.Errorf("[CrossRefGateway] CheckAuth request error: %w", err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("slib.uz/1.0 (mailto:%s)", this.depositorEmail))

	resp, err := this.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("[CrossRefGateway] CheckAuth error: %w", err)
	}
	defer utils.Closer(resp.Body)()

	return resp.StatusCode == http.StatusOK, nil
}

func (this *CrossRefGatewayImpl) Deposit(params *gw.CrossRefDepositParams) (*entity.CrossRefDepositResultEntity, error) {
	xmlContent, err := buildCrossRefXML(params, this.depositorEmail, this.articleResourceURL)
	if err != nil {
		return nil, err
	}

	result := &entity.CrossRefDepositResultEntity{
		Status:      enum.DoiDepositStatusFailure,
		RequestBody: string(xmlContent),
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("operation", "doMDUpload")
	_ = writer.WriteField("usr", params.Username)
	_ = writer.WriteField("pwd", params.Password)

	part, _ := writer.CreateFormFile("mdFile", "deposit.xml")
	_, _ = part.Write(xmlContent)
	_ = writer.Close()

	req, err := http.NewRequest("POST", this.depositURL, &body)
	if err != nil {
		result.Message = err.Error()
		return result, fmt.Errorf("[CrossRefGateway] create request error: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", fmt.Sprintf("slib.uz/1.0 (mailto:%s)", this.depositorEmail))

	resp, err := this.client.Do(req)
	if err != nil {
		result.Message = err.Error()
		return result, fmt.Errorf("[CrossRefGateway] request error: %w", err)
	}
	defer utils.Closer(resp.Body)()

	respBody, _ := io.ReadAll(resp.Body)
	result.ResponseBody = string(respBody)

	if resp.StatusCode == http.StatusForbidden {
		result.Message = string(respBody)
		return result, appResponse.NewOptionalResponse(200, appResponse.CodeCrossRefForbidden, nil, result.Message)
	}
	if resp.StatusCode != http.StatusOK {
		result.Message = fmt.Sprintf("unexpected status %d", resp.StatusCode)
		return result, fmt.Errorf("[CrossRefGateway] unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	if diag, err := crossrefResp.ParseCrossRefDiagnostic(string(respBody)); err == nil {
		result.Status = diag.Status
		result.DOI = diag.DOI
		result.Message = diag.Message
		result.SubmissionID = diag.SubmissionID
	}

	return result, nil
}

func (this *CrossRefGatewayImpl) GetPrefixInfo(prefix string) (*entity.CrossRefPrefixInfoEntity, error) {
	url := fmt.Sprintf("https://api.crossref.org/prefixes/%s", prefix)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("[CrossRefGateway] GetPrefixInfo request error: %w", err)
	}
	req.Header.Set("User-Agent", fmt.Sprintf("slib.uz/1.0 (mailto:%s)", this.depositorEmail))

	resp, err := this.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[CrossRefGateway] GetPrefixInfo error: %w", err)
	}
	defer utils.Closer(resp.Body)()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[CrossRefGateway] GetPrefixInfo unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Message struct {
			Member string `json:"member"`
			Name   string `json:"name"`
			Prefix string `json:"prefix"`
		} `json:"message"`
	}
	respBody, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(respBody, &body); err != nil {
		return nil, fmt.Errorf("[CrossRefGateway] GetPrefixInfo parse error: %w", err)
	}

	return &entity.CrossRefPrefixInfoEntity{
		MemberURL: body.Message.Member,
		Name:      body.Message.Name,
		Prefix:    body.Message.Prefix,
	}, nil
}
