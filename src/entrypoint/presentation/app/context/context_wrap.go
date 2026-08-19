package context

import (
	"time"

	"github.com/labstack/echo/v4"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity"
)

type Context struct {
	echo.Context
	User        *entity.UserBasicEntity
	Client      *entity.ClientEntity
	AnonymousID string

	// TokenID va TokenExp joriy access tokenning jti si va muddati.
	// JWT middleware to'ldiradi; logout ularsiz tokenni bekor qila olmaydi.
	TokenID  string
	TokenExp time.Time
}

func NewContext(context echo.Context) *Context {
	return &Context{Context: context}
}

// func (this *Context) super() echo.Context {
// 	return this.Context
// }

func (this *Context) JsonResponse(status int, data any) error {

	payload := response.NewResponse(status, nil, nil, "", data)

	return this.JSON(status, payload)
}

func (this *Context) SetClient(client *entity.ClientEntity) {
	this.Client = client
}

func (this *Context) GetClient() *entity.ClientEntity {
	return this.Client
}
