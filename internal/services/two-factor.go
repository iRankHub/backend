package services

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base32"
	"fmt"

	"github.com/pquerna/otp/totp"

	"github.com/iRankHub/backend/internal/models"
)

type TwoFactorService struct {
	queries *models.Queries
}

func NewTwoFactorService(queries *models.Queries) *TwoFactorService {
	return &TwoFactorService{queries: queries}
}

func (s *TwoFactorService) GenerateSecret() (string, error) {
	bytes := make([]byte, 20)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %v", err)
	}
	return base32.StdEncoding.EncodeToString(bytes), nil
}

func (s *TwoFactorService) EnableTwoFactor(ctx context.Context, userID int32) (string, string, error) {
    secret, err := s.GenerateSecret()
    if err != nil {
        return "", "", err
    }

    err = s.queries.UpdateUserTwoFactorSecret(ctx, models.UpdateUserTwoFactorSecretParams{
        Userid:          userID,
        TwoFactorSecret: sql.NullString{String: secret, Valid: true},
    })
    if err != nil {
        return "", "", fmt.Errorf("failed to update user two factor secret: %v", err)
    }

    key, err := totp.Generate(totp.GenerateOpts{
        Issuer:      "iRankHub",
        AccountName: fmt.Sprintf("user-%d", userID),
        Secret:      []byte(secret),  // Convert string to byte slice
    })
    if err != nil {
        return "", "", fmt.Errorf("failed to generate TOTP key: %v", err)
    }

    return secret, key.URL(), nil
}

func (s *TwoFactorService) ValidateCode(secret, code string) bool {
	return totp.Validate(code, secret)
}