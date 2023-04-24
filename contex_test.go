package gocontext

import (
	"context"
	"fmt"
	"testing"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println("Background:", background)
	todo := context.TODO()
	fmt.Println("TODO:", todo)
}
