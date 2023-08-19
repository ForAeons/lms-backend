package finepolicy

import (
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/abilities"
	"lms-backend/internal/policy/commonpolicy"
)

func DeletePolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanDeleteFine.Name,
		),
	)
}

func ReadPolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanReadFine.Name,
		),
	)
}

func SettlePolicy(loanID, fineID int64) policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanSettleFine.Name,
		),
		AllowIfFineBelongsToUser(loanID, fineID),
	)
}
