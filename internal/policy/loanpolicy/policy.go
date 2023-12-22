package loanpolicy

import (
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/abilities"
	"lms-backend/internal/policy/commonpolicy"
)

func ReadPolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanReadLoan.Name,
		),
	)
}

func ListPolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanReadLoan.Name,
		),
		AllowIfSelf(),
	)
}

func DeletePolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanDeleteLoan.Name,
		),
	)
}

func LoanPolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanLoanBook.Name,
		),
	)
}

func CreatePolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanLoanBook.Name,
		),
	)
}

func ReturnPolicy(loanID int64) policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanReturnBook.Name,
		),
		AllowIfLoanBelongsToUser(loanID),
	)
}

func RenewPolicy(loanID int64) policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(
			abilities.CanManageAll.Name,
			abilities.CanManageBookRecords.Name,
			abilities.CanRenewBook.Name,
		),
		AllowIfLoanBelongsToUser(loanID),
	)
}
