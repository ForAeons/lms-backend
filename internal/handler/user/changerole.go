package userhandler

import (
	"fmt"
	"lms-backend/internal/api"
	audit "lms-backend/internal/auditlog"
	"lms-backend/internal/dataaccess/user"
	"lms-backend/internal/database"
	"lms-backend/internal/params/userparams"
	"lms-backend/internal/policy"
	"lms-backend/internal/policy/userpolicy"
	"lms-backend/internal/view/userview"
	"lms-backend/pkg/error/externalerrors"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	changeRoleAction = "change role"
)

func HandleChangeRole(c *fiber.Ctx) error {
	param := c.Params("user_id")
	userID, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return externalerrors.BadRequest(fmt.Sprintf("%s is not a valid user id.", param))
	}

	var params userparams.UpdateRoleParams
	err = c.BodyParser(&params)
	if err != nil {
		return err
	}

	if err := params.Validate(); err != nil {
		return err
	}

	err = policy.Authorize(c, changeRoleAction, userpolicy.UpdateRolePolicy(userID, params.RoleID))
	if err != nil {
		return err
	}

	db := database.GetDB()
	username, err := user.GetUserName(db, userID)

	tx, rollBackOrCommit := audit.Begin(
		c, fmt.Sprintf("updating user %s's role to role %d", username, params.RoleID),
	)
	defer func() { rollBackOrCommit(err) }()

	usr, err := user.UpdateRoles(tx, userID, params.RoleID)
	if err != nil {
		return err
	}

	abilities, err := user.GetAbilities(tx, userID)
	if err != nil {
		return err
	}

	return c.JSON(api.Response{
		Data: userview.ToView(usr, abilities...),
		Messages: api.Messages(
			api.SuccessMessage(fmt.Sprintf(
				"Successfully updated user %s's role to role %d.", username, params.RoleID),
			),
		),
	})
}
