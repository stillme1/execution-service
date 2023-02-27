package run

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func GetFiles(problemId, submissionId string) (int, int, string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		println("Error with env file", err.Error())
		return -1, -1, "", err
	}

	// connect to database
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		println(err.Error())
		return -1, -1, "", err
	}
	defer db.Close()

	// get data using primary key: probelemId
	var data []byte
	var lang string
	var timeLimit int
	var memoryLimit int

	err = db.QueryRow("SELECT lang, timeLimit, memoryLimit, data FROM my_table WHERE problemId=$1", problemId).Scan(&lang, &timeLimit, &memoryLimit, &data)
	if err != nil {
		println(err.Error())
		return -1, -1, "", err
	}

	file, err := os.Create(submissionId + ".zip")
	if err != nil {
		println(err.Error())
		return -1, -1, "", err
	}
	defer file.Close()

	// Write the zip data to the file
	_, err = file.Write(data)
	if err != nil {
		println(err.Error())
		return -1, -1, "", err
	}
	return timeLimit, memoryLimit, lang, nil
}
