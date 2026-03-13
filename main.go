package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
)

func main() {
	server := &ssh.Server{
		Addr: ":2222",
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			log.Printf("Inlogpoging: gebruiker=%s, wachtwoord=%s, vanaf=%s",
				ctx.User(), password, ctx.RemoteAddr())
			return true
		},
		Handler: func(s ssh.Session) {
			log.Printf("Nieuwe sessie van %s", s.RemoteAddr())

			pty, winCh, isPty := s.Pty()
			if !isPty {
				fmt.Fprintf(s, "Geen pseudo-terminal beschikbaar.\n")
				s.Exit(1)
				return
			}
			log.Printf("Pseudo-terminal aangevraagd: %v", pty)

			go func() {
				for win := range winCh {
					log.Printf("Window resize: %d x %d", win.Width, win.Height)
				}
			}()

			fmt.Fprintf(s, "Welkom bij de SSH Honeypot! (Dit is een val.)\r\n")

			buf := make([]byte, 1)
			var line []byte

			// Toon eerste prompt
			fmt.Fprint(s, "$ ")

			for {
				_, err := s.Read(buf)
				if err != nil {
					if err == io.EOF {
						s.Exit(0)
					} else {
						log.Printf("Fout bij lezen: %v", err)
						s.Exit(1)
					}
					return
				}

				b := buf[0]

				switch {
				case b == '\r' || b == '\n':
					// Enter ingedrukt: verwerk de regel
					fmt.Fprint(s, "\r\n")
					input := strings.TrimSpace(string(line))
					line = nil

					if input != "" {
						log.Printf("Commando uitgevoerd: %s", input)
					}

					switch input {
					case "ls":
						fmt.Fprintf(s, "fake_file1.txt  fake_file2.txt  fake_directory\r\n")
					case "whoami":
						fmt.Fprintf(s, "root\r\n")
					case "pwd":
						fmt.Fprintf(s, "/home/fake_user\r\n")
					case "exit":
						fmt.Fprintf(s, "exit\r\n")
						s.Exit(0)
						return
					case "":
						// lege regel, doe niets
					default:
						fmt.Fprintf(s, "%s: command not found\r\n", input)
					}

					fmt.Fprint(s, "$ ")

				case b == 127 || b == 8:
					// Backspace
					if len(line) > 0 {
						line = line[:len(line)-1]
						// Verwijder karakter van terminal
						fmt.Fprint(s, "\b \b")
					}

				case b >= 32 && b < 127:
					// Normaal afdrukbaar karakter
					line = append(line, b)
					// Echo terug naar client
					s.Write([]byte{b})
				}
			}
		},
	}

	log.Println("SSH Honeypot luistert op :2222...")
	log.Fatal(server.ListenAndServe())
}
