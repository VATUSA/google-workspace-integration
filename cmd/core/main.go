package main

import (
	"github.com/VATUSA/google-workspace-integration/internal/database"
	"github.com/VATUSA/google-workspace-integration/internal/workflow"
	"log"
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
		log.Print("Starting process loop")

		err = workflow.WorkflowMain()
		if err != nil {
			log.Printf("Error occurred in WorkflowMain: %v", err)
			return
		}

		log.Print("End of process loop -- Sleep")
		time.Sleep(5 * time.Minute)
	}
}
