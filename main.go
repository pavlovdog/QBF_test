package main
import (
	"fmt"
	"app"
	"config"
)

func main() {
	// Initialize the application context
	// - Load the configuration
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)

	// Run the application
	fmt.Println("magic is happening on port 8081")

	app.Run()
}