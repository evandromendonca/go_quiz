package util

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"golang.org/x/term"
)

var baseUrl string = "http://localhost:8080"

func ReadPassword() string {
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return string(bytepw)
}

func PostJson(path string, object any) *http.Response {
	postBody, _ := json.Marshal(object)
	reqBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(fmt.Sprintf("%s/%s", baseUrl, path), "application/json", reqBody)
	if err != nil {
		fmt.Printf("An Error Occured %v", err)
		os.Exit(1)
	}

	return resp
}

func GetJson[T any](path string) (*http.Response, T) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseUrl, path), nil)
	if err != nil {
		fmt.Println("error building request:", err)
		os.Exit(1)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error calling GET:", err)
		os.Exit(1)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var object T

	if IsSuccessStatusCode(response.StatusCode) {
		json.Unmarshal(responseBody, &object)
	}

	return response, object
}

func GetJsonAuth[T any](path, username, password string) (*http.Response, T) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s", baseUrl, path), nil)
	if err != nil {
		fmt.Println("error building request:", err)
		os.Exit(1)
	}

	req.SetBasicAuth(username, password)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error calling GET:", err)
		os.Exit(1)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var object T

	if IsSuccessStatusCode(response.StatusCode) {
		json.Unmarshal(responseBody, &object)
	}

	return response, object
}

func PostJsonAuth[T any](path, username, password string, postObject any) (*http.Response, T) {
	var body io.Reader = nil
	if postObject != nil {
		jsonObject, _ := json.Marshal(postObject)
		body = bytes.NewBuffer(jsonObject)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s", baseUrl, path), body)
	if err != nil {
		fmt.Println("error building request:", err)
		os.Exit(1)
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("error calling POST:", err)
		os.Exit(1)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var objectResult T

	if IsSuccessStatusCode(response.StatusCode) {
		json.Unmarshal(responseBody, &objectResult)
	}

	return response, objectResult
}

func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

func EnsureAuthorized(statusCode int) {
	if statusCode == 401 {
		fmt.Println("Unauthorized access. Halt!")
		os.Exit(1)
	}
}

func GetInputNumber(prompt string, min, max int) int {
	scanner := bufio.NewScanner(os.Stdin)

	number := 0
	for {
		fmt.Printf("%s [%d..%d]:\n", prompt, min, max)
		scanner.Scan()
		var err error
		number, err = strconv.Atoi(scanner.Text())
		if err != nil || number < min || number > max {
			fmt.Printf("Number must be between %d to %d. Try again\n", min, max)
			continue
		}
		break
	}

	return number
}
