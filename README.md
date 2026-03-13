# 🍯 SSH Honeypot in Go

Een lichtgewicht SSH-honeypot geschreven in Go, die inlogpogingen en commando's van aanvallers logt.

## 📋 Functies

- Accepteert **alle** inlogpogingen (gebruikersnaam + wachtwoord worden gelogd)
- Simuleert een echte Linux-shell met nep-commando's
- Ondersteunt zowel **interactieve PTY-sessies** als **geautomatiseerde (non-PTY) verbindingen**
- Logt alle activiteit naar **terminal én `honeypot.log`**

## 🖥️ Gesimuleerde commando's

| Commando | Output |
|----------|--------|
| `ls` | Nep-bestanden |
| `whoami` | `root` |
| `id` | `uid=0(root)` |
| `uname -a` | Nep-kernelversie |
| `ps aux` | Nep-processen |
| `cat /etc/passwd` | Nep-gebruikerslijst |
| `history` | Nep-commandogeschiedenis |
| `wget` / `curl` | Logt de URL, geeft foutmelding |
| `sudo` | Doet alsof het werkt |
| `ifconfig` / `ip` | Nep-netwerkinterfaces |
| `uptime` | Nep-uptime |

## 🚀 Installatie & Gebruik

### Vereisten
- Go 1.18+

### Installeren
```bash
git clone https://github.com/B2TheEe/ssh-honeypot.git
cd ssh-honeypot
go mod tidy
```

### Starten
```bash
go run main.go
```

De honeypot luistert standaard op **poort 2222**.

### Verbinden (testen)
```bash
ssh -p 2222 testuser@localhost
# Wachtwoord: alles werkt
```

### Geautomatiseerd testen
```bash
./test.sh
```

## 📁 Logbestand

Alle activiteit wordt opgeslagen in `honeypot.log`:
```
2026/03/13 17:06:54 Inlogpoging: gebruiker=testuser, wachtwoord=lol, vanaf=127.0.0.1:34550
2026/03/13 17:06:54 Nieuwe sessie van 127.0.0.1:34550
2026/03/13 17:06:54 Commando uitgevoerd: id
```

Live meekijken:
```bash
tail -f honeypot.log
```

## 📦 Afhankelijkheden

- [gliderlabs/ssh](https://github.com/gliderlabs/ssh) — SSH-server bibliotheek voor Go

## ⚠️ Disclaimer

Dit project is uitsluitend bedoeld voor **educatieve doeleinden** en **eigen netwerken**. Gebruik het nooit op systemen zonder toestemming.

## 📄 Licentie

MIT License
