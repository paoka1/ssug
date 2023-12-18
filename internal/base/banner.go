package base

import (
	"fmt"
)

var banner = `
     _        _        _        _    
   _( )__   _( )__   _( )__   _( )__ 
 _|     _|_|     _|_|     _|_|     _|
(_ S _ (_(_ S _ (_(_ U _ (_(_ G _ (_ 
  |_( )__| |_( )__| |_( )__| |_( )__|
`

func GetBanner() string {
	return banner +
		"\n                      " +
		"Shauio's short URL generator %s\n"
}

func PrintBanner() {
	fmt.Println(fmt.Sprintf(GetBanner(), Version))
}
