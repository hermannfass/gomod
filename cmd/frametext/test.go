package main

import(
	"fmt"
	"github.com/hermannfass/gomod/textboxes"
)

func main() {
	s := `
Hallo
das hier ist Zeile zwei
und das ist Zeile drei und die längste
vor der Schlusszeile
`
	fmt.Println(textboxes.FrameText("mixed", s))
}
