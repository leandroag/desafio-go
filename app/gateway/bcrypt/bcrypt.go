package bcrypt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type cryptService struct {
	jwtSecret []byte
}

func NewCryptService(jwtSecret []byte) *cryptService {
	return &cryptService{
		jwtSecret: jwtSecret,
	}
}

func (service cryptService) GenerateToken(accountID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": accountID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(service.jwtSecret))
}

func (service cryptService) GetAccountByToken(token string) (string, error) {
	// Extrair as informações do token
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.jwtSecret), nil
	})

	if err != nil {
		return "", err
	}

	// Obter a identificação da conta a partir das informações do token
	accountID, ok := claims["account_id"].(string)
	if !ok {
		return "", err
	}

	return accountID, nil
}

// Gera um hash do secret da conta
func (service cryptService) HashSecret(secret string) (string, error) {
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedSecret), nil
}
