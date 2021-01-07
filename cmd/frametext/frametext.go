package main

import(
	"os"
	"fmt"
	"github.com/hermannfass/gomod/textboxes"
)

func main() {
	var t string
	if len(os.Args) > 1 {
		t = os.Args[1]
	} else {
		t = "This is a box.\nSpecify your text as parameter, or:\n" +
		    "Use the function textboxes.FrameText() from your code.\n\n" +
		    "The call for what you see here would be:\n" +
		    "textboxes.FrameText(\"single\", \"This is a box ...\")"
	}
	fmt.Println(textboxes.FrameText("single", t))
}
