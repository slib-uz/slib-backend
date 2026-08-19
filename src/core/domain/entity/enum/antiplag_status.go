package enum

type AntiPlagStatus int

const (
	AntiPlagStatusLoading              AntiPlagStatus = 0 // Yuklanmoqda
	AntiPlagStatusLoaded               AntiPlagStatus = 1 // Yuklandi
	AntiPlagStatusChecking             AntiPlagStatus = 2 // Tekshirilmoqda
	AntiPlagStatusFailed               AntiPlagStatus = 3 // Muvafaqqiyatsiz
	AntiPlagStatusSuccess              AntiPlagStatus = 4 // Muvaffaqiyatli
	AntiPlagStatusShortReportPreparing AntiPlagStatus = 5 // Qisqa hisobot tayyorlanmoqda
	AntiPlagStatusFullReportPreparing  AntiPlagStatus = 6 // To'liq hisobot tayyorlanmoqda
	AntiPlagStatusReady                AntiPlagStatus = 7 // Tayyor
)
