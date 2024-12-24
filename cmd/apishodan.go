package cmd

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/projectdiscovery/mapcidr"
	"github.com/spf13/cobra"
)

type apishodanShodanResponse struct {
	IPStr   string   `json:"ip_str"`
	Ports   []int    `json:"ports"`
	ASN     string   `json:"asn"`
	Org     string   `json:"org"`
	Domains []string `json:"domains"`
}

func apishodanfetchShodanData(ip string) (*apishodanShodanResponse, error) {
	url := fmt.Sprintf("https://api.shodan.io/shodan/host/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var data apishodanShodanResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func apishodanprintPlainText(data *apishodanShodanResponse) {
	for _, port := range data.Ports {
		fmt.Printf("%s:%d [AS%s] [%s]\n", data.IPStr, port, data.ASN, data.Org)
	}
}

func apishodanprintJSON(data *apishodanShodanResponse) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func apishodanprocessInput(input io.Reader, apishodanjsonFlag bool) {
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		ip := strings.TrimSpace(scanner.Text())

		// Remove protocol prefix if present
		ip = strings.TrimPrefix(ip, "http://")
		ip = strings.TrimPrefix(ip, "https://")

		if ip == "" {
			continue
		}

		// If input is CIDR, expand to individual IPs
		if strings.Contains(ip, "/") {
			ips, err := mapcidr.IPAddresses(ip)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error processing CIDR %s: %v\n", ip, err)
				continue
			}

			// Fetch Shodan data for each IP in the CIDR range
			for _, ipAddr := range ips {
				data, err := apishodanfetchShodanData(ipAddr)
				if err != nil {
					if apishodanverboseFlag {
						fmt.Fprintf(os.Stderr, "Error fetching data for IP %s: %v\n", ipAddr, err)
					}
					continue
				}

				if apishodanjsonFlag {
					apishodanprintJSON(data)
				} else {
					apishodanprintPlainText(data)
				}
			}
		} else {
			// Process a single IP
			data, err := apishodanfetchShodanData(ip)
			if err != nil {
				if apishodanverboseFlag {
					fmt.Fprintf(os.Stderr, "Error fetching data for IP %s: %v\n", ip, err)
				}
				continue
			}

			if apishodanjsonFlag {
				apishodanprintJSON(data)
			} else {
				apishodanprintPlainText(data)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

// apishodanCmd represents the apishodan command
var apishodanCmd = &cobra.Command{
	Use:   "apishodan",
	Short: "A brief description of your command uses https://api.shodan.io/shodan/host/",
	Long: `A longer description of your command uses https://api.shodan.io/shodan/host/.

Examples:
 echo "1.2.3.4" | portmap apishodan
 echo "1.2.3.4/24" | portmap apishodan
 cat ips.txt | portmap apishodan
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(flag.Args()) == 0 {
			// No arguments, read from stdin
			apishodanprocessInput(os.Stdin, apishodanjsonFlag)
		} else {
			// Process files given as arguments
			for _, fileName := range flag.Args() {
				file, err := os.Open(fileName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", fileName, err)
					continue
				}
				apishodanprocessInput(file, apishodanjsonFlag)
				file.Close()
			}
		}
	},
}

var apishodanjsonFlag bool
var apishodanverboseFlag bool

func init() {
	rootCmd.AddCommand(apishodanCmd)

	apishodanCmd.Flags().BoolVar(&apishodanjsonFlag, "json", false, "Output in JSON format")
	apishodanCmd.Flags().BoolVar(&apishodanverboseFlag, "verbose", false, "enable verbose mode")
}
