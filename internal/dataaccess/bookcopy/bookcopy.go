package bookcopy

import (
	"lms-backend/internal/dataaccess/book"
	"lms-backend/internal/dataaccess/loan"
	"lms-backend/internal/dataaccess/reservation"
	"lms-backend/internal/dataaccess/user"
	"lms-backend/internal/model"
	"lms-backend/internal/orm"
	"lms-backend/pkg/error/externalerrors"

	"gorm.io/gorm"
)

func preloadBook(db *gorm.DB) *gorm.DB {
	return db.Preload("Book")
}

func preloadAssociations(db *gorm.DB) *gorm.DB {
	return db.Scopes(preloadBook).
		Preload("Reservations").
		Preload("Loans").
		Preload("Loans.LoanHistories").
		Preload("Loans.Fines")
}

func Read(db *gorm.DB, id int64) (*model.BookCopy, error) {
	var b model.BookCopy
	result := db.Model(&model.BookCopy{}).
		Where("id = ?", id).
		First(&b)
	if err := result.Error; err != nil {
		if orm.IsRecordNotFound(err) {
			return nil, orm.ErrRecordNotFound(model.BookModelName)
		}
		return nil, err
	}

	return &b, nil
}

func ReadWithBook(db *gorm.DB, id int64) (*model.BookCopy, error) {
	var b model.BookCopy
	result := db.Model(&model.BookCopy{}).
		Scopes(preloadBook).
		Where("id = ?", id).
		First(&b)
	if err := result.Error; err != nil {
		if orm.IsRecordNotFound(err) {
			return nil, orm.ErrRecordNotFound(model.BookModelName)
		}
		return nil, err
	}

	return &b, nil
}

func ReadDetailed(db *gorm.DB, id int64) (*model.BookCopy, error) {
	var b model.BookCopy
	result := db.Model(&model.BookCopy{}).
		Scopes(preloadAssociations).
		Where("id = ?", id).
		First(&b)
	if err := result.Error; err != nil {
		if orm.IsRecordNotFound(err) {
			return nil, orm.ErrRecordNotFound(model.BookModelName)
		}
		return nil, err
	}

	return &b, nil
}

func Create(db *gorm.DB, bookID int64) (*model.BookCopy, error) {
	b := model.BookCopy{
		BookID: uint(bookID),
		Status: model.BookStatusAvailable,
	}

	if err := b.Create(db); err != nil {
		return nil, err
	}

	return ReadWithBook(db, int64(b.ID))
}

// CreateMultiple creates multiple book copies with the same book ID.
//
// This function will not preload Book
func CreateMultiple(db *gorm.DB, bookID, count int64) ([]model.BookCopy, error) {
	var bookCopies []model.BookCopy

	for i := int64(0); i < count; i++ {
		bookCopy := model.BookCopy{
			BookID: uint(bookID),
			Status: model.BookStatusAvailable,
		}
		bookCopies = append(bookCopies, bookCopy)
	}

	result := db.Create(&bookCopies)
	if err := result.Error; err != nil {
		return nil, err
	}

	return bookCopies, nil
}

func Delete(db *gorm.DB, id int64) (*model.BookCopy, error) {
	b, err := ReadWithBook(db, id)
	if err != nil {
		return nil, err
	}

	if err := b.Delete(db); err != nil {
		return nil, err
	}

	return b, nil
}

func LoanCopy(db *gorm.DB, userID, id int64) (*model.Loan, error) {
	b, err := Read(db, id)
	if err != nil {
		return nil, err
	}

	// Check if user has exceeded max loan
	hasExceededMaxLoan, err := user.HasExceededMaxLoan(db, userID)
	if err != nil {
		return nil, err
	}
	if hasExceededMaxLoan {
		return nil, externalerrors.BadRequest("You have reached the maximum number of loans")
	}

	// Check if user has loaned or reserved this book
	loanCount, err := book.CountNumberOfCopiesLoanedByUser(db, userID, int64(b.BookID))
	if err != nil {
		return nil, err
	}
	if loanCount > 0 {
		return nil, externalerrors.BadRequest("You have already loaned a copy of this book")
	}

	// Check if book is on loan
	if b.Status == model.BookStatusOnLoan {
		return nil, externalerrors.BadRequest("Book is already on loan")
	}

	// Check if book is on reserve
	if b.Status == model.BookStatusOnReserve {
		// check if the book is reserved by the same user
		var r model.Reservation
		result := db.Model(&model.Reservation{}).
			Where("user_id = ? AND book_copy_id = ? AND status = ?", userID, id, model.ReservationStatusPending).
			First(&r)
		if err := result.Error; err != nil {
			if orm.IsRecordNotFound(err) {
				return nil, externalerrors.BadRequest("Book is currently on reserve by another user")
			}
			return nil, err
		}

		// This book is reserved by the same user, so we can loan it
		// Fulfill the reservation
		r.Status = model.ReservationStatusFulfilled
		if err := r.Update(db); err != nil {
			return nil, err
		}

		// Proceed to loan
	}

	ln, err := loan.Loan(db, userID, id)
	if err != nil {
		return nil, err
	}

	// Update book status
	b.Status = model.BookStatusOnLoan
	if err := b.Update(db); err != nil {
		return nil, err
	}

	return ln, nil
}

func ReturnCopy(db *gorm.DB, loanID int64) (*model.Loan, error) {
	ln, err := loan.Read(db, loanID)
	if err != nil {
		return nil, err
	}

	b, err := Read(db, int64(ln.BookCopyID))
	if err != nil {
		return nil, err
	}

	if b.Status != model.BookStatusOnLoan {
		return nil, externalerrors.BadRequest("Book is not on loan")
	}

	ln, err = loan.ReturnLoan(db, loanID)
	if err != nil {
		return nil, err
	}

	// Update book status
	b.Status = model.BookStatusAvailable
	if err := b.Update(db); err != nil {
		return nil, err
	}

	return ln, nil
}

func ReturnByBookCopyID(db *gorm.DB, bookCopyID int64) (*model.Loan, error) {
	var ln model.Loan

	result := db.Model(&model.Loan{}).
		Where("book_copy_id = ?", bookCopyID).
		Where("status = ?", model.LoanStatusBorrowed).
		First(&ln)

	if err := result.Error; err != nil {
		if orm.IsRecordNotFound(err) {
			return nil, externalerrors.BadRequest("Book is not on loan")
		}
		return nil, result.Error
	}

	b, err := Read(db, int64(ln.BookCopyID))
	if err != nil {
		return nil, err
	}

	if b.Status != model.BookStatusOnLoan {
		return nil, externalerrors.BadRequest("Book is not on loan")
	}

	returnedLn, err := loan.ReturnLoan(db, int64(ln.ID))
	if err != nil {
		return nil, err
	}

	// Update book status
	b.Status = model.BookStatusAvailable
	if err := b.Update(db); err != nil {
		return nil, err
	}

	return returnedLn, nil
}

func RenewCopy(db *gorm.DB, loanID int64) (*model.Loan, error) {
	ln, err := loan.Read(db, loanID)
	if err != nil {
		return nil, err
	}

	b, err := Read(db, int64(ln.BookCopyID))
	if err != nil {
		return nil, err
	}

	if b.Status != model.BookStatusOnLoan {
		return nil, externalerrors.BadRequest("Book is not on loan")
	}

	renewedLn, err := loan.RenewLoan(db, loanID)
	if err != nil {
		return nil, err
	}

	// Status remain "loaned"

	return renewedLn, nil
}

func ReserveCopy(db *gorm.DB, userID, id int64) (*model.Reservation, error) {
	b, err := Read(db, id)
	if err != nil {
		return nil, err
	}

	if b.Status == model.BookStatusOnLoan {
		return nil, externalerrors.BadRequest("Book is currently on loan")
	}

	if b.Status == model.BookStatusOnReserve {
		return nil, externalerrors.BadRequest("Book is currently on reserve")
	}

	hasExceededMaxReservation, err := user.HasExceededMaxReservation(db, userID)
	if err != nil {
		return nil, err
	}
	if hasExceededMaxReservation {
		return nil, externalerrors.BadRequest("You have reached the maximum number of reservations")
	}

	loanCount, err := book.CountNumberOfCopiesLoanedByUser(db, userID, int64(b.BookID))
	if err != nil {
		return nil, err
	}
	if loanCount > 0 {
		return nil, externalerrors.BadRequest("You have already loaned a copy of this book")
	}

	resCount, err := book.CountNumberOfCopiesReservedByUser(db, userID, int64(b.BookID))
	if err != nil {
		return nil, err
	}
	if resCount > 0 {
		return nil, externalerrors.BadRequest("You have already reserved a copy of this book")
	}

	res, err := reservation.ReserveBook(db, userID, id)
	if err != nil {
		return nil, err
	}

	// Update book status
	b.Status = model.BookStatusOnReserve
	if err := b.Update(db); err != nil {
		return nil, err
	}

	return res, nil
}

// CancelReservationCopy cancels the reservation
func CancelReservationCopy(db *gorm.DB, resID int64) (*model.Reservation, error) {
	res, err := reservation.Read(db, resID)
	if err != nil {
		return nil, err
	}

	b, err := Read(db, int64(res.BookCopyID))
	if err != nil {
		return nil, err
	}

	if b.Status != model.BookStatusOnReserve {
		return nil, externalerrors.BadRequest("Book is not on reserve")
	}

	// Fulfill the reservation
	res, err = reservation.FullfilReservation(db, resID)
	if err != nil {
		return nil, err
	}

	// Update book status
	b.Status = model.BookStatusAvailable
	if err := b.Update(db); err != nil {
		return nil, err
	}

	return res, nil
}

func Count(db *gorm.DB) (int64, error) {
	var count int64

	result := orm.CloneSession(db).
		Model(&model.BookCopy{}).
		Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func List(db *gorm.DB) ([]model.BookCopy, error) {
	var bs []model.BookCopy

	result := db.Model(&model.BookCopy{}).
		Find(&bs)
	if result.Error != nil {
		return nil, result.Error
	}

	return bs, nil
}

// Preloads Book
func ListDetailed(db *gorm.DB) ([]model.BookCopy, error) {
	var bs []model.BookCopy

	result := db.Model(&model.BookCopy{}).
		Scopes(preloadBook).
		Find(&bs)
	if result.Error != nil {
		return nil, result.Error
	}

	return bs, nil
}

func GetBookTitle(db *gorm.DB, id int64) (string, error) {
	var title string

	result := db.Model(&model.BookCopy{}).
		Select("books.title").
		Joins("JOIN books ON book_copies.book_id = books.id").
		Where("book_copies.id = ?", id).
		First(&title)
	if err := result.Error; err != nil {
		if orm.IsRecordNotFound(err) {
			return "", orm.ErrRecordNotFound(model.BookModelName)
		}
		return "", err
	}

	return title, nil
}
