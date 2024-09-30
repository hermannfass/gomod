package publisher

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
)

func MakeUnixLineEnds(s string) string {
	// Unix: \n=LF=0x0A=\012=10  Mac: \r=CR=0x0D=\015=13  Dos: \r\n=CRLF
	re := regexp.MustCompile("\r\n?")
	s = re.ReplaceAllString(s, "\n")
	return s
}

	//              pre\/  start  \text  /   end   \post
	p := "(?sm)" + "(.*)(" + a + ")(.*)(" + b + ")(.*)"
	re, err := regexp.Compile(p)
	if err != nil {
		panic("Could not compile regexp: " + p)
	}
	m := re.FindStringSubmatch(s)
	r["pre"] = m[1] + m[2] // Before match including start tag
	r["tagged"] = m[3]
	r["post"] = m[4] + m[5] // After match prepended by end tag
	return (r)
}

func ReadHttpBody(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return (string(body))
}

/* Write the top of an HTML page. I WANT TO HARDCODE this for now! */
func WritePageTop(w io.Writer) {
	fmt.Fprintln(w,
		"<!DOCTYPE html>\n",
		"<html lang=\"de\">\n<head>\n<meta charset=\"utf-8\" />",
		"<title>Publisher</title>\n",
		"</head>\n<body>\n")
}
func WritePageBottom(w io.Writer) {
	fmt.Fprintln(w, "</body>\n</html>")
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("In defaultHandler")
	w.Header().Set("Content-Type", "text/html")
	WritePageTop(w)
	fmt.Fprintln(w,
		`<form action="/edit" method="get"><p>`,
		`<input type="text" name="u" value="http://www.astridco.de/band.html" /><br>`,
		`<input type="text" name="x" placeholder="vorname" autofocus />`,
		`<input type="submit" value="Seite aufrufen" />`,
		`</p></form>`)
	WritePageBottom(w)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	// Download page and present it here as a form.
	fmt.Println("In editHandler")
	// u := r.URL.Query()["u"][0]
	// x := r.URL.Query()["x"][0]
	u, x := "http://www.vonabiszet.de/", "Hermann"
	fmt.Println("Presenting editable section ", x)
	h := ReadHttpBody(u) // Get full HTML of that page (HTTP resp. body)
	in := `<!--` + x + `_start-->`
	out := `<!--` + x + `_end-->`
	parts := splitAtTags(h, in, out)
	w.Header().Set("Content-Type", "text/html")
	WritePageTop(w)
	fmt.Fprintln(w,
		parts["pre"]+
			`<form action="/process" method="post">`+
			`<input type="hidden" name="u" value="`+u+" />\n"+
			in+
			`<textarea name="update" rows="50" cols="70">`+
			parts["tagged"]+`</textarea>`+
			out+
			parts["post"])
	WritePageBottom(w)
}

func StartWebServer() *http.Server {
	// Yes, I roll my own
	m := http.NewServeMux()
	m.HandleFunc("/", defaultHandler)
	m.HandleFunc("/edit", editHandler)
	s := &http.Server{Addr: ":8080", Handler: m}
	go s.ListenAndServe()
	return (s)
}
