package util

import (
	"log"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() string {
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln(err)
	}

	return string(bytepw)
}
