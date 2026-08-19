package enum

type UserRole int

const (
	RoleAdmin            UserRole = 10 // Super Admin
	RoleExpert           UserRole = 20 // Journal Expert
	RolePublisherAdmin   UserRole = 30 // Publisher Admin
	RoleChiefEditor      UserRole = 40 // Chief Editor of a Journal
	RoleSecretary        UserRole = 50 // Journal Secretary
	RoleAccountant       UserRole = 60 // Buxgalter
	RoleAuditor          UserRole = 70 // Auditor
	RoleInstitutionAdmin UserRole = 80 // Institution Admin
)
