package main

import (
	"flag"
	"fmt"
	"techstacks.cn/techstacks/router"
)

var jwtSecret = flag.String("jwtSecret", "123qwe", "Input JwtSecret")

var cliFlag int

func Init() {
	flag.IntVar(&cliFlag, "flagname", 1234, "Just for demo")
}

func main() {
	Init()
	flag.Parse()
	fmt.Printf("args=%s, num=%d\n", flag.Args(), flag.NArg())
	for i := 0; i != flag.NArg(); i++ {
		fmt.Printf("arg[%d]=%s\n", i, flag.Arg(i))
	}
	fmt.Println("flagname=", cliFlag)
	fmt.Println("jwtSecret=", *jwtSecret)

	r := router.InitRouter()
	r.Run(":3000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
