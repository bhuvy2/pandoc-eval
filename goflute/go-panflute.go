package goflute

type Elem struct {
	Type string `json:"t"`
	Content []interface{} `json:"c"`
}

type Document struct {
	PandocApiVersion []int64 `json:"pandoc-api-version"`
	Metadata map[string]interface{} `json:"meta"`
	Blocks []Elem `json:"blocks"`
}

func must(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func RunFilter() {
	var record Document;
	err := json.NewDecoder(os.Stdin).Decode(&record)
	must(err)

	fmt.Fprintf(os.Stderr, "%v\n", record)
}
