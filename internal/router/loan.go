package router

import (
	loanhandler "lms-backend/internal/handler/loan"

	"github.com/gofiber/fiber/v2"
)

func LoanRoutes(r fiber.Router) {
	r.Get("/", loanhandler.HandleList)
	r.Post("/", loanhandler.HandleCreate)
	r.Post("/book", loanhandler.HandleCreateByBook)

	Route(r, "/:loan_id", func(r fiber.Router) {
		r.Get("/", loanhandler.HandleRead)
		r.Patch("/return", loanhandler.HandleReturn)
		r.Patch("/renew", loanhandler.HandleRenew)
	})
}
