package enum

type Bucket string

const (
	BucketPublic  = Bucket("PUBLIC")
	BucketPrivate = Bucket("PRIVATE")
)

type StorageFolder string

const (
	FolderCertificates        = StorageFolder("certificates")
	FolderAntiPlagCertificate = StorageFolder("antiplag_certificate")
	FolderApplications        = StorageFolder("applications")
	FolderArticles            = StorageFolder("articles")
	FolderNews                = StorageFolder("news")
	FolderImages              = StorageFolder("images")
	FolderSpellCheck          = StorageFolder("spellcheck")
	FolderTickets             = StorageFolder("tickets")
	FolderEditions            = StorageFolder("editions")

	FolderArticle                 = StorageFolder("article")
	FolderArticleExpertConclusion = StorageFolder("expert_conclusion")
)
