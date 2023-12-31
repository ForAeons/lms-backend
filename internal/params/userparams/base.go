package userparams

import (
	"lms-backend/internal/model"
	"regexp"
)

var (
	emailReg = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

type BaseUserParams struct {
	Username string `json:"username"`
	// Email    string `json:"email"`
	Password string `json:"password"`
}

func (b *BaseUserParams) ToModel() *model.User {
	return &model.User{
		Username: b.Username,
		// Email:             b.Email,
		EncryptedPassword: b.Password,
	}
}

func (b *BaseUserParams) Validate() error {
	// if len(b.Email) > 0 && !emailReg.MatchString(b.Email) {
	// 	return fiber.NewError(fiber.StatusBadRequest, "invalid email format")
	// }

	return nil
}
