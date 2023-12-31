package abilities

import (
	"lms-backend/internal/model"
)

var (
	CanReadLoan model.Ability = model.Ability{
		Name:        "canReadLoan",
		Description: "can read loan",
	}
	CanLoanBook model.Ability = model.Ability{
		Name:        "canBorrowBook",
		Description: "can borrow book",
	}
	CanReturnBook model.Ability = model.Ability{
		Name:        "canReturnBook",
		Description: "can return book",
	}
	CanRenewBook model.Ability = model.Ability{
		Name:        "canRenewBook",
		Description: "can renew book",
	}
	CanDeleteLoan model.Ability = model.Ability{
		Name:        "canDeleteLoan",
		Description: "can delete loan",
	}
)
