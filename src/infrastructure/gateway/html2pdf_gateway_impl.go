package gateway

import (
	"fmt"
	"io"
	"net/url"
	"strconv"

	"slib.uz/src/core/domain/entity"
	"slib.uz/src/core/domain/ports/gateway"
	env "slib.uz/src/infrastructure/config"
	"slib.uz/src/infrastructure/gateway/network"
)

type Html2PdfGatewayImpl struct {
	client  *network.CHTTpClient
	BaseURL string
}

// @inject
func NewHtml2PdfGatewayImpl(client *network.CHTTpClient, env *env.Config) gateway.Html2PdfGateway {
	return &Html2PdfGatewayImpl{client: client, BaseURL: env.Html2PdfBaseURL}
}

func (this *Html2PdfGatewayImpl) Convert(content string) (*entity.Stream, error) {

	const maxPDFBytes = 100 << 20 // 100 MiB

	// ex: http://localhost:8700/to-pdf?page-size=A4&margin-top=10mm&margin-bottom=10mm&margin-left=10mm&margin-right=10mm

	payload := map[string]string{"content": content}
	headers := map[string]string{"Accept": "application/pdf"}

	resp, err := this.client.Post(this.buildUrl(), payload, nil, headers)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[Html2PdfGatewayImpl] unexpected status code: %d", resp.StatusCode)
	}

	// Content-Length -> size
	size := int64(-1)
	if cl := resp.Header.Get("Content-Length"); cl != "" {
		if n, e := strconv.ParseInt(cl, 10, 64); e == nil {
			size = n
			if n > maxPDFBytes {
				_ = resp.Body.Close()
				return nil, fmt.Errorf("pdf too large: %d > %d", n, maxPDFBytes)
			}
		}
	}

	var rc = resp.Body
	rc = struct {
		io.Reader
		io.Closer
	}{
		Reader: io.LimitReader(resp.Body, maxPDFBytes+1),
		Closer: resp.Body,
	}

	return &entity.Stream{
		Body:        rc,
		Size:        size,
		ContentType: resp.Header.Get("Content-Type"),
	}, nil

}

func (this *Html2PdfGatewayImpl) buildUrl() string {
	u, _ := url.Parse(this.BaseURL)
	u.Path = "/to-pdf"

	q := u.Query()
	// q.Set("page-size", "A4")
	// q.Set("margin-top", "15mm")
	// q.Set("margin-bottom", "20mm")
	// q.Set("margin-left", "20mm")
	// q.Set("margin-right", "15mm")

	u.RawQuery = q.Encode()

	return u.String()
}
