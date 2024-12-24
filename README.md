## portmap

portmap is a fast portscan tool, uses shodan public data for port scan used internetdb.shodan.io and api.shodan.io/shodan/host

## Installation
```
go install -v github.com/rix4uni/portmap@latest
```

## Download prebuilt binaries
```
wget https://github.com/rix4uni/portmap/releases/download/v0.0.2/portmap-linux-amd64-0.0.2.tgz
tar -xvzf portmap-linux-amd64-0.0.2.tgz
rm -rf portmap-linux-amd64-0.0.2.tgz
mv portmap ~/go/bin/portmap
```
Or download [binary release](https://github.com/rix4uni/portmap/releases) for your platform.

## Compile from source
```
git clone --depth 1 github.com/rix4uni/portmap.git
cd portmap; go install
```

## Usage
```
                         __
    ____   ____   _____ / /_ ____ ___   ____ _ ____
   / __ \ / __ \ / ___// __// __  __ \ / __  // __ \
  / /_/ // /_/ // /   / /_ / / / / / // /_/ // /_/ /
 / .___/ \____//_/    \__//_/ /_/ /_/ \__,_// .___/
/_/                                        /_/
                    Current portmap version v0.0.2

A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.

Usage:
  portmap [flags]
  portmap [command]

Available Commands:
  apishodan   A brief description of your command uses https://api.shodan.io/shodan/host/
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  internetdb  A brief description of your command uses https://internetdb.shodan.io/

Flags:
  -h, --help      help for portmap
  -v, --version   Print the version of the tool and exit.
```

## Usage Example

Single IP:
```
# Basic Usage
▶ echo "147.249.56.149" | portmap internetdb
147.249.56.149:443
147.249.56.149:8080
147.249.56.149:8843

# Advanced Usage
▶ echo "147.249.56.149" | portmap apishodan
147.249.56.149:8000 [ASAS6419] [Fidelity National Information Services, Inc.]
147.249.56.149:8080 [ASAS6419] [Fidelity National Information Services, Inc.]
147.249.56.149:443 [ASAS6419] [Fidelity National Information Services, Inc.]
147.249.56.149:8843 [ASAS6419] [Fidelity National Information Services, Inc.]
147.249.56.149:8443 [ASAS6419] [Fidelity National Information Services, Inc.]

# CIDR range
▶ echo "1.2.3.4/24" | portmap apishodan
1.2.3.4:80
1.2.3.5:80

# Get JSON response
▶ echo "147.249.56.149" | portmap apishodan --json
{
  "ip_str": "147.249.56.149",
  "ports": [
    8000,
    8080,
    443,
    8843,
    8443
  ],
  "asn": "AS6419",
  "org": "Fidelity National Information Services, Inc.",
  "domains": [
    "automatedfinancial.com"
  ]
}
```

Multiple IPs:
```
▶ cat ips.txt
104.18.36.214
104.18.39.102
147.249.56.149
1.2.3.4/24
```

```
▶ cat ips.txt | portmap apishodan
▶ cat ips.txt | portmap internetdb
```
