package gateway

import "slib.uz/src/core/domain/entity"

type Html2PdfGateway interface {
	Convert(content string) (*entity.Stream, error)
}
