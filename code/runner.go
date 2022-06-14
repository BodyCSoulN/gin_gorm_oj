package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	// go run user-code/main.go
	cmd := exec.Command("go", "run", "user-code/main.go")
	var out, err bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &err
	pipe, err2 := cmd.StdinPipe()
	if err2 != nil {
		log.Fatalln(err2)
	}
	io.WriteString(pipe, "23 11\n")
	// 根据测试的输入案例进行运行得到输出结果，和标准的输出结果进行比对
	if err2 = cmd.Run(); err2 != nil {
		log.Fatalln(err2, err.String())
	}
	fmt.Println(out.String())

	println(out.String() == "34\n")
}
