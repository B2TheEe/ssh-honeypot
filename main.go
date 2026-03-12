package main

import (
	"log"
	"github.com/gliderlabs/ssh"
)

func main() {
	ssh.Handle(func(s ssh.Session) {
		log.Printf("Nieuwe verbinding van %s", s.RemoteAddr())
		s.Write([]byte("Welkom bij de SSH Honeypot! (Dit is een val.)\n"))
		s.Exit(0)
	})

	log.Println("SSH Honeypot luistert op :2222...")
	log.Fatal(ssh.ListenAndServe(":2222", nil))
}

