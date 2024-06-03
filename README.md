# portmap
 
# Installation
```
go install -v github.com/rix4uni/portmap@latest
```

# Usage
```
Usage of portmap:
  -asn
        Show ASN
  -c int
        Number of threads to use (default 8)
  -exclude-ports string
        Exclude ports (comma-separated)
  -ip
        Show IP address
  -o string
        Output file
  -v    Prints current version
```

# Usage Example

Single URL:
```
echo "hackerone.com" | portmap
hackerone.com:8080
hackerone.com:2082
hackerone.com:2083
hackerone.com:2053
hackerone.com:2086
hackerone.com:2087
hackerone.com:80
hackerone.com:8880
hackerone.com:8443
hackerone.com:443

echo "hackerone.com" | portmap -c 100 -asn
hackerone.com:8080 AS13335
hackerone.com:2082 AS13335
hackerone.com:2083 AS13335
hackerone.com:2053 AS13335
hackerone.com:2086 AS13335
hackerone.com:2087 AS13335
hackerone.com:80 AS13335
hackerone.com:8880 AS13335
hackerone.com:8443 AS13335
hackerone.com:443 AS13335

echo "hackerone.com" | portmap -c 100 -ip
hackerone.com:8080  104.18.36.214
hackerone.com:2082  104.18.36.214
hackerone.com:2083  104.18.36.214
hackerone.com:2053  104.18.36.214
hackerone.com:2086  104.18.36.214
hackerone.com:2087  104.18.36.214
hackerone.com:80  104.18.36.214
hackerone.com:8880  104.18.36.214
hackerone.com:8443  104.18.36.214
hackerone.com:443  104.18.36.214

echo "hackerone.com" | portmap -c 100 -asn -ip
hackerone.com:8080 AS13335 104.18.36.214
hackerone.com:2082 AS13335 104.18.36.214
hackerone.com:2083 AS13335 104.18.36.214
hackerone.com:2053 AS13335 104.18.36.214
hackerone.com:2086 AS13335 104.18.36.214
hackerone.com:2087 AS13335 104.18.36.214
hackerone.com:80 AS13335 104.18.36.214
hackerone.com:8880 AS13335 104.18.36.214
hackerone.com:8443 AS13335 104.18.36.214
hackerone.com:443 AS13335 104.18.36.214

## remove port 80,443
echo "hackerone.com" | portmap -c 100 -exclude-ports 80,443
hackerone.com:8080
hackerone.com:2082
hackerone.com:2083
hackerone.com:2053
hackerone.com:2086
hackerone.com:2087
hackerone.com:8880
hackerone.com:8443

## save output in a file
echo "hackerone.com" | portmap -c 100 -exclude-ports 80,443 -o output.txt
hackerone.com:8080
hackerone.com:2082
hackerone.com:2083
hackerone.com:2053
hackerone.com:2086
hackerone.com:2087
hackerone.com:8880
hackerone.com:8443
```

Multiple URLs:
```
cat urls.txt | portmap -c 100
```
subs.txt contains:
```
hackerone.com
notion.so
```

Output:
```
hackerone.com:8080
hackerone.com:2082
hackerone.com:2083
hackerone.com:2053
hackerone.com:2086
hackerone.com:2087
hackerone.com:80
hackerone.com:8880
hackerone.com:8443
hackerone.com:443
notion.so:2096
notion.so:2082
notion.so:2083
notion.so:2053
notion.so:2086
notion.so:2087
notion.so:2095
notion.so:80
notion.so:8880
notion.so:8080
notion.so:8443
notion.so:443
```
