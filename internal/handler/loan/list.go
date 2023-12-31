package loanhandler

import (
	"lms-backend/internal/api"
	"lms-backend/internal/dataaccess/loan"
	"lms-backend/internal/database"
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/loanpolicy"
	"lms-backend/internal/view/loanview"
	collection "lms-backend/pkg/collectionquery"

	"github.com/gofiber/fiber/v2"
)

const (
	listBookLoanAction = "list book on loan"
)

func HandleList(c *fiber.Ctx) error {
	err := policy.Authorize(c, listBookLoanAction, loanpolicy.ListPolicy())
	if err != nil {
		return err
	}

	cq := collection.GetCollectionQueryFromParam(c)
	db := database.GetDB()

	totalCount, err := loan.Count(db)
	if err != nil {
		return err
	}

	dbFiltered := cq.Filter(db, loan.Filters(), loan.JoinBookCopy, loan.JoinBook)

	filteredCount, err := loan.Count(dbFiltered)
	if err != nil {
		return err
	}

	dbSorted := cq.Sort(dbFiltered, loan.Sorters())
	dbPaginated := cq.Paginate(dbSorted)
	lns, err := loan.ListWithBookUser(dbPaginated)
	if err != nil {
		return err
	}

	var view = []loanview.DetailedView{}
	for _, l := range lns {
		//nolint:gosec // loop does not modify struct
		view = append(view, *loanview.ToDetailedView(&l))
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
