package loanpolicy

import (
	"lms-backend/internal/database"
	"lms-backend/internal/model"
	"lms-backend/internal/policy"
	"lms-backend/internal/session"

	"github.com/gofiber/fiber/v2"
)

type LoanBelongsToUser struct {
	LoanID int64
	BookID int64
}

func AllowIfLoanBelongsToUser(loanID, bookID int64) *LoanBelongsToUser {
	return &LoanBelongsToUser{
		LoanID: loanID,
		BookID: bookID,
	}
}

func (p *LoanBelongsToUser) Validate(c *fiber.Ctx) (policy.Decision, error) {
	userID, err := session.GetLoginSession(c)
	if err != nil {
		return policy.Deny, err
	}

	db := database.GetDB()

	var exists int64
	result := db.Model(&model.Loan{}).
		Where("id = ? AND user_id = ? AND book_id = ?", p.LoanID, userID, p.BookID).
		Count(&exists)
	if result.Error != nil {
		return policy.Deny, result.Error
	}

	if exists == 0 {
		return policy.Deny, nil
	}

	return policy.Allow, nil
}