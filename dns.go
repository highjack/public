// receiver.go
package main

import (
    "encoding/base64"
    "fmt"
    "github.com/miekg/dns"
    "os"
    "strings"
    "sync"
)

var (
    dataStore = make(map[string][]string)
    mutex     = &sync.Mutex{}
)

func handleDNSRequest(w dns.ResponseWriter, r *dns.Msg) {
    msg := dns.Msg{}
    msg.SetReply(r)

    for _, q := range r.Question {
        if q.Qtype == dns.TypeA {
            domain := strings.ToLower(q.Name)
            parts := strings.Split(domain, ".")

            if len(parts) < 3 {
                continue // Need at least 1 chunk + filename + domain
            }

            // Last two parts are filename and the receiver hostname
            filename := parts[len(parts)-2]
            chunks := parts[:len(parts)-2]

            // Rejoin the data from subdomain segments (already base64-safe)
            chunkData := strings.Join(chunks, "")
            chunkData = strings.ReplaceAll(chunkData, "-", "") // in case hyphens were used

            mutex.Lock()
            dataStore[filename] = append(dataStore[filename], chunkData)
            mutex.Unlock()

            fmt.Printf("[+] Received chunk for %s (%d bytes)\n", filename, len(chunkData))
        }
    }

    w.WriteMsg(&msg)
}

func saveFiles() {
    mutex.Lock()
    defer mutex.Unlock()

    for filename, chunks := range dataStore {
        fullData := strings.Join(chunks, "")

        // Ensure proper base64 padding
        missing := len(fullData) % 4
        if missing != 0 {
            fullData += strings.Repeat("=", 4-missing)
        }

        decoded, err := base64.StdEncoding.DecodeString(fullData)
        if err != nil {
            fmt.Printf("[-] Failed to decode base64 for %s: %s\n", filename, err)
            continue
        }

        err = os.WriteFile(filename, decoded, 0644)
        if err != nil {
            fmt.Printf("[-] Failed to write file %s: %s\n", filename, err)
        } else {
            fmt.Printf("[+] Successfully reconstructed file: %s\n", filename)
        }
    }
}

func startServer() {
    dns.HandleFunc(".", handleDNSRequest)
    server := &dns.Server{Addr: ":53", Net: "udp"}
    fmt.Println("[*] DNS server listening on port 53...")
    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("[-] DNS server error: %s\n", err)
        os.Exit(1)
    }
}

func main() {
    go startServer()
    fmt.Println("[*] Press Enter to stop and save reconstructed files...")
    fmt.Scanln()
    saveFiles()
}
