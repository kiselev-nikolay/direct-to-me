package storage

type Redirect struct {
	FromURI         string `json:"from"`
	ToURL           string `json:"to"`
	RedirectAfter   string `json:"after"`
	URLTemplate     string `json:"urlTemplate"`
	MethodTemplate  string `json:"methodTemplate"`
	HeadersTemplate string `json:"headersTemplate"`
	BodyTemplate    string `json:"bodyTemplate"`
}
