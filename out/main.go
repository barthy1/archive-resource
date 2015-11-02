package main

import (
	"encoding/base64"
	"encoding/json"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/concourse/archive-resource/models"
)

func main() {
	if len(os.Args) < 2 {
		println("usage: " + os.Args[0] + " <sourceDirectory>")
		os.Exit(1)
	}

	sourceDirectory := os.Args[1]

	var request models.OutRequest

	err := json.NewDecoder(os.Stdin).Decode(&request)
	if err != nil {
		fatal("reading request", err)
	}

	sourceURL, err := url.Parse(request.Source.URI)
	if err != nil {
		fatal("parsing uri", err)
	}

	authHeader := "Authorization: " + request.Source.Authorization

	directoriesToInclude := request.Params.DirectoriesToInclude
	curlPipe := exec.Command(
		"sh",
		"-c",
		`tar -C "$1" -czf - $2 | curl -H "$3" -X PUT "$3" -T -`,
		"sh",
		sourceDirectory,
		strings.Join(directoriesToInclude, " "),
		sourceURL.String(),
		authHeader,
	)

	curlPipe.Stdout = os.Stderr
	curlPipe.Stderr = os.Stderr

	err = curlPipe.Run()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	json.NewEncoder(os.Stdout).Encode(models.InResponse{})
}

func fatal(doing string, err error) {
	println("error " + doing + ": " + err.Error())
	os.Exit(1)
}

func basicAuth(user *url.Userinfo) string {
	username := user.Username()
	password, _ := user.Password()
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}
