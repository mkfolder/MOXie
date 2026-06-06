package auth

import (
	"github.com/pquerna/otp/totp"
)

type TOTPSetup struct {
	Secret string `json:"secret"`
	URI    string `json:"uri"`
}

func GenerateTOTPSecret(email string) (*TOTPSetup, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "MOXie",
		AccountName: email,
	})
	if err != nil {
		return nil, err
	}

	return &TOTPSetup{
		Secret: key.Secret(),
		URI:    key.URL(),
	}, nil
}

func ValidateTOTPCode(secret, code string) bool {
	return totp.Validate(code, secret)
}
