package groups

import (
	"github.com/labstack/echo/v4"
	social2 "slib.uz/src/entrypoint/presentation/handlers/social"
)

//var SocialGroup = app.GetInstance().Group(
//	"/social",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//)

//func init() {
//	SocialGroup.GET("/:profile_id/detail/:id", container.InitUserSocialDetailHandler().Handle)
//	SocialGroup.POST("/:profile_id/create/", container.InitUserSocialCreateHandler().Handle)
//	SocialGroup.DELETE("/:profile_id/delete/:id", container.InitUserSocialDeleteHandler().Handle)
//	SocialGroup.PATCH("/:profile_id/update/:id", container.InitUserSocialUpdateHandler().Handle)
//	SocialGroup.PUT("/:profile_id/update/:id", container.InitUserSocialUpdateHandler().Handle)
//
//	SocialGroup.GET("/all", container.InitSocialAllHandler().Handle)
//	SocialGroup.POST("/create/", container.InitSocialInitHandler().Handle)
//	SocialGroup.PUT("/update/:id", container.InitSocialUpdateHandler().Handle)
//	SocialGroup.PATCH("/update/:id", container.InitSocialUpdateHandler().Handle)
//	SocialGroup.DELETE("/delete/:id", container.InitSocialDestroyHandler().Handle)
//}

type SocialGroup struct {
	createHandler *social2.UserSocialCreateHandler
	deleteHandler *social2.UserSocialDeleteHandler
	updateHandler *social2.UserSocialUpdateHandler
	allHandler    *social2.SocialAllHandler
}

// @inject
func NewSocialGroup(createHandler *social2.UserSocialCreateHandler, deleteHandler *social2.UserSocialDeleteHandler, updateHandler *social2.UserSocialUpdateHandler, allHandler *social2.SocialAllHandler) *SocialGroup {
	return &SocialGroup{createHandler: createHandler, deleteHandler: deleteHandler, updateHandler: updateHandler, allHandler: allHandler}
}

func (this *SocialGroup) RegisterRoutes(group *echo.Group) {
	group.POST("/create/", this.createHandler.Handle)
	group.DELETE("/delete/:id", this.deleteHandler.Handle)
	group.PATCH("/update/:id", this.updateHandler.Handle)
	group.PUT("/update/:id", this.updateHandler.Handle)

	group.GET("/all", this.allHandler.Handle)
}
