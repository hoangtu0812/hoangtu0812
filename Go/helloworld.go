package main //in Go, every program is part of a package. We define this using the package keyword.
// In this example, the program belongs to the main package.
import "fmt" //lets us import files included in the fmt package. The fmt package contains functions for formatting text, including printing to the console.
//func main() {} is a function. Any code inside its curly brackets {} will be executed.
func main() {
	//fmt.Println() is a function made available from the fmt package. It is used to output/print text. In our example it will output "Hello World!".
	//fmt.Println("Hello World!") is a statement.
	//In Go, statements are separated by ending a line (hitting the Enter key) or by a semicolon ";".
	//The left curly bracket { cannot come at the start of a line.
	fmt.Println("Hello, World!")
}
