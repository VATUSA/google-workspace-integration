package main

import (
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/integration"
)

func main() {
	err := database.Connect()
	if err != nil {
		return
	}
	err = database.MigrateDB()
	if err != nil {
		return
	}

	err = integration.ProcessPasswordEmails()
	if err != nil {
		println(err.Error())
		return
	}
}
