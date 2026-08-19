package storage

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/entity/enum"
	"slib.uz/src/core/domain/ports/storage"
	"slib.uz/src/core/utils"
	"slib.uz/src/infrastructure/config"
)

type MinioStorage struct {
	minio *minio.Client
	env   *config.Config
}

// @inject
func NewMinioStorage(minio *minio.Client, env *config.Config) storage.FileStorage {
	return &MinioStorage{minio: minio, env: env}
}

func (this *MinioStorage) Upload(folder enum.StorageFolder, localFilePath string, bucket enum.Bucket) (string, error) {

	objectPath := fmt.Sprintf("%s/%s", folder, filepath.Base(localFilePath))

	// Bu yo'l hozir o'lik: uni faqat UploadIfExists chaqiradi, u esa portda
	// izohga olingan. Shunga qaramay allow-list'ga bog'landi — kimdir uni
	// tiklasa, mime.TypeByExtension orqali allow-list chetlab o'tilardi.
	uploadType, ok := enum.LookupUploadType(strings.ToLower(filepath.Ext(localFilePath)))
	if !ok {
		return "", response.InvalidFileError
	}

	if _, err := this.minio.FPutObject(
		context.Background(),
		this.getBucketName(bucket),
		objectPath,
		localFilePath,
		minio.PutObjectOptions{
			ContentType:        uploadType.ContentType,
			ContentDisposition: uploadType.Disposition,
		},
	); err != nil {
		return "", err
	}

	return objectPath, nil

}

func (this *MinioStorage) PostPolicyPresignedUrl(
	ctx context.Context,
	bucket enum.Bucket,
	objectName string,
	ut *enum.UploadType,
	maxSize int64,
) (*url.URL, map[string]string, error) {
	policy, err := BuildPostPolicy(
		this.getBucketName(bucket),
		objectName,
		ut,
		maxSize,
		time.Now().Add(30*time.Minute),
	)
	if err != nil {
		return nil, nil, err
	}
	return this.minio.PresignedPostPolicy(ctx, policy)
}

func (this *MinioStorage) getBucketName(bucket enum.Bucket) string {
	if bucket == enum.BucketPublic {
		return this.env.PublicBucketName
	}
	return this.env.MinioBucketName
}

func (this *MinioStorage) SaveTempFile(file *multipart.FileHeader) (string, error) {

	src, err := file.Open()
	if err != nil {
		return "", response.InvalidFileError
	}

	defer utils.Closer(src)()

	if file.Size > int64(this.env.UploadedFileMaxSize) {
		return "", response.EntityToLargeError
	}

	filename := utils.MakeFileName(utils.SanitizeFilename(file.Filename), utils.RandomHex(1))

	dstPath := filepath.Join("uploads", filename)

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		panic(err)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}

	defer utils.Closer(dst)()

	if _, err = io.Copy(dst, src); err != nil {
		panic(err)
	}

	return dstPath, nil

}

func (this *MinioStorage) IsTemFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}

	return true
}

func (this *MinioStorage) UploadIfExists(folder enum.StorageFolder, localFilePath string, bucket enum.Bucket) (string, error) {
	if !this.IsTemFileExists(localFilePath) {
		return "", response.InvalidFileError
	}
	return this.Upload(folder, localFilePath, bucket)
}

func (this *MinioStorage) GetObject(bucket enum.Bucket, objectName string) ([]byte, error) {
	obj, err := this.minio.GetObject(context.Background(), this.getBucketName(bucket), objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer utils.Closer(obj)()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (this *MinioStorage) GetObjectAsBase64(bucket enum.Bucket, objectName string) (string, error) {
	obj, err := this.GetObject(bucket, objectName)
	if err != nil {
		return "", err
	}

	info, err := this.minio.StatObject(context.Background(), this.getBucketName(bucket), objectName, minio.StatObjectOptions{})
	if err != nil {
		return "", fmt.Errorf("stat object failed: %w", err)
	}

	contentType := info.ContentType
	if contentType == "" {
		contentType = http.DetectContentType(obj)
	}

	encoded := base64.StdEncoding.EncodeToString(obj)

	dataURL := fmt.Sprintf("data:%s;base64,%s", contentType, encoded)

	return dataURL, nil
}

func (this *MinioStorage) PutObject(bucket enum.Bucket, folder enum.StorageFolder, objectName string, b []byte) (string, error) {
	objectPath := fmt.Sprintf("%s/%s", folder, objectName)

	// Content-Type allow-list'dan olinadi, mime.TypeByExtension dan emas:
	// u allow-list'dan kengroq va shuning uchun xavfsizlik qarori uchun
	// yaroqsiz. Bitta yuklash yo'li boshqa qoidada qolmasligi kerak.
	contentType := "application/octet-stream"
	disposition := enum.DispositionAttachment
	if ut, ok := enum.LookupUploadType(strings.ToLower(filepath.Ext(objectName))); ok {
		contentType = ut.ContentType
		disposition = ut.Disposition
	}

	if _, err := this.minio.PutObject(
		context.Background(),
		this.getBucketName(bucket),
		objectPath,
		bytes.NewReader(b),
		int64(len(b)),
		minio.PutObjectOptions{
			ContentType:        contentType,
			ContentDisposition: disposition,
		},
	); err != nil {
		return "", err
	}

	return objectPath, nil
}

func (this *MinioStorage) PresignedURL(bucket enum.Bucket, objectName string, expires time.Duration) (string, error) {
	reqParams := ResponseParamsForObject(objectName)

	presignedURL, err := this.minio.PresignedGetObject(
		context.Background(),
		this.getBucketName(bucket),
		objectName,
		expires,
		reqParams,
	)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

func (this *MinioStorage) GetLocalFilePath(localFilePath string) ([]byte, error) {

	objectPath := fmt.Sprintf("%s/%s", "uploads", filepath.Base(localFilePath))

	if !this.IsTemFileExists(objectPath) {
		println("objectPath", objectPath)
		return nil, response.InvalidFileError
	}

	data, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (this *MinioStorage) PutStream(ctx context.Context, objectPath string, r io.Reader, size int64, contentType string, bucket enum.Bucket) error {
	_, err := this.minio.PutObject(ctx, this.getBucketName(bucket), objectPath, r, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}
