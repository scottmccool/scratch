package main

import "fmt"

func main() {
	fizzBuzz(100)
}

func fizzBuzz(x int) {
	for i := 1; i <= x; i++ {
		match := false
		if i%3 == 0 {
			fmt.Print("Fizz")
			match = true
		}
		if i%5 == 0 {
			fmt.Print("Buzz")
			match = true
		}
		if match == true {
			fmt.Println()
		} else {
			fmt.Println(i)
		}
	}
}
