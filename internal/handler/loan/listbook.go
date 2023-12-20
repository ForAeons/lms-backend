package loanhandler

import (
	"lms-backend/internal/api"
	"lms-backend/internal/dataaccess/book"
	"lms-backend/internal/database"
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/loanpolicy"
	"lms-backend/internal/view/bookview"
	collection "lms-backend/pkg/collectionquery"
	"lms-backend/pkg/error/internalerror"

	"github.com/gofiber/fiber/v2"
)

const (
	listBookLoanAction = "list book on loan"
)

// @Summary List books on loan by current user
// @Description Lists books on loan by current user
// @Tags loan
// @Accept */*
// @Produce application/json
// @Success 200 {object} api.SwgResponse[bookview.BookLoanView]
// @Failure 400 {object} api.SwgErrResponse
// @Router /api/v1/loan/book [get]
func HandleListBook(c *fiber.Ctx) error {
	err := policy.Authorize(c, listBookLoanAction, loanpolicy.ReadBookPolicy())
	if err != nil {
		return err
	}

	cq := collection.GetCollectionQueryFromParam(c)
	db := database.GetDB()

	totalCount, err := book.Count(db)
	if err != nil {
		return err
	}

	dbFiltered := cq.Filter(db, book.LoanFilters())

	filteredCount, err := book.Count(dbFiltered)
	if err != nil {
		return err
	}

	dbSorted := cq.Sort(dbFiltered, book.Sorters())
	dbPaginated := cq.Paginate(dbSorted)
	books, err := book.ListWithLoan(dbPaginated)
	if err != nil {
		return err
	}

	var view = []*bookview.BookLoanView{}
	for _, b := range books {
		if (b.Loans == nil) || (len(b.Loans) == 0) {
			return internalerror.InternalServerError("book does not have any loan")
		}
		//nolint:gosec // loop does not modify struct
		view = append(view, bookview.ToBookLoanView(&b, &b.Loans[0]))
	}

	return c.JSON(api.Response{
		Data: view,
		Meta: api.Meta{
			TotalCount:    totalCount,
			FilteredCount: filteredCount,
		},
		Messages: api.Messages(
			api.SilentMessage("books listed successfully"),
		),
	})
}
