# Commander

Commander is a simple chat command parser.

### Simple Usage

````
func main() {
	c := commander.New()
	c.Use(logMiddleware)

	c.Command("/hello {times:int} {text+}", helloCommand)

	ok, err := c.Execute("/hello 5 world and everyone")
	fmt.Println("Success?", ok, "Error:", err)
}
````
