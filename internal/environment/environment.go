package environment

import "github.com/joho/godotenv"

// Init loads in the environment
func Init(filename string) error {
	if err := godotenv.Load(filename); err != nil {
		return err;
	}
	return nil
}
