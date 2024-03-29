package bookpolicy

import (
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/abilities"
	"lms-backend/internal/policy/commonpolicy"
)

func ReadPolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.AllowAll(),
	)
}

func ListPolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.AllowAll(),
	)
}

func CreatePolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(abilities.CanManageAll.Name, abilities.CanCreateBook.Name),
	)
}

func UpdatePolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(abilities.CanManageAll.Name, abilities.CanUpdateBook.Name),
	)
}

func DeletePolicy() policy.Policy {
	return commonpolicy.Any(
		commonpolicy.HasAnyAbility(abilities.CanManageAll.Name, abilities.CanDeleteBook.Name),
	)
}
