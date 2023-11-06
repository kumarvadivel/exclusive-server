package validators

import (
	"go-server/models"
	"log"
)

func RegisterValidator(body models.User) (bool, error) {
	log.Fatal(body.Email)
	return false, nil
}
