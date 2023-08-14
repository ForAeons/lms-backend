package abilities

import (
	"lms-backend/internal/model"
)

var (
	CanDeleteReservation model.Ability = model.Ability{
		Name:        "canDeleteReservation",
		Description: "can delete reservation",
	}
	CanCreateReservation model.Ability = model.Ability{
		Name:        "canCreateReservation",
		Description: "can create reservation",
	}
	CanCheckoutReservation model.Ability = model.Ability{
		Name:        "canCheckoutReservation",
		Description: "can checkout reservation",
	}
	CanCancelReservation model.Ability = model.Ability{
		Name:        "canCancelReservation",
		Description: "can cancel reservation",
	}
)