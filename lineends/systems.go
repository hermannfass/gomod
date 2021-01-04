package lineends

var LineEnd = make(map[string]string, 20)  // LineEnd["darwin"] = "\r"

func init() {
	LineEnd["unix"] = "\n"
	LineEnd["macos"] = "\n"
	LineEnd["classicmac"] = "\r"
	LineEnd["msdos"] = "\r\n"
	LineEnd["windows"] = LineEnd["msdos"]
}


