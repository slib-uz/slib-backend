package uploadusecases_test

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"

	"slib.uz/src/core/application/response"
	"slib.uz/src/core/application/usecase/uploadusecases"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/infrastructure/config"
)

// fakeStorage faqat PostPolicyPresignedUrl ni kuzatadi; qolgan metodlar
// chaqirilsa test panic bilan yiqiladi — bu ataylab.
type fakeStorage struct {
	storage.FileStorage
	calls      int
	gotBucket  enum.Bucket
	gotObject  string
	gotType    *enum.UploadType
	gotMaxSize int64
	err        error
}

func (f *fakeStorage) PostPolicyPresignedUrl(
	_ context.Context, bucket enum.Bucket, objectName string,
	ut *enum.UploadType, maxSize int64,
) (*url.URL, map[string]string, error) {
	f.calls++
	f.gotBucket = bucket
	f.gotObject = objectName
	f.gotType = ut
	f.gotMaxSize = maxSize
	if f.err != nil {
		return nil, nil, f.err
	}
	u, _ := url.Parse("https://storage.example/slibbucketprivate")
	return u, map[string]string{"key": objectName, "Content-Type": ut.ContentType}, nil
}

func newUseCase(s *fakeStorage) *uploadusecases.UploadFileUseCase {
	return uploadusecases.NewUploadFileUseCase(s, &config.Config{UploadedFileMaxSize: 10 << 20})
}

func TestDangerousExtensionNeverReachesStorage(t *testing.T) {
	for _, name := range []string{"evil.html", "evil.svg", "shell.php", "a.exe", "noext"} {
		s := &fakeStorage{}
		_, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, name)
		if err == nil {
			t.Errorf("%s: xato kutilgandi", name)
		}
		if !errors.Is(err, response.InvalidFileError) {
			t.Errorf("%s: InvalidFileError kutilgandi, %v keldi", name, err)
		}
		if s.calls != 0 {
			t.Errorf("%s: storage chaqirilmasligi kerak edi, %d marta chaqirildi", name, s.calls)
		}
	}
}

func TestAllowedExtensionBuildsRequest(t *testing.T) {
	s := &fakeStorage{}
	got, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, "maqola.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.calls != 1 {
		t.Fatalf("storage bir marta chaqilishi kerak edi, %d", s.calls)
	}
	if s.gotType == nil || s.gotType.ContentType != "application/pdf" {
		t.Errorf("application/pdf turi uzatilishi kerak edi, %v", s.gotType)
	}
	if s.gotMaxSize != 10<<20 {
		t.Errorf("maksimal hajm uzatilishi kerak edi, %d keldi", s.gotMaxSize)
	}
	if got.Fields == nil {
		t.Error("javobda POST policy form maydonlari bo'lishi kerak edi")
	}
	if got.ObjectName != s.gotObject {
		t.Errorf("javobdagi ObjectName storage'ga uzatilgani bilan mos kelishi kerak")
	}
}

func TestClientDirectoryPathIsStripped(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, "../../../etc/passwd.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if strings.Contains(s.gotObject, "..") {
		t.Errorf("obyekt kalitida '..' qolmasligi kerak: %q", s.gotObject)
	}
	if !strings.HasPrefix(s.gotObject, string(enum.FolderArticles)+"/") {
		t.Errorf("obyekt kaliti %q papkasi bilan boshlanishi kerak: %q", enum.FolderArticles, s.gotObject)
	}
	if !strings.HasSuffix(s.gotObject, ".pdf") {
		t.Errorf("kengaytma saqlanishi kerak: %q", s.gotObject)
	}
}

func TestObjectNamesAreUnique(t *testing.T) {
	s1, s2 := &fakeStorage{}, &fakeStorage{}
	_, _ = newUseCase(s1).Execute(context.Background(), enum.FolderArticles, "a.pdf")
	_, _ = newUseCase(s2).Execute(context.Background(), enum.FolderArticles, "a.pdf")
	if s1.gotObject == s2.gotObject {
		t.Errorf("bir xil nom uchun kalitlar farq qilishi kerak edi: %q", s1.gotObject)
	}
}

func TestPublicFolderGoesToPublicBucket(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderNews, "muqova.png")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.gotBucket != enum.BucketPublic {
		t.Errorf("news papkasi PUBLIC bucket'ga tushishi kerak, %q keldi", s.gotBucket)
	}
}

func TestPrivateFolderGoesToPrivateBucket(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderApplications, "ariza.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.gotBucket != enum.BucketPrivate {
		t.Errorf("applications papkasi PRIVATE bucket'ga tushishi kerak, %q keldi", s.gotBucket)
	}
}

// Nashr fayli ommaviy: jurnal soni saytda hech qanday ruxsatsiz ochiladi.
// Ilgari u `applications` papkasiga yozilardi va PRIVATE bucket'ga tushib
// ko'rinmay qolardi.
func TestEditionsFolderGoesToPublicBucket(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderEditions, "nashr.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.gotBucket != enum.BucketPublic {
		t.Errorf("editions papkasi PUBLIC bucket'ga tushishi kerak, %q keldi", s.gotBucket)
	}
}

// Ekspert xulosasi ham ommaviy. Bu integratsiya yo'li bilan tekislaydi:
// integration_article_create_usecase.go allaqachon shu papkaga PUBLIC yozadi,
// presigned yo'l esa PRIVATE yozardi — bitta papka, ikki xil bucket edi.
func TestExpertConclusionFolderGoesToPublicBucket(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderArticleExpertConclusion, "xulosa.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.gotBucket != enum.BucketPublic {
		t.Errorf("expert_conclusion papkasi PUBLIC bucket'ga tushishi kerak, %q keldi", s.gotBucket)
	}
}

// Jurnal va OAK sertifikatlari ommaviy: ular jurnal sahifasida ko'rsatiladi.
// 16-avgustgacha mijoz ularni PUBLIC bucket'ga yozardi; papka xaritaga
// PRIVATE bo'lib tushgach ochilmay qoldi.
func TestCertificatesFolderGoesToPublicBucket(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderCertificates, "sertifikat.pdf")
	if err != nil {
		t.Fatalf("kutilmagan xato: %v", err)
	}
	if s.gotBucket != enum.BucketPublic {
		t.Errorf("certificates papkasi PUBLIC bucket'ga tushishi kerak, %q keldi", s.gotBucket)
	}
}

// Antiplagiat sertifikati va imlo tekshiruvi natijasini server o'zi yozadi va
// ikkalasini ham PUBLIC bucket'ga qo'yadi. Xarita PRIVATE deb tursa, bir kun
// mijoz shu papkaga yuklaganda fayl boshqa bucket'ga tushib, o'qish kodi uni
// topolmaydi. Xarita yozuv joyi bilan bir xil bo'lishi shart.
func TestServerWrittenFoldersMatchTheirWriteBucket(t *testing.T) {
	for _, folder := range []enum.StorageFolder{
		enum.FolderAntiPlagCertificate,
		enum.FolderSpellCheck,
	} {
		s := &fakeStorage{}
		_, err := newUseCase(s).Execute(context.Background(), folder, "hujjat.pdf")
		if err != nil {
			t.Fatalf("%s: kutilmagan xato: %v", folder, err)
		}
		if s.gotBucket != enum.BucketPublic {
			t.Errorf("%s papkasi PUBLIC bucket'ga tushishi kerak, %q keldi", folder, s.gotBucket)
		}
	}
}

func TestUnknownFolderIsRejected(t *testing.T) {
	s := &fakeStorage{}
	_, err := newUseCase(s).Execute(context.Background(), enum.StorageFolder("../secrets"), "a.pdf")
	if !errors.Is(err, response.InvalidFileError) {
		t.Errorf("InvalidFileError kutilgandi, %v keldi", err)
	}
	if s.calls != 0 {
		t.Errorf("noma'lum papka uchun storage chaqirilmasligi kerak edi")
	}
}

func TestStorageErrorIsPropagated(t *testing.T) {
	boom := errors.New("minio ishlamayapti")
	s := &fakeStorage{err: boom}
	_, err := newUseCase(s).Execute(context.Background(), enum.FolderArticles, "a.pdf")
	if !errors.Is(err, boom) {
		t.Errorf("storage xatosi uzatilishi kerak edi, %v keldi", err)
	}
}
