package main

import (
	"fmt"
	"io"
	"net/http"
)

// Mock server is used to test streaming
func main2() {
	port := "12345"
	fmt.Println("Mock server accepting connections at localhost:" + port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Println("Method not supported on this endpoint!")
			return
		}

		fmt.Println("Mock server received a connection!")

		buf := make([]byte, 8096)

		for {
			n, err := r.Body.Read(buf)
			if n > 0 {
				fmt.Print(string(buf[:n]))
			}

			if err == io.EOF {
				r.Body.Close()
				break
			}
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	})

	http.ListenAndServe(":"+port, nil)
}
