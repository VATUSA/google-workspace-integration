package main

import (
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/integration"
	"time"
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
	for {
		println("Starting process loop")
		err = integration.CreateFacilityGroups()
		if err != nil {
			println(err.Error())
			return
		}

		err = integration.DoProcessStaff()
		if err != nil {
			println(err.Error())
			return
		}
		println("End of process loop -- Sleep")
		time.Sleep(5 * time.Minute)
	}
}
