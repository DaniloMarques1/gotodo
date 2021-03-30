package main

func main() {
	var app App
	app.Initialize(ConnectionDb{
		Host: "0.0.0.0",
		Port: "5433",
		User: "fitz",
		Password: "123456",
		Dbname: "todoexample",
	})

	app.Run(":8080")
}
