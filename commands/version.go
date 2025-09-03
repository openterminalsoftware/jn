package commands

import (
	"fmt"
	"jn/colors"
)

func Version() {
	fmt.Println("You are using " + colors.Bold + colors.Magenta + "jn 1.0.0" + colors.Reset)
	fmt.Println("Created by OpenTerminalSoftware (OTS)")
	fmt.Println(colors.Bold + "https://github.com/openterminalsoftware/jn" + colors.Reset)
}
