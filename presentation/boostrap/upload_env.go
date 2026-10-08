package boostrap

import (
	"fmt"
	"github.com/joho/godotenv"
)

func UploadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Could not load .env file X_x.")
	}
	fmt.Println("Env loaded successfully O_o")
}
