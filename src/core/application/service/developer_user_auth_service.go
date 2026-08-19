package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"slib.uz/src/core/application/response"
	"slib.uz/src/core/domain/ports/repository"
)

type DeveloperUserAuthService struct {
	repository repository.DeveloperUserRepository
}

// @inject
func NewDeveloperUserAuthService(repository repository.DeveloperUserRepository) *DeveloperUserAuthService {
	return &DeveloperUserAuthService{repository: repository}
}

func (this *DeveloperUserAuthService) Authenticate(username string, password string) (uint, error) {
	user, err := this.repository.GetByUsername(username)
	if err != nil {
		return 0, err
	}
	ok, err := this.checkPassword(password, user.Password)
	if err != nil {
		return 0, response.NewFailResponse(401, "Invalid username or password")
	}
	if !ok {
		return 0, response.NewFailResponse(401, "Invalid username or password")
	}
	return user.ID, nil
}

func (*DeveloperUserAuthService) checkPassword(password, passwordHash string) (bool, error) {
	// bcrypt
	if strings.HasPrefix(passwordHash, "$2a$") || strings.HasPrefix(passwordHash, "$2b$") || strings.HasPrefix(passwordHash, "$2y$") {
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
		return err == nil, nil
	}

	// pbkdf2_sha256
	parts := strings.Split(passwordHash, "$")
	if len(parts) != 4 {
		return false, errors.New("invalid password hash format")
	}
	algorithm := parts[0]
	iterations, err := strconv.Atoi(parts[1])
	if err != nil {
		return false, err
	}
	salt := parts[2]
	hash := parts[3]

	if algorithm != "pbkdf2_sha256" {
		return false, errors.New("unsupported algorithm: " + algorithm)
	}

	key := pbkdf2.Key([]byte(password), []byte(salt), iterations, 32, sha256.New)
	b64key := base64.StdEncoding.EncodeToString(key)
	return b64key == hash, nil
}
