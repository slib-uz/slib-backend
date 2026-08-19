package app

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	. "slib.uz/src/entrypoint/presentation/groups"
	"slib.uz/src/entrypoint/presentation/headers"
	"slib.uz/src/entrypoint/presentation/interceptor/middlewares"
	"slib.uz/src/entrypoint/presentation/interceptor/permissions"
	"slib.uz/src/infrastructure/config"
)

type App struct {
	echo *echo.Echo
	env  *config.Config

	// groups
	adminGroup                    *AdminGroup
	applicationGroup              *ArticleApplicationGroup
	articlesGroup                 *ArticlesGroup
	authGroup                     *AuthGroup
	authorGroup                   *AuthorGroup
	etaqrizGroup                  *EtaqrizGroup
	journalApplicationsGroup      *JournalApplicationsGroup
	journalGroup                  *JournalGroup
	journalStatisticsGroup        *JournalStatisticsGroup
	journalManageGroup            *JournalManageGroup
	languageGroup                 *LanguageGroup
	myApplicationGroup            *MyApplicationGroup
	publisherGroup                *PublisherGroup
	roleGroup                     *RoleGroup
	socialGroup                   *SocialGroup
	statsGroup                    *StatsGroup
	studyFieldGroup               *StudyFieldGroup
	studyFieldManageGroup         *StudyFieldManageGroup
	uploadGroup                   *UploadGroup
	userGroup                     *UserGroup
	userFindGroup                 *UserFindGroup
	developmentGroup              *DevelopmentGroup
	commentGroup                  *CommentGroup
	notificationGroup             *NotificationGroup
	organizationGroup             *OrganizationGroup
	institutionGroup              *InstitutionGroup
	notificationTestGroup         *NotificationTestGroup
	deadlineGroup                 *DeadlineGroup
	guideGroup                    *GuidesGroup
	publicOfferGroup              *PublicOfferGroup
	partnerGroup                  *PartnerGroup
	projectGroup                  *ProjectGroup
	newsGroup                     *NewsGroup
	newsCategoryGroup             *NewsCategoryGroup
	supportDialogQuestionGroup    *SupportDialogQuestionGroup
	supportDialogAnswerGroup      *SupportDialogAnswerGroup
	reportGroup                   *ReportGroup
	ratingGroup                   *JournalRatingGroup
	chiefEditorGroup              *ChiefEditorGroup
	secretaryGroup                *SecretaryGroup
	researchMetricGroup           *ResearchMetricGroup
	articleAuthorAffiliationGroup *ArticleAuthorAffiliationGroup
	siteGroup                     *SiteGroup
	userStatsGroup                *UserStatsGroup
	draftGroup                    *DraftGroup
	// roiGroup                      *RoiGroup
	aboutGroup            *AboutGroup
	soatoGroup            *SoatoGroup
	uzsciGroup            *UzsciGroup
	authorshipClaimGroup  *AuthorshipClaimGroup
	ethicPolicyGroup      *EthicsPolicyGroup
	journalConfigGroup    *JournalConfigGroup
	authV2Group           *AuthV2Group
	editionGroup          *EditionGroup
	journalNewsGroup      *JournalNewsGroup
	journalEditorialGroup *JournalEditorialGroup
	instructionGroup      *InstructionGroup
	doiSettingGroup       *DoiSettingGroup
	spellCheckGroup       *SpellCheckGroup
	antiPlagGroup         *AntiPlagGroup
	aiDetectGroup         *AiDetectGroup
	integrationGroup      *IntegrationGroup
	locationGroup         *LocationGroup
	llmToolsGroup         *LlmToolsGroup

	// middlewares
	jwtAuthMiddleware  *middlewares.JwAuthMiddleware
	developerBasicAuth *middlewares.DeveloperUserBasicAuthMiddleware
	responseMiddleware *middlewares.ResponseMiddleware
	recoveryMiddleware *middlewares.RecoveryAndLoggerMiddleware
	debugLogMiddleware *middlewares.DebugLogMiddleware
}

// @inject
func NewApp(
	echo *echo.Echo,
	env *config.Config,
	adminGroup *AdminGroup,
	applicationGroup *ArticleApplicationGroup,
	articlesGroup *ArticlesGroup,
	authGroup *AuthGroup,
	authorGroup *AuthorGroup,
	etaqrizGroup *EtaqrizGroup,
	journalApplicationsGroup *JournalApplicationsGroup,
	journalGroup *JournalGroup,
	journalStatisticsGroup *JournalStatisticsGroup,
	journalManageGroup *JournalManageGroup,
	languageGroup *LanguageGroup,
	myApplicationGroup *MyApplicationGroup,
	publisherGroup *PublisherGroup,
	roleGroup *RoleGroup,
	socialGroup *SocialGroup,
	statsGroup *StatsGroup,
	studyFieldGroup *StudyFieldGroup,
	studyFieldManageGroup *StudyFieldManageGroup,
	uploadGroup *UploadGroup,
	userGroup *UserGroup,
	userFindGroup *UserFindGroup,
	developmentGroup *DevelopmentGroup,
	commentGroup *CommentGroup,
	notificationGroup *NotificationGroup,
	organizationGroup *OrganizationGroup,
	institutionGroup *InstitutionGroup,
	notificationTestGroup *NotificationTestGroup,
	deadlineGroup *DeadlineGroup,
	guideGroup *GuidesGroup,
	publicOfferGroup *PublicOfferGroup,
	partnerGroup *PartnerGroup,
	projectGroup *ProjectGroup,
	newsGroup *NewsGroup,
	newsCategoryGroup *NewsCategoryGroup,
	supportDialogQuestionGroup *SupportDialogQuestionGroup,
	supportDialogAnswerGroup *SupportDialogAnswerGroup,
	reportGroup *ReportGroup,
	ratingGroup *JournalRatingGroup,
	chiefEditorGroup *ChiefEditorGroup,
	secretaryGroup *SecretaryGroup,
	researchMetricGroup *ResearchMetricGroup,
	articleAuthorAffiliationGroup *ArticleAuthorAffiliationGroup,
	siteGroup *SiteGroup,
	userStatsGroup *UserStatsGroup,
	draftGroup *DraftGroup,
	// roiGroup *RoiGroup,
	aboutGroup *AboutGroup,
	soatoGroup *SoatoGroup,
	uzsciGroup *UzsciGroup,
	authorshipClaimGroup *AuthorshipClaimGroup,
	ethicPolicyGroup *EthicsPolicyGroup,
	journalConfigGroup *JournalConfigGroup,
	authV2Group *AuthV2Group,
	editionGroup *EditionGroup,
	journalNewsGroup *JournalNewsGroup,
	journalEditorialGroup *JournalEditorialGroup,
	instructionGroup *InstructionGroup,
	doiSettingGroup *DoiSettingGroup,
	spellCheckGroup *SpellCheckGroup,
	antiPlagGroup *AntiPlagGroup,
	aiDetectGroup *AiDetectGroup,
	integrationGroup *IntegrationGroup,
	locationGroup *LocationGroup,
	llmToolsGroup *LlmToolsGroup,
	jwtAuthMiddleware *middlewares.JwAuthMiddleware,
	developerBasicAuth *middlewares.DeveloperUserBasicAuthMiddleware,
	responseMiddleware *middlewares.ResponseMiddleware,
	recoveryMiddleware *middlewares.RecoveryAndLoggerMiddleware,
	debugLogMiddleware *middlewares.DebugLogMiddleware,
) *App {
	return &App{
		echo:                          echo,
		env:                           env,
		adminGroup:                    adminGroup,
		applicationGroup:              applicationGroup,
		articlesGroup:                 articlesGroup,
		authGroup:                     authGroup,
		authorGroup:                   authorGroup,
		etaqrizGroup:                  etaqrizGroup,
		journalApplicationsGroup:      journalApplicationsGroup,
		journalGroup:                  journalGroup,
		journalStatisticsGroup:        journalStatisticsGroup,
		journalManageGroup:            journalManageGroup,
		languageGroup:                 languageGroup,
		myApplicationGroup:            myApplicationGroup,
		publisherGroup:                publisherGroup,
		roleGroup:                     roleGroup,
		socialGroup:                   socialGroup,
		statsGroup:                    statsGroup,
		studyFieldGroup:               studyFieldGroup,
		studyFieldManageGroup:         studyFieldManageGroup,
		uploadGroup:                   uploadGroup,
		userGroup:                     userGroup,
		userFindGroup:                 userFindGroup,
		developmentGroup:              developmentGroup,
		commentGroup:                  commentGroup,
		notificationGroup:             notificationGroup,
		organizationGroup:             organizationGroup,
		institutionGroup:              institutionGroup,
		notificationTestGroup:         notificationTestGroup,
		deadlineGroup:                 deadlineGroup,
		guideGroup:                    guideGroup,
		publicOfferGroup:              publicOfferGroup,
		partnerGroup:                  partnerGroup,
		projectGroup:                  projectGroup,
		newsGroup:                     newsGroup,
		newsCategoryGroup:             newsCategoryGroup,
		supportDialogQuestionGroup:    supportDialogQuestionGroup,
		supportDialogAnswerGroup:      supportDialogAnswerGroup,
		reportGroup:                   reportGroup,
		ratingGroup:                   ratingGroup,
		chiefEditorGroup:              chiefEditorGroup,
		secretaryGroup:                secretaryGroup,
		researchMetricGroup:           researchMetricGroup,
		articleAuthorAffiliationGroup: articleAuthorAffiliationGroup,
		siteGroup:                     siteGroup,
		userStatsGroup:                userStatsGroup,
		draftGroup:                    draftGroup,
		// roiGroup:                      roiGroup,
		aboutGroup:            aboutGroup,
		soatoGroup:            soatoGroup,
		uzsciGroup:            uzsciGroup,
		authorshipClaimGroup:  authorshipClaimGroup,
		ethicPolicyGroup:      ethicPolicyGroup,
		journalConfigGroup:    journalConfigGroup,
		authV2Group:           authV2Group,
		editionGroup:          editionGroup,
		journalNewsGroup:      journalNewsGroup,
		journalEditorialGroup: journalEditorialGroup,
		instructionGroup:      instructionGroup,
		doiSettingGroup:       doiSettingGroup,
		jwtAuthMiddleware:     jwtAuthMiddleware,
		developerBasicAuth:    developerBasicAuth,
		responseMiddleware:    responseMiddleware,
		recoveryMiddleware:    recoveryMiddleware,
		debugLogMiddleware:    debugLogMiddleware,
		spellCheckGroup:       spellCheckGroup,
		antiPlagGroup:         antiPlagGroup,
		aiDetectGroup:         aiDetectGroup,
		integrationGroup:      integrationGroup,
		locationGroup:         locationGroup,
		llmToolsGroup:         llmToolsGroup,
	}
}

func (this *App) Init() {
	this.initMiddlewares()
	this.initGroups()
}

// SecurityHeadersMiddleware barcha javoblarga xavfsizlik sarlavhalarini qo'shadi.
//
// nosniff bu yerda eng qimmatlisi: API binar javoblar ham qaytaradi
// (statistika Excel eksporti, base64 yuklab olish) va brauzer ularning
// turini taxmin qilishga urinmasligi kerak.
//
// XSSProtection "0" ATAYLAB: eski X-XSS-Protection sarlavhasi zamonaviy
// brauzerlarda o'chirilgan va yoqilganda o'zi zaiflik manbai bo'lishi mumkin.
//
// Diqqat: bu ekspertiza 2.2.8 (Clickjacking) ni YOPMAYDI — u frontend
// sahifalariga tegishli, bu backend esa ularni xizmat qilmaydi.
func SecurityHeadersMiddleware() echo.MiddlewareFunc {
	const apiCSP = "default-src 'none'; frame-ancestors 'none'"
	// Swagger UI o'z JS/CSS/inline skriptlarini yuklaydi; default-src 'none'
	// brauzerda oppoq oyna beradi.
	const swaggerDocsCSP = "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none'"

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			h := c.Response().Header()
			h.Set("X-XSS-Protection", "0")
			h.Set("X-Content-Type-Options", "nosniff")
			h.Set("X-Frame-Options", "DENY")
			if isSwaggerDocsPath(c.Request().URL.Path) {
				h.Set(echo.HeaderContentSecurityPolicy, swaggerDocsCSP)
			} else {
				h.Set(echo.HeaderContentSecurityPolicy, apiCSP)
			}
			return next(c)
		}
	}
}

func isSwaggerDocsPath(path string) bool {
	return path == "/docs" || strings.HasPrefix(path, "/docs/")
}

func (this *App) initMiddlewares() {
	this.echo.Use(this.debugLogMiddleware.Wrap)

	this.echo.Use(SecurityHeadersMiddleware())

	this.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// AllowOrigins:  this.env.CorsAllowOrigin,
		AllowMethods:  []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE, echo.OPTIONS},
		AllowHeaders:  []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, headers.HeaderAnonymToken, headers.HeaderIdempotencyKey},
		ExposeHeaders: []string{headers.HeaderAnonymToken, headers.HeaderIdempotencyKey},
	}))

	this.echo.Use(this.responseMiddleware.Call)

	this.echo.Use(this.recoveryMiddleware.Wrap)

	this.echo.Use(middlewares.ContextMiddleware)
	this.echo.Use(this.jwtAuthMiddleware.Call)

}

func (this *App) initGroups() {
	this.developmentGroup.RegisterRoutes(this.echo.Group(""), this.developerBasicAuth.Wrap)
	this.adminGroup.RegisterRoutes(this.Group("/admin", permissions.AdminPermission))
	this.applicationGroup.RegisterRoutes(this.Group("/article-application"))
	this.articlesGroup.RegisterRoutes(this.Group("/articles"))
	this.authGroup.RegisterRoutes(this.Group("/auth"))
	this.authorGroup.RegisterRoutes(this.Group("/author"))
	this.commentGroup.RegisterRoutes(this.Group("/comment"))
	this.etaqrizGroup.RegisterRoutes(this.Group("/etaqriz"))
	this.journalApplicationsGroup.RegisterRoutes(this.Group("/journal-applications", permissions.AuthenticatedPermission))
	this.journalGroup.RegisterRoutes(this.Group("/journal"))
	this.journalStatisticsGroup.RegisterRoutes(this.Group("/journal-statistics"))
	this.journalManageGroup.RegisterRoutes(this.Group("/journal-manage", permissions.AuthenticatedPermission))
	this.languageGroup.RegisterRoutes(this.Group("/language"))
	this.myApplicationGroup.RegisterRoutes(this.Group("/my-application", permissions.AuthenticatedPermission))
	this.notificationGroup.RegisterRoutes(this.Group("/notification", permissions.AuthenticatedPermission))
	this.organizationGroup.RegisterRoutes(this.Group("/organization"))
	this.institutionGroup.RegisterRoutes(this.Group("/institution"))
	this.publisherGroup.RegisterRoutes(this.Group("/publisher", permissions.AuthenticatedPermission))
	this.roleGroup.RegisterRoutes(this.Group("/role", permissions.AuthenticatedPermission))
	this.socialGroup.RegisterRoutes(this.Group("/social", permissions.AuthenticatedPermission))
	this.statsGroup.RegisterRoutes(this.Group("/stats", permissions.AuthenticatedPermission))
	this.studyFieldManageGroup.RegisterRoutes(this.Group("/studyfield/manage", permissions.AdminPermission))
	this.studyFieldGroup.RegisterRoutes(this.Group("/studyfield"))
	this.uploadGroup.RegisterRoutes(this.Group("/upload"))
	this.userFindGroup.RegisterRoutes(this.Group("/users/find", permissions.AuthenticatedPermission))
	this.userGroup.RegisterRoutes(this.Group("/user", permissions.AuthenticatedPermission))
	//this.commentGroup.RegisterRoutes(this.Group("/comment"))
	this.guideGroup.RegisterRoutes(this.Group("/guides"))
	this.publicOfferGroup.RegisterRoutes(this.Group("/public-offer"))
	this.partnerGroup.RegisterRoutes(this.Group("/partner"))
	this.projectGroup.RegisterRoutes(this.Group("/project"))
	this.newsGroup.RegisterRoutes(this.Group("/news"))
	this.newsCategoryGroup.RegisterRoutes(this.Group("/news-category"))
	this.supportDialogQuestionGroup.RegisterRoutes(this.Group("/support-dialog-question", permissions.AuthenticatedPermission))
	this.supportDialogAnswerGroup.RegisterRoutes(this.Group("/support-dialog-answer", permissions.AdminPermission))
	this.reportGroup.RegisterRoutes(this.Group("/report", permissions.AuthenticatedPermission))
	this.ratingGroup.RegisterRoutes(this.Group("/journal-rating"))
	this.chiefEditorGroup.RegisterRoutes(this.Group("/chief-editor"))
	this.secretaryGroup.RegisterRoutes(this.Group("/secretary"))
	this.researchMetricGroup.RegisterRoutes(this.Group("/research-metric", permissions.AuthenticatedPermission))
	this.notificationTestGroup.RegisterRoutes(this.Group("/notification-test", this.developerBasicAuth.Wrap))
	this.articleAuthorAffiliationGroup.RegisterRoutes(this.Group("/article-author-affiliation", permissions.AuthenticatedPermission))
	this.userStatsGroup.RegisterRoutes(this.Group("/user-stats"))

	this.deadlineGroup.RegisterRoutes(this.Group("/deadline", permissions.AuthenticatedPermission))
	this.siteGroup.RegisterRoutes(this.Group("/site"))

	this.draftGroup.RegisterRoutes(this.Group("/draft", permissions.AuthenticatedPermission))
	// this.roiGroup.RegisterRoutes(this.Group("/roi"))
	this.aboutGroup.RegisterRoutes(this.Group("/about"))
	this.soatoGroup.RegisterRoutes(this.Group("/soato"))
	this.uzsciGroup.RegisterRoutes(this.Group("/uzsci"))
	this.authorshipClaimGroup.RegisterRoutes(this.Group("/authorship-claims", permissions.AuthenticatedPermission))
	this.ethicPolicyGroup.RegisterRoutes(this.Group("/ethics_policy"))
	this.journalConfigGroup.RegisterRoutes(this.Group("/journal-config"))

	this.authV2Group.RegisterRoutes(this.Group("/auth-v2"))
	this.editionGroup.RegisterRoutes(this.Group("/edition"))
	this.journalNewsGroup.RegisterRoutes(this.Group("/journal-news"))
	this.journalEditorialGroup.RegisterRoutes(this.Group("/journal-editorial"))
	this.instructionGroup.RegisterRoutes(this.Group("/instruction"))

	this.doiSettingGroup.RegisterRoutes(this.Group("/doi-settings", permissions.AuthenticatedPermission))

	this.spellCheckGroup.RegisterRoutes(this.Group("/spellcheck"))

	this.antiPlagGroup.RegisterRoutes(this.Group("/antiplag"))

	this.aiDetectGroup.RegisterRoutes(this.Group("/ai-detect"))

	this.integrationGroup.RegisterRoutes(this.Group("/integration"))
	this.locationGroup.RegisterRoutes(this.Group("/location"))
	this.llmToolsGroup.RegisterRoutes(this.Group("/llm-tools"))
}

func (app *App) Start() {
	app.echo.Logger.Fatal(app.echo.Start(":8080"))
}

func (this *App) Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group {
	path := "/api" + prefix
	return this.echo.Group(path, m...)
}
