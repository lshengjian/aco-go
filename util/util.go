package util
import (
	"os"
	"fmt"
)
//prints error and exits on abnormal conditions
func PrintErrorAndExit(err error) {
	fmt.Print(err)
	os.Exit(2)
}

