package loanhandler

import (
	"fmt"
	"lms-backend/internal/api"
	audit "lms-backend/internal/auditlog"
	"lms-backend/internal/dataaccess/bookcopy"
	"lms-backend/internal/dataaccess/user"
	"lms-backend/internal/database"
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/loanpolicy"
	"lms-backend/internal/session"
	"lms-backend/internal/view/loanview"
	"lms-backend/pkg/error/externalerrors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	returnBookAction = "return book"
)

func HandleReturn(c *fiber.Ctx) error {
	param := c.Params("loan_id")
	loanID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return externalerrors.BadRequest(fmt.Sprintf("%s is not a valid loan id.", param))
	}

	err = policy.Authorize(c, returnBookAction, loanpolicy.ReturnPolicy(loanID))
	if err != nil {
		return err
	}

	userID, err := session.GetLoginSession(c)
	if err != nil {
		return err
	}

	db := database.GetDB()

	username, err := user.GetUserName(db, userID)
	if err != nil {
		return err
	}

	tx, rollBackOrCommit := audit.Begin(
		c, fmt.Sprintf("%s returning loan id - \"%d\"", username, loanID),
	)
	defer func() { rollBackOrCommit(err) }()

	ln, err := bookcopy.ReturnCopy(tx, loanID)
	if err != nil {
		return err
	}

	return c.JSON(api.Response{
		Data: loanview.ToDetailedView(ln),
		Messages: api.Messages(
			api.SuccessMessage(fmt.Sprintf(
				"Loan id - \"%d\" has been returned.", loanID,
			))),
	})
}
