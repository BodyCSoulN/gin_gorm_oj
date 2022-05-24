package test

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"testing"
)

func TestGenerateUuid(t *testing.T) {
	s := uuid.NewV4()

	fmt.Println(s)
}
