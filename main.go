package main

func main() {
	app := NewApp()
	app.Init()

	app.Run(":8080")
}
