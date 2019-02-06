package aeridya

import "fmt"

const (
	NAME     = "Aeridya"
	MAJORVER = "0"
	MINORVER = "4"
	VERTAG   = "-alpha"
	DESC     = "Server and CMS"
)

// Version returns a formatted string of the name/version number
func Version() string {
	return fmt.Sprintf("%s v%s.%s%s", NAME, MAJORVER, MINORVER, VERTAG)
}

// Info returns a formatted string of Version and the Description
func Info() string {
	return fmt.Sprintf("%s\n\t%s", Version(), DESC)
}
