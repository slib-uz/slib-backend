package client

import (
	"crypto/tls"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"slib.uz/src/infrastructure/config"
)

// @inject
func NewMinioClient(env *config.Config) *minio.Client {
	minioClient, err := minio.New(
		env.MinioEndpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(env.MinioAccessKey, env.MinioSecretKey, ""),
			Secure: env.MinioUseSSL,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	)

	if err != nil {
		panic(err)
	}

	return minioClient

}
