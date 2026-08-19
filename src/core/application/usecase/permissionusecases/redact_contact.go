package permissionusecases

import "slib.uz/src/core/domain/entity"

// RedactUserContact so'rovchi egasi yoki privilegiyalangan (admin) bo'lmasa
// UserEntity dagi aloqa ma'lumotlarini bo'shatadi.
//
// PINFL va tug'ilgan sana bu yerda kerak emas — ular json:"-" bilan baribir
// chiqmaydi. Bu yordamchi faqat telefon va email uchun, chunki ular ba'zi
// kontekstlarda (o'z yozuvi, admin) ko'rinishi kerak.
func RedactUserContact(user *entity.UserEntity, requesterID uint, isPrivileged bool) {
	if user == nil || user.ID == requesterID || isPrivileged {
		return
	}
	user.PhoneNumber = ""
	user.Email = ""
}

// RedactBasicContact UserBasicEntity uchun aynan shu qoidani qo'llaydi.
// UserBasicEntity da email yo'q, shuning uchun faqat telefon bo'shatiladi.
func RedactBasicContact(user *entity.UserBasicEntity, requesterID uint, isPrivileged bool) {
	if user == nil || user.ID == requesterID || isPrivileged {
		return
	}
	user.PhoneNumber = ""
}

// RedactApplicationContacts ariza detali javobining uchala maxfiy yo'lini
// tozalaydi. UserEntity yo'llarida (ega, affiliation muallif) telefon+email,
// UserBasicEntity yo'lida (reviewer) telefon. PINFL/birth_date json:"-" bilan
// baribir chiqmaydi.
//
// ATAYLAB QAMROV TASHQARISIDA: ReviewStageResponseEntity.Previous (rekursiv
// oldingi bosqich, o'zining Reviewer'i bilan) va ReviewStageResponseEntity
// .Application.User (ariza egasi, ichma-ich) bu yerda tozalanmaydi. Rekursiya
// ataylab qo'shilmagan (YAGNI) — hozircha yagona chaqiruvchi
// (article_applications_usecases.ApplicationDetailUsecase, GetByIDWithRelations
// orqali) bu ikki maydonni Preload qilmaydi, shuning uchun ular doim nil keladi
// va amalda hech narsa sizmaydi. LEKIN bu latent xavf: kelajakda kimdir
// ReviewStages.Previous yoki ReviewStages.Application Preload qo'shsa, yoki
// boshqa chaqiruvchi bu maydonlarni to'ldirilgan struct bilan uzatsa, ushbu
// qorovul buni ushlamaydi va telefon/email sizib chiqadi. Preload qo'shilganda
// bu funksiya ham yangilanishi SHART (Previous ustida rekursiv chaqiruv va
// Application.User ustida RedactUserContact).
func RedactApplicationContacts(resp *entity.ApplicationResponseEntity, requesterID uint, isPrivileged bool) {
	if resp == nil {
		return
	}

	RedactUserContact(resp.User, requesterID, isPrivileged)

	if resp.Article != nil {
		for _, aff := range resp.Article.CoAuthorsWithAffiliation {
			if aff != nil && aff.Author != nil {
				RedactUserContact(aff.Author.User, requesterID, isPrivileged)
			}
		}
	}

	for _, stage := range resp.ReviewStages {
		if stage != nil {
			RedactBasicContact(stage.Reviewer, requesterID, isPrivileged)
		}
	}
}
