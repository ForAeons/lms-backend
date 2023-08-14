package bookhandler

import (
	"fmt"
	"lms-backend/internal/api"
	audit "lms-backend/internal/auditlog"
	"lms-backend/internal/dataaccess/book"
	"lms-backend/internal/database"
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/bookpolicy"
	"lms-backend/pkg/error/externalerrors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func HandleDelete(c *fiber.Ctx) error {
	err := policy.Authorize(c, createBookAction, bookpolicy.DeletePolicy())
	if err != nil {
		return err
	}

	param := c.Params("id")
	bookID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return externalerrors.BadRequest(fmt.Sprintf("%s is not a valid book id.", param))
	}

	db := database.GetDB()
	tx, rollBackOrCommit := audit.Begin(
		c, db, fmt.Sprintf("Deleting a book in library - ID: %d", bookID),
	)
	defer func() { rollBackOrCommit(err) }()

	bookModel, err := book.Delete(tx, bookID)
	if err != nil {
		return err
	}

	return c.JSON(api.Response{
		Messages: api.Messages(
			api.SuccessMessage(fmt.Sprintf(
				"\"%s\" removed from library.", bookModel.Title,
			))),
	})
}