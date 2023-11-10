package main

import (
    "crypto/md5"
    "encoding/hex"
    "flag"
    "fmt"
    "math/rand"
    "strings"
    "time"
)

var (
    location  = flag.String("l", "UTC", "Location (country/city) for timezone")
    md5Format = flag.Bool("m", false, "Use MD5 format for Message-ID")
    smallD    = flag.Bool("i", false, "Use small 'd' in Message-Id")
    domain    = flag.String("d", "", "Domain part")
)

func main() {
    flag.Parse()

    if flag.NFlag() == 0 {
        printUsage()
        return
    }

    var domainComponent string

    if *domain != "" {
        domainComponent = *domain
    } else {
        fmt.Print("Enter domain: ")
        fmt.Scanln(&domainComponent)
    }

    headers := messageID(domainComponent)
    fmt.Println(headers)
}

func messageID(domainComponent string) (headers string) {
    loc, err := time.LoadLocation(*location)
    if err != nil {
        fmt.Printf("Invalid location: %v\n", err)
        return
    }

    now := time.Now().In(loc)
    dateComponent := now.Format("20060102.150405")
    randomComponent := hex.EncodeToString(randBytes(4))

    messageId := ""
    if *md5Format {
        h := md5.New()
        h.Write([]byte(dateComponent + "." + randomComponent))
        messageId = fmt.Sprintf(
            "Message-ID: <%x@%s>",
            h.Sum(nil),
            domainComponent,
        )
    } else {
        messageId = fmt.Sprintf(
            "Message-ID: <%s.%s@%s>",
            dateComponent,
            randomComponent,
            domainComponent,
        )
    }

    if *smallD {
        messageId = strings.Replace(messageId, "Message-ID", "Message-Id", 1)
    }

    dateHeader := fmt.Sprintf(
        "Date: %s",
        now.Format("Mon, 02 Jan 2006 15:04:05 -0700"),
    )

    headers = messageId + "\n" + dateHeader
    return
}

func randBytes(n int) []byte {
    b := make([]byte, n)
    _, err := rand.Read(b)
    if err != nil {
        panic(err)
    }
    return b
}

func printUsage() {
    fmt.Println("Usage: mid -l <location> [-m] [-i] [-d <domain>] or provide domain as input")
    flag.PrintDefaults()
}
