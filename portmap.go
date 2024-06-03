package main

import (
    "bufio"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "net"
    "net/http"
    "os"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
)

const version = "v0.0.1"

type ShodanResponse struct {
    Ports []int  `json:"ports"`
    ASN   string `json:"asn"`
}

func checkDomain(domain string) (string, error) {
    ip, err := net.LookupIP(domain)
    if err != nil {
        return "", err
    }
    if len(ip) > 0 {
        return ip[0].String(), nil
    }
    return "", fmt.Errorf("no IP found for domain")
}

func fetchShodanData(ipAddress string) (*ShodanResponse, error) {
    url := fmt.Sprintf("https://api.shodan.io/shodan/host/%s", ipAddress)
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    var data ShodanResponse
    err = json.Unmarshal(body, &data)
    if err != nil {
        return nil, err
    }
    return &data, nil
}

func checkDomains(domains []string, threads int, showASN, showIP bool, excludePorts []int, output io.Writer) {
    var wg sync.WaitGroup
    sem := make(chan struct{}, threads)

    for _, domain := range domains {
        wg.Add(1)
        go func(domain string) {
            defer wg.Done()
            sem <- struct{}{}
            defer func() { <-sem }()

            // Remove protocol prefix if present
            domain = strings.TrimPrefix(domain, "http://")
            domain = strings.TrimPrefix(domain, "https://")

            ipAddress, err := checkDomain(domain)
            if err != nil {
                return
            }

            data, err := fetchShodanData(ipAddress)
            if err != nil {
                time.Sleep(1 * time.Second)
                return
            }

            if data.Ports != nil {
                asn := ""
                if showASN && data.ASN != "" {
                    asn = data.ASN
                }
                ip := ""
                if showIP {
                    ip = ipAddress
                }
                ports := filterPorts(data.Ports, excludePorts)
                for _, port := range ports {
                    outputStr := fmt.Sprintf("%s:%d %s %s\n", domain, port, asn, ip)
                    fmt.Fprint(output, outputStr)
                }
            }
        }(domain)
    }
    wg.Wait()
}

func filterPorts(ports []int, excludePorts []int) []int {
    sort.Ints(excludePorts)
    filteredPorts := []int{}
    for _, port := range ports {
        if !contains(excludePorts, port) {
            filteredPorts = append(filteredPorts, port)
        }
    }
    return filteredPorts
}

func contains(s []int, e int) bool {
    i := sort.SearchInts(s, e)
    return i < len(s) && s[i] == e
}

func parsePorts(s string) ([]int, error) {
    if s == "" {
        return nil, nil
    }
    parts := strings.Split(s, ",")
    ports := make([]int, 0, len(parts))
    for _, p := range parts {
        port, err := strconv.Atoi(p)
        if err != nil {
            return nil, err
        }
        ports = append(ports, port)
    }
    return ports, nil
}

func main() {
    threads := flag.Int("c", 8, "Number of threads to use")
    showASN := flag.Bool("asn", false, "Show ASN")
    showIP := flag.Bool("ip", false, "Show IP address")
    excludePortsStr := flag.String("exclude-ports", "", "Exclude ports (comma-separated)")
    outputFile := flag.String("o", "", "Output file")
    versionFlag := flag.Bool("v", false, "Prints current version")
    flag.Parse()

    if *versionFlag {
        fmt.Printf("portmap version: %s\n", version)
        os.Exit(0)
    }

    excludePorts, err := parsePorts(*excludePortsStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error parsing exclude-ports: %v\n", err)
        os.Exit(1)
    }

    var output io.Writer
    if *outputFile != "" {
        file, err := os.Create(*outputFile)
        if err != nil {
            fmt.Fprintf(os.Stderr, "error creating output file: %v\n", err)
            os.Exit(1)
        }
        defer file.Close()
        output = io.MultiWriter(os.Stdout, file) // Write to both stdout and file
    } else {
        output = os.Stdout
    }

    scanner := bufio.NewScanner(os.Stdin)
    var domains []string
    for scanner.Scan() {
        domains = append(domains, scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Fprintln(os.Stderr, "reading standard input:", err)
        os.Exit(1)
    }

    checkDomains(domains, *threads, *showASN, *showIP, excludePorts, output)
}