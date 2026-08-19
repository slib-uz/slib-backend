package storage

import (
	"context"
	"io"
	"mime/multipart"
	"net/url"
	"time"

	"slib.uz/src/core/domain/entity/enum"
)

type FileStorage interface {
	SaveTempFile(file *multipart.FileHeader) (string, error)
	//UploadIfExists(folder enum.StorageFolder, localFilePath string, bucket enum.Bucket) (string, error)
	GetObject(bucket enum.Bucket, objectName string) ([]byte, error)
	GetObjectAsBase64(bucket enum.Bucket, objectName string) (string, error)
	PutObject(bucket enum.Bucket, folder enum.StorageFolder, objectName string, b []byte) (string, error)
	PresignedURL(bucket enum.Bucket, objectName string, expires time.Duration) (string, error)
	PutStream(ctx context.Context, objectPath string, r io.Reader, size int64, contentType string, bucket enum.Bucket) error
	// PostPolicyPresignedUrl imzolangan POST policy qaytaradi.
	// Content-Type, hajm oralig'i va obyekt kaliti policy ichida qat'iy
	// belgilanadi — mijoz ularni o'zgartira olmaydi, MinIO mos kelmagan
	// yuklashni rad etadi.
	PostPolicyPresignedUrl(ctx context.Context, bucket enum.Bucket, objectName string, ut *enum.UploadType, maxSize int64) (*url.URL, map[string]string, error)
}
