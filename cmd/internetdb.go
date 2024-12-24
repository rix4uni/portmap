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

type ShodanResponse struct {
	IPS       string   `json:"ip"`
	Ports     []int    `json:"ports"`
	Hostnames []string `json:"hostnames"`
}

func fetchShodanData(ip string) (*ShodanResponse, error) {
	url := fmt.Sprintf("https://internetdb.shodan.io/%s", ip)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data: %s", resp.Status)
	}

	var data ShodanResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func printPlainText(data *ShodanResponse) {
	for _, port := range data.Ports {
		fmt.Printf("%s:%d\n", data.IPS, port)
	}
}

func printJSON(data *ShodanResponse) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating JSON: %v\n", err)
		return
	}
	fmt.Println(string(jsonData))
}

func processInput(input io.Reader, jsonFlag bool) {
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
				data, err := fetchShodanData(ipAddr)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error fetching data for IP %s: %v\n", ipAddr, err)
					continue
				}

				if jsonFlag {
					printJSON(data)
				} else {
					printPlainText(data)
				}
			}
		} else {
			// Process a single IP
			data, err := fetchShodanData(ip)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error fetching data for IP %s: %v\n", ip, err)
				continue
			}

			if jsonFlag {
				printJSON(data)
			} else {
				printPlainText(data)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

// internetdbCmd represents the internetdb command
var internetdbCmd = &cobra.Command{
	Use:   "internetdb",
	Short: "A brief description of your command uses https://internetdb.shodan.io/",
	Long: `A longer description of your command uses https://internetdb.shodan.io/.

Examples:
 echo "1.2.3.4" | portmap internetdb
 echo "1.2.3.4/24" | portmap internetdb
 cat ips.txt | portmap internetdb
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(flag.Args()) == 0 {
			// No arguments, read from stdin
			processInput(os.Stdin, jsonFlag)
		} else {
			// Process files given as arguments
			for _, fileName := range flag.Args() {
				file, err := os.Open(fileName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error opening file %s: %v\n", fileName, err)
					continue
				}
				processInput(file, jsonFlag)
				file.Close()
			}
		}
	},
}

var jsonFlag bool

func init() {
	rootCmd.AddCommand(internetdbCmd)

	internetdbCmd.Flags().BoolVar(&jsonFlag, "json", false, "Output in JSON format")
}
