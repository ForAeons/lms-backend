package peopleparams

import (
	"lms-backend/internal/model"

	"github.com/gofiber/fiber/v2"
)

type CreateParams struct {
	FullName           string `json:"full_name"`
	PreferredName      string `json:"preferred_name"`
	LanguagePreference string `json:"language_preference"`
}

func (b *CreateParams) ToModel() *model.Person {
	return &model.Person{
		FullName:           b.FullName,
		PreferredName:      b.PreferredName,
		LanguagePreference: b.LanguagePreference,
	}
}

func (b *CreateParams) Validate() error {
	if b.FullName == "" {
		return fiber.NewError(fiber.StatusBadRequest, "full_name is required")
	}

	if len(b.LanguagePreference) != 2 {
		return fiber.NewError(fiber.StatusBadRequest, "language_preference must be 2 characters long")
	}

	return nil
}
