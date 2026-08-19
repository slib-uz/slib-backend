package purchase

import (
	"github.com/labstack/echo/v4"
	. "slib.uz/src/core/application/usecase/article_purchase_usecases"
	context2 "slib.uz/src/entrypoint/presentation/app/context"
)

type UserPurchasedArticlesListHandler struct {
	uc *UserPurchasedArticlesListUsecase
}

// @inject
func NewUserPurchasedArticlesListHandler(uc *UserPurchasedArticlesListUsecase) *UserPurchasedArticlesListHandler {
	return &UserPurchasedArticlesListHandler{uc: uc}
}

// Handle godoc
// @Summary      Get my purchased articles list
// @Description  Get my purchased articles list
// @Tags         purchase
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        search  query  string  false  "Search by title or author"
// @Param        page    query  int     false  "Page number" default(1)
// @Param        page_size  query  int     false  "Page size" default(10)
// @Success      200  {array}  entity.ArticlePurchaseEntity
// @Router 	 /purchase/my-purchased-articles [get]
func (this *UserPurchasedArticlesListHandler) Handle(ctx echo.Context) error {
	c := ctx.(*context2.Context)

	page, pageSize := context2.GetPagingParams(ctx)
	search := c.QueryParam("search")

	paging, err := this.uc.Execute(c.User.ID, page, pageSize, search)
	if err != nil {
		return err
	}

	return c.JsonResponse(200, paging)
}
