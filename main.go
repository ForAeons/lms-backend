package main

import (
	_ "lms-backend/docs" // docs is generated by Swag CLI, you have to import it.
	"lms-backend/internal/app"
)

// @title Library Mangement System API
func main() {
	// Setup and run the app
	err := app.SetupAndRunApp()
	if err != nil {
		panic(err)
	}
}
