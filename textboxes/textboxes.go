package textboxes

import(
	"fmt"  // testing only
	"strings"
)

var TextFields = []string{"Title", "Subtitle", "Author", "Date"}

func maxWidth(ss []string) int {
	var max int = 0
	var l int = 0
	for _, s := range(ss) {
		l = len([]rune(s))
		if l > max {
			max = l
		}
	}
	return max
}

func CenterText(w int, t string) string {
	spcL := strings.Repeat(" ", (w-len([]rune(t)))/2)
	spcR := strings.Repeat(" ", w-len([]rune(t))-len([]rune(spcL)))
	return(spcL + t + spcR)
}

// Shorten a text to the total length as specified (l) including one
// appended ellipsis character (\u2026, making it not ASCII safe!).
func shortenTo(l int, s string) string {
	fmt.Println("Shortening", s, "to", l)
	r := []rune(s)
	t := string(r[:l-1]) + "\u2026"
	fmt.Printf("Shortened text\n%s\n", t)
	return t
}

// Distribute two strings, justified in a line
// with the specified total width (tw).
func AlignText(tw int, t1, t2 string) string {
	// Available width: Deduct 2× border, 2× padding
	d := 3       // Min. distance, space between t1 and t2
	w := tw - d  // Space for t1 and t2
	l1 := len([]rune(t1))
	l2 := len([]rune(t2))
	if l1 + l2 + d > w {
		switch {
			case l1 < w/2:
				// leave t1 as is and shorten t2
				// t2 = string([]rune(t2)[:l1-1]) + "\u2026" // cut 1 for"…"
				t2 = shortenTo(l1, t2)
				l2 = len([]rune(t2))
			case l2 < w/2:
				// t1 = string([]rune(t1)[:w-l2-1]) + "\u2026" // cut 1 for"…"
				t1 = shortenTo(w-l2, t1)
				l1 = len([]rune(t1))
			default:
				t1 = shortenTo(w/2, t1)
				l1 = len([]rune(t1))
				fmt.Println("Shortened t1:", t1)
				t2 = shortenTo(w-l1, t2)
				fmt.Println("Shortened t2:", t2)
				l2 = len([]rune(t2))
		}
	}
	dist := tw - l1 - l2 // Space between t1 to t2 
	return(t1 + strings.Repeat(" ", dist) + t2)
}

func space(w int) string {
	return strings.Repeat(" ", w)
}

func line(left, content, right string, pad int) string {
	return left + space(pad) + content + space(pad) + right
}

func HeaderBox(bw int, styleName string, t map[string]string) string {
	s := makeStyle(styleName)
	border := 1
	pad := 1
	iw := bw - 2*border         // inner width (no border)
	tw := iw - 2*pad            // text width (without border & padding)
	lc := 7  // How many lines
	withSubtitle := t["Subtitle"] != ""
	if withSubtitle {
		lc = 9
	}
	l := make([]string, 0, lc)
	spaceLine := line(s.v, space(tw), s.v, pad)
	l = append(l, line(s.a, strings.Repeat(s.h, iw), s.b, 0))    // top
	l = append(l, spaceLine)
	if len(t["Title"]) > tw {
		t["Title"] = shortenTo(tw, t["Title"])
	}
	l = append(l, line(s.v, CenterText(tw, t["Title"]), s.v, pad)) // title
	if withSubtitle {
		if len(t["Subtitle"]) > bw-4 {
			t["Subtitle"] = shortenTo(tw, t["Subtitle"])
		}
		l = append(l, spaceLine)
		l = append(l, line(s.v, CenterText(tw, t["Subtitle"]), s.v, pad))
	}
	l = append(l, spaceLine)
	l = append(l, line(s.l, strings.Repeat(s.s, iw), s.r, 0))     // hline
	l = append(l, line(s.v, AlignText(tw, t["Author"], t["Date"]), s.v, pad))
	l = append(l, line(s.c, strings.Repeat(s.h, iw), s.d, 0))
	return strings.Join(l, "\n") + "\n"
}

func FrameText(styleName string, t string) string {
	hpad := 2  // horizontal padding
	vpad := 1  // vertical padding
	s := makeStyle(styleName)
	t = strings.TrimSpace(t)
	tl := strings.Split(t, "\n")
	tw := maxWidth(tl)  // text width
	iw := tw + 2*hpad   // inner width, i.e. without borders
	l := make([]string, 0, len(tl) + 2 + 2*vpad) // text, border, padding
	spaceLine := line(s.v, space(tw), s.v, hpad)
	l = append(l, line(s.a, strings.Repeat(s.h, iw), s.b, 0))
	l = append(l, spaceLine)
	for _, txt := range(tl) {
		l = append(l, line(s.v, CenterText(tw, txt), s.v, hpad))
	}
	l = append(l, spaceLine)
	l = append(l, line(s.c, strings.Repeat(s.h, iw), s.d, 0))
	return strings.Join(l, "\n") + "\n"
}

