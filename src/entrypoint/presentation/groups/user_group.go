package groups

import (
	"github.com/labstack/echo/v4"
	job2 "slib.uz/src/entrypoint/presentation/handlers/job"
	profile2 "slib.uz/src/entrypoint/presentation/handlers/profile"
	"slib.uz/src/entrypoint/presentation/handlers/role"
	user2 "slib.uz/src/entrypoint/presentation/handlers/user"
)

//var UserGroup = app.GetInstance().Group(
//	"/user",
//	container.InitJwtAuthMiddleware().Call,
//	permissions.AuthenticatedPermission,
//)

//func init() {
//	UserGroup.GET("/me", container.InitUserMeHandler().Handle)
//	UserGroup.GET("/profile", container.InitUserProfileHandler().Handle)
//	UserGroup.GET("/roles", container.InitUserRoleListHandler().Handle)
//	UserGroup.GET("/profile/:id", container.InitUserProfileByIDHandler().Handle)
//	UserGroup.POST("/profile/update", container.InitProfileUpdateHandler().Handle)
//
//	UserGroup.GET("/doctorate", container.InitDoctorateHandler().Handle)
//	UserGroup.GET("/jobs", container.InitJobsHandler().Handle)
//	UserGroup.GET("/degree", container.InitDegreeHandler().Handle)
//}

type UserGroup struct {
	meHandler            *profile2.UserMeHandler
	profileHandler       *profile2.ProfileHandler
	rolesHandler         *role.UserRoleListHandler
	profileByIDHandler   *profile2.ProfileByIDHandler
	profileUpdateHandler *profile2.ProfileUpdateHandler

	doctorateHandler       *profile2.DoctoratesHandler
	jobCreateUpdateHandler *job2.JobCreateUpdateHandler
	jobDeleteHandler       *job2.JobDeleteHandler
	degreeHandler          *profile2.DegreeHandler

	articleStatisticsByYearHandler      *user2.UserArticleStatisticsByYearHandler
	articleStatisticsByYearRangeHandler *user2.UserArticleStatisticsByYearRangeHandler
	userListHandler                     *user2.UserListHandler
	userDetailHandler                   *user2.UserDetailHandler
}

// @inject
func NewUserGroup(meHandler *profile2.UserMeHandler, profileHandler *profile2.ProfileHandler, rolesHandler *role.UserRoleListHandler, profileByIDHandler *profile2.ProfileByIDHandler, profileUpdateHandler *profile2.ProfileUpdateHandler, doctorateHandler *profile2.DoctoratesHandler, jobCreateUpdateHandler *job2.JobCreateUpdateHandler, jobDeleteHandler *job2.JobDeleteHandler, degreeHandler *profile2.DegreeHandler, articleStatisticsByYearHandler *user2.UserArticleStatisticsByYearHandler, articleStatisticsByYearRangeHandler *user2.UserArticleStatisticsByYearRangeHandler, userListHandler *user2.UserListHandler, userDetailHandler *user2.UserDetailHandler) *UserGroup {
	return &UserGroup{meHandler: meHandler, profileHandler: profileHandler, rolesHandler: rolesHandler, profileByIDHandler: profileByIDHandler, profileUpdateHandler: profileUpdateHandler, doctorateHandler: doctorateHandler, jobCreateUpdateHandler: jobCreateUpdateHandler, jobDeleteHandler: jobDeleteHandler, degreeHandler: degreeHandler, articleStatisticsByYearHandler: articleStatisticsByYearHandler, articleStatisticsByYearRangeHandler: articleStatisticsByYearRangeHandler, userListHandler: userListHandler, userDetailHandler: userDetailHandler}
}

func (this *UserGroup) RegisterRoutes(group *echo.Group) {
	group.GET("/me", this.meHandler.Handle)
	group.GET("/profile", this.profileHandler.Handle)
	group.GET("/roles", this.rolesHandler.Handle)
	group.GET("/profile/:id", this.profileByIDHandler.Handle)
	group.PATCH("/profile/update", this.profileUpdateHandler.Handle)

	group.GET("/doctorate", this.doctorateHandler.Handle)
	group.POST("/job", this.jobCreateUpdateHandler.Handle)
	group.DELETE("/job", this.jobDeleteHandler.Handle)
	group.GET("/degree", this.degreeHandler.Handle)

	group.GET("/article/statistic/by-year/:year", this.articleStatisticsByYearHandler.Handle)
	group.GET("/article/statistic/by-year-range", this.articleStatisticsByYearRangeHandler.Handle)

	group.GET("/list", this.userListHandler.Handle)
	group.GET("/:id", this.userDetailHandler.Handle)
}
