package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"
)

func main() {

	// start tcp server
	fmt.Println("start tcp server ===>")
	listener, _ := net.Listen("tcp", ":5555")
	defer listener.Close()

	// tunneling by ngrok
	fmt.Println("tunneling by ngrok ===>")
	exec.Command(
		path.Join(os.Getenv("LAMBDA_TASK_ROOT"), "ngrok"),
		"tcp",
		"5555",
		"--authtoken",
		os.Getenv("NGROK_AUTH_TOKEN")).Start()

	for {
		// get the next event
		fmt.Println("get the next event ===>")
		resp, err := http.Get(fmt.Sprintf(
			"http://%s/2018-06-01/runtime/invocation/next",
			os.Getenv("AWS_LAMBDA_RUNTIME_API"),
		))
		if err != nil {
			panic(err)
		}

		// read the event data
		fmt.Println("read the event data ===>")
		fmt.Println(resp.Header)
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))

		// connection open
		fmt.Println("connection open ===>")
		connection, _ := listener.Accept()
		defer connection.Close()

		// bash start
		fmt.Println("bash start ===>")
		shell := exec.Command("/bin/bash", "-i")
		shell.Stdin = bufio.NewReader(connection)
		shell.Stdout = bufio.NewWriter(connection)
		shell.Stderr = bufio.NewWriter(connection)
		if err = shell.Start(); err != nil {
			panic(err)
		}

		// keep the connection until the end file is created
		for {
			time.Sleep(1 * time.Second)
			if _, err := os.Stat("/tmp/exit"); !os.IsNotExist(err) {
				if err := os.Remove("/tmp/exit"); err != nil {
					panic(err)
				}
				break
			}
		}
		fmt.Println("close connection ===>")
		connection.Close()

		// post the result
		fmt.Println("post the result ===>")
		resp, err = http.Post(fmt.Sprintf(
			"http://%s/2018-06-01/runtime/invocation/%s/response",
			os.Getenv("AWS_LAMBDA_RUNTIME_API"),
			resp.Header["Lambda-Runtime-Aws-Request-Id"][0],
		), "application/json", strings.NewReader("{\"statusCode\": 200, \"body\": \"ok\"}"))
		if err != nil {
			panic(err)
		}
		// read the event data
		fmt.Println("get the result ===>")
		body, _ = ioutil.ReadAll(resp.Body)
		fmt.Println(string(body))
	}

}
