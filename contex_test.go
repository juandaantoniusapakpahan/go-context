package gocontext

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"
)

func TestContext(t *testing.T) {
	background := context.Background()
	fmt.Println("Background:", background)
	todo := context.TODO()
	fmt.Println("TODO:", todo)
}

func TestContextWithValue(t *testing.T) {
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextC, "g", "G")

	fmt.Println(contextA)
	fmt.Println(contextB)
	fmt.Println(contextC)
	fmt.Println(contextD)
	fmt.Println(contextE)
	fmt.Println(contextF)
	fmt.Println(contextG)
}

func TestContextGetValue(t *testing.T) {
	contextA := context.Background()
	contextB := context.WithValue(contextA, "b", "B")
	contextC := context.WithValue(contextA, "c", "C")

	contextD := context.WithValue(contextB, "d", "D")
	contextE := context.WithValue(contextB, "e", "E")

	contextF := context.WithValue(contextC, "f", "F")
	contextG := context.WithValue(contextC, "g", "G")

	contextH := context.WithValue(contextD, "h", "H")
	contextI := context.WithValue(contextF, "f", "F")

	fmt.Println(contextA.Value("b"))
	fmt.Println(contextB.Value("e"))
	fmt.Println(contextC.Value(""))
	fmt.Println(contextB.Value("g"))
	fmt.Println(contextC.Value("d"))
	fmt.Println(contextD.Value("b"))
	fmt.Println(contextE.Value("c"))
	fmt.Println(contextG.Value("c"))
	fmt.Println(contextF.Value("g"))
	fmt.Println(contextH.Value("b"))
	fmt.Println(contextI.Value("c"))

}

func CreateCounterLeak() chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1
		for {
			destination <- counter
			counter++
		}
	}()
	return destination
}

func TestGoroutineLeak(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
	destination := CreateCounterLeak() // Goroutine still run
	for value := range destination {
		fmt.Println("Counter: ", value)
		if value == 10 {
			break
		}
	}

	fmt.Println("Total Goroutine:", runtime.NumGoroutine())
}

func CreateCounterWithCancelContext(ctx context.Context) chan int {
	destination := make(chan int)

	go func() {
		defer close(destination)
		counter := 1

		for {
			select {
			case <-ctx.Done():
				return
			default:
				destination <- counter
				counter++
			}
		}
	}()

	return destination
}

func TestWithCancelContext(t *testing.T) {
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	parent := context.Background()
	ctx, cancel := context.WithCancel(parent)

	destination := CreateCounterWithCancelContext(ctx)
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

	for value := range destination {
		fmt.Println("Counter: ", value)
		if value > 7 {
			break
		}
	}
	cancel()
	time.Sleep(1 * time.Second)
	fmt.Println("Total Goroutine:", runtime.NumGoroutine())

}
