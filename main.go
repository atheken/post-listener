package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
)

//go:embed landing_page.gohtml
var landing_page string

func main() {
	done := make(chan error)

	t := template.New("form")
	t, _ = t.Parse(landing_page)

	inputNames := os.Args[1:]

	renderOut := func(w io.Writer) {
		t.Execute(w, inputNames)
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "9700"
	}

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		renderOut(w)
		defer r.Body.Close()
	})

	http.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			renderOut(w)
			return
		}

		vars := []string{}

		for _, name := range inputNames {
			vars = append(vars, fmt.Sprintf("%s='%s'", name, r.FormValue(name)))
		}

		defer r.Body.Close()
		fmt.Println(strings.Join(vars, ";"))
		close(done)
	})

	server := &http.Server{Addr: fmt.Sprintf(":%s", port)}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
			os.Exit(1)
		}
	}()

	err := <-done
	server.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	} else {
		os.Exit(0)
	}
}
