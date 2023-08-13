package bookhandler

import (
	"fmt"
	"lms-backend/internal/api"
	audit "lms-backend/internal/auditlog"
	"lms-backend/internal/dataaccess/book"
	"lms-backend/internal/database"
	"lms-backend/internal/params/bookparams"
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/bookpolicy"
	"lms-backend/internal/view/bookview"
	"lms-backend/pkg/error/externalerrors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func HandleUpdate(c *fiber.Ctx) error {
	err := policy.Authorize(c, createBookAction, bookpolicy.UpdatePolicy())
	if err != nil {
		return err
	}

	param := c.Params("id")
	bookID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return externalerrors.BadRequest(fmt.Sprintf("%s is not a valid book id.", param))
	}

	var bookParams bookparams.UpdateParams
	if err := c.BodyParser(&bookParams); err != nil {
		return err
	}

	if err := bookParams.Validate(bookID); err != nil {
		return err
	}

	db := database.GetDB()
	tx, rollBackOrCommit := audit.Begin(
		c, db, fmt.Sprintf("Updating existing book in library: %s.", bookParams.Title),
	)
	defer func() { rollBackOrCommit(err) }()

	bookModel := bookParams.ToModel()
	bookModel, err = book.Update(tx, bookModel)
	if err != nil {
		return err
	}

	view := bookview.ToView(bookModel)

	return c.JSON(api.Response{
		Data: view,
		Messages: []api.Message{
			api.SuccessMessage(fmt.Sprintf(
				"\"%s\" modified successfully.", bookModel.Title,
			))},
	})
}
