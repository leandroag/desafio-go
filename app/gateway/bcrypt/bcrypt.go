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

func (service cryptService) GenerateToken(accountID int32) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id": accountID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(service.jwtSecret))
}

func (service cryptService) GetAccountByToken(token string) (int32, error) {
	// Extrair as informações do token
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	// Obter a identificação da conta a partir das informações do token
	var accountID int32 = int32(claims["account_id"].(float64))

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
