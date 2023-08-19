package policy

import (
	"fmt"
	"lms-backend/util/ternary"

	"github.com/gofiber/fiber/v2"
)

type Decision int

const (
	Allow Decision = iota
	Deny
)

type Policy interface {
	Validate(c *fiber.Ctx) (Decision, error)
}

func Authorize(c *fiber.Ctx, action string, policy Policy) error {
	decision, err := policy.Validate(c)
	if err != nil || decision == Deny {
		return fiber.NewError(fiber.StatusForbidden,
			ternary.If[string](err != nil).
				LazyThen(func() string { return err.Error() }).
				LazyElse(func() string { return fmt.Sprintf("You are not authorized to %s.", action) }),
		)
	}

	return nil
}
