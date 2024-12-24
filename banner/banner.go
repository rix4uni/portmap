package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.2"

func PrintVersion() {
	fmt.Printf("Current portmap version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                         __                         
    ____   ____   _____ / /_ ____ ___   ____ _ ____ 
   / __ \ / __ \ / ___// __// __  __ \ / __  // __ \
  / /_/ // /_/ // /   / /_ / / / / / // /_/ // /_/ /
 / .___/ \____//_/    \__//_/ /_/ /_/ \__,_// .___/ 
/_/                                        /_/`
	fmt.Printf("%s\n%50s\n\n", banner, "Current portmap version "+version)
}
