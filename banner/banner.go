package banner

import (
	"fmt"
)

// prints the version message
const version = "v0.0.3"

func PrintVersion() {
	fmt.Printf("Current paramfinder version %s\n", version)
}

// Prints the Colorful banner
func PrintBanner() {
	banner := `
                                           ____ _             __           
    ____   ____ _ _____ ____ _ ____ ___   / __/(_)____   ____/ /___   _____
   / __ \ / __  // ___// __  // __  __ \ / /_ / // __ \ / __  // _ \ / ___/
  / /_/ // /_/ // /   / /_/ // / / / / // __// // / / // /_/ //  __// /    
 / .___/ \__,_//_/    \__,_//_/ /_/ /_//_/  /_//_/ /_/ \__,_/ \___//_/
/_/
`
	fmt.Printf("%s\n%55s\n\n", banner, "Current paramfinder version "+version)
}
