package reservationview

import (
	"lms-backend/internal/model"
	"lms-backend/internal/view/sharedview"
)

type DetailedView struct {
	BaseView
	User *sharedview.UserView `json:"user"`
	Book *sharedview.BookView `json:"book"`
}

func ToDetailedView(res *model.Reservation) *DetailedView {
	return &DetailedView{
		BaseView: *ToView(res),
		User:     sharedview.ToUserView(res.User),
		Book:     sharedview.ToBookView(res.BookCopy.Book),
	}
}
