package main

func main() {

	// Inject db later
	// db := connectToDB() // returns *sql.DB
	// postRepo := postgres.NewPostgresRepo(db)

	// Call logger here and input in every layer later
	// logFile := "logs/app.log"
	// loggerObj := logger.GetLoggerObject(logFile)

	// postRepo := NewPostgresPostRepo(db, loggerObj)

	/* Make sure service is also using triple-s uploader
		import imageuploader "1337b04rd/internal/service/image_uploader"

		uploader := imageuploader.NewTripleSUploader(baseURL, "data/")

		postService := service.NewPostService(repo, uploader, logger)

	*/

}
