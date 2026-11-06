package models

import (
	"fmt"
	"testing"

	"github.com/antonisgkamitsios/verese/internal/database"
)

func TestUser(t *testing.T) {
	db := database.SetupTestDb(t)

	model := UserModel{DB: db}

	fmt.Printf("%+v\n", model)
}
