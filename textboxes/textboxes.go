package textboxes

import(
	"strings"
)

var TextFields = []string{"Title", "Subtitle", "Author", "Date"}

func centerText(t string) string {
	w := bw - 2 // Available width. Deduct borders as out of scope
	spcL := strings.Repeat(" ", (w-len([]rune(t)))/2)
	spcR := strings.Repeat(" ", w-len([]rune(t))-len([]rune(spcL)))
	return(spcL + t + spcR)
}

func alignText(t1, t2 string) string {
	w := bw - 2 // Available width. Deduct borders as out of scope
	tw := len([]rune(t1)) + len([]rune(t2)) // total text width
	dist := w - tw - 2 // Distance t1 to t2 (2 spaces, 1 at each border)
	interspace := strings.Repeat(" ", dist)
	return(" " + t1 + interspace + t2 + " ")
}
	
func HeaderBox(styleName string, t map[string]string) string {
	s := makeStyle(styleName)
	lc := 7  // How many lines
	withSubtitle := t["Subtitle"] != ""
	if withSubtitle { lc = 9 }
	l := make([]string, 0, lc)
	spaceLine := s.v + strings.Repeat(" ", bw-2) + s.v
	l = append(l, s.a + strings.Repeat(s.h, bw-2) + s.b)               // top
	l = append(l, spaceLine)
	l = append(l, s.v + centerText(t["Title"]) + s.v)          // title
	if withSubtitle {
		l = append(l, spaceLine)
		l = append(l, s.v + centerText(t["Subtitle"]) + s.v)
	}
	l = append(l, spaceLine)
	l = append(l, s.l + strings.Repeat(s.s, bw-2) + s.r)         // separator 
	l = append(l, s.v + alignText(t["Author"], t["Date"]) + s.v)    // QA
	l = append(l, s.c + strings.Repeat(s.h, bw-2) + s.d)            // bottom
	return strings.Join(l, "\n") + "\n"
}

