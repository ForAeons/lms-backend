package bookmarkhandler

import (
	"lms-backend/internal/api"
	"lms-backend/internal/dataaccess/bookmark"
	"lms-backend/internal/database"
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/bookmarkpolicy"
	"lms-backend/internal/view/bookmarkview"
	collection "lms-backend/pkg/collectionquery"

	"github.com/gofiber/fiber/v2"
)

const (
	listBookmarkAction = "list book mark"
)

func HandleList(c *fiber.Ctx) error {
	err := policy.Authorize(c, listBookmarkAction, bookmarkpolicy.ListPolicy())
	if err != nil {
		return err
	}

	cq := collection.GetCollectionQueryFromParam(c)
	db := database.GetDB()

	totalCount, err := bookmark.Count(db)
	if err != nil {
		return err
	}

	dbFiltered := cq.Filter(db, bookmark.Filters())

	filteredCount, err := bookmark.Count(dbFiltered)
	if err != nil {
		return err
	}

	dbSorted := cq.Sort(dbFiltered, bookmark.Sorters())
	dbPaginated := cq.Paginate(dbSorted)

	fns, err := bookmark.ListDetailed(dbPaginated)
	if err != nil {
		return err
	}

	var view = []bookmarkview.DetailedView{}
	for _, f := range fns {
		//nolint:gosec // loop does not modify struct
		view = append(view, *bookmarkview.ToDetailedView(&f))
	}

	return c.JSON(api.Response{
		Data: view,
		Meta: api.Meta{
			TotalCount:    totalCount,
			FilteredCount: filteredCount,
		},
		Messages: api.Messages(
			api.SilentMessage("bookmarks listed successfully"),
		),
	})
}
