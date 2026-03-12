package main

import (
    "log"
    "github.com/gliderlabs/ssh"
)

func main() {
    // Configureer de SSH-server
    server := &ssh.Server{
        Addr: ":2222",
        PasswordHandler: func(ctx ssh.Context, password string) bool {
            log.Printf("Inlogpoging: gebruiker=%s, wachtwoord=%s, vanaf=%s",
                ctx.User(), password, ctx.RemoteAddr())
            // Accepteer altijd het wachtwoord voor de honeypot
            return true
        },
    }

    // Voeg een handler toe voor de sessie
    server.Handler = func(s ssh.Session) {
        log.Printf("Nieuwe sessie van %s", s.RemoteAddr())
        s.Write([]byte("Welkom bij de SSH Honeypot! (Dit is een val.)\n"))
        s.Exit(0)
    }

    // Start de server
    log.Println("SSH Honeypot luistert op :2222...")
    log.Fatal(server.ListenAndServe())
}

