package groups

import (
	"github.com/labstack/echo/v4"
	"slib.uz/src/entrypoint/presentation/handlers/article_application"
	"slib.uz/src/entrypoint/presentation/handlers/spellcheck"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
)

type ArticleApplicationGroup struct {
	submit            *article_application.ApplicationSubmitHandler
	list              *article_application.ApplicationsListHandler
	detail            *article_application.ApplicationDetailHandler
	technicalReview   *article_application.TechnicalStageReviewHandler
	spellCheck        *spellcheck.SpellCheckHandler
	peerReviewSubmit  *article_application.SubmitForPeerReviewHandler
	peerReviewList    *article_application.PeerReviewSubmissionsListHandler
	peerReviewStage   *article_application.PeerReviewStageHandler
	publish           *article_application.ApplicationPublishHandler
	antiPlagResults   *article_application.AntiPlagResultListHandler
	antiPlagCheck     *article_application.CheckAntiPlagHandler
	aiDetectCheck     *article_application.CheckAiDetectHandler
	aiDetectResults   *article_application.AiDetectResultListHandler
	spellCheckResults *article_application.SpellCheckResultsListHandler
	peerReviewReviews *article_application.PeerReviewSubmissionReviewsHandler
	lastReviewStage   *article_application.ApplicationCheckForPaymentHandler
	confirmPayment    *article_application.ConfirmPaymentHandler

	clientBasicAuthMiddleware *middlewares.ClientBasicAuthMiddleware
}

// @inject
func NewArticleApplicationGroup(
	submit *article_application.ApplicationSubmitHandler,
	list *article_application.ApplicationsListHandler,
	detail *article_application.ApplicationDetailHandler,
	technicalReview *article_application.TechnicalStageReviewHandler,
	spellCheck *spellcheck.SpellCheckHandler,
	peerReviewSubmit *article_application.SubmitForPeerReviewHandler,
	peerReviewList *article_application.PeerReviewSubmissionsListHandler,
	peerReviewStage *article_application.PeerReviewStageHandler,
	publish *article_application.ApplicationPublishHandler,
	antiPlagResults *article_application.AntiPlagResultListHandler,
	antiPlagCheck *article_application.CheckAntiPlagHandler,
	aiDetectCheck *article_application.CheckAiDetectHandler,
	aiDetectResults *article_application.AiDetectResultListHandler,
	spellCheckResults *article_application.SpellCheckResultsListHandler,
	peerReviewReviews *article_application.PeerReviewSubmissionReviewsHandler,
	lastReviewStage *article_application.ApplicationCheckForPaymentHandler,
	confirmPayment *article_application.ConfirmPaymentHandler,
	clientBasicAuthMiddleware *middlewares.ClientBasicAuthMiddleware,
) *ArticleApplicationGroup {
	return &ArticleApplicationGroup{
		submit:                    submit,
		list:                      list,
		detail:                    detail,
		technicalReview:           technicalReview,
		spellCheck:                spellCheck,
		peerReviewSubmit:          peerReviewSubmit,
		peerReviewList:            peerReviewList,
		peerReviewStage:           peerReviewStage,
		publish:                   publish,
		antiPlagResults:           antiPlagResults,
		antiPlagCheck:             antiPlagCheck,
		aiDetectCheck:             aiDetectCheck,
		aiDetectResults:           aiDetectResults,
		spellCheckResults:         spellCheckResults,
		peerReviewReviews:         peerReviewReviews,
		lastReviewStage:           lastReviewStage,
		confirmPayment:            confirmPayment,
		clientBasicAuthMiddleware: clientBasicAuthMiddleware,
	}
}

func (this *ArticleApplicationGroup) RegisterRoutes(group *echo.Group) {
	authenticated := group.Group("", permissions.AuthenticatedPermission)

	authenticated.POST("/submit/:journalId", this.submit.Handle)
	authenticated.GET("/list", this.list.Handle)
	authenticated.GET("/detail/:applicationId", this.detail.Handle)
	authenticated.POST("/review/technical-review-stage", this.technicalReview.Handle)
	authenticated.POST("/spellcheck", this.spellCheck.Handle)
	authenticated.POST("/peer-review/submit", this.peerReviewSubmit.Handle, permissions.IdempotencyKeyRequiredPermission)
	authenticated.GET("/:applicationId/peer-review/submissions/list", this.peerReviewList.Handle)
	authenticated.GET("/peer-review/reviews/:external_id", this.peerReviewReviews.Handle)
	authenticated.POST("/review/peer-review-stage", this.peerReviewStage.Handle)
	authenticated.POST("/publish", this.publish.Handle)
	authenticated.GET("/:applicationId/antiplag/results", this.antiPlagResults.Handle)
	authenticated.POST("/:stageId/antiplag/check", this.antiPlagCheck.Handle)
	authenticated.POST("/:stageId/ai-detect/check", this.aiDetectCheck.Handle)
	authenticated.GET("/:applicationId/ai-detect/results", this.aiDetectResults.Handle)
	authenticated.GET("/:applicationId/spellcheck/results", this.spellCheckResults.Handle)

	group.GET("/:applicationId/check-for-payment", this.lastReviewStage.Handle, this.clientBasicAuthMiddleware.Wrap)
	group.POST("/:applicationId/confirm-payment", this.confirmPayment.Handle, this.clientBasicAuthMiddleware.Wrap)
}
