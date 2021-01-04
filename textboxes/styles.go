package textboxes

const bw int = 76  // box width, i.e. text width all inclusive to be used

type style struct {
	h string // Horizontal line
	v string // Vertical line
	a string // Top left corner
	b string // Top right corner
	c string // Bottom left corner
	d string // Bottom right corner
	l string // Left separator connection
	r string // Right separator connection
	s string // Separator line
}

func makeStyle(styleName string) style {
	switch styleName {
	case "ascii":
		//            h    v    a    b    c    d    l    r    s
		return style{"_", "|", "_", "_", "|", "|", "|", "|", "_"} 
	case "single":
		//            h    v    a    b    c    d    l    r    s
		return style{"─", "│", "┌", "┐", "└", "┘", "├", "┤", "─"} 
	case "double":
		//            h    v    a    b    c    d    l    r    s
		return style{"═", "║", "╔", "╗", "╚", "╝", "╠", "╣", "═"} 
	default: // "mixed"
		//            h    v    a    b    c    d    l    r    s
		return style{"═", "║", "╔", "╗", "╚", "╝", "╟", "╢", "─"}
	}
	
}
/* Layout:
	    a┌───h─────────┐b
	     │             │
	     v             │
	     │             │
	    l├────s────────┤r
	     │             │
	    c└─────────────┘d
*/


