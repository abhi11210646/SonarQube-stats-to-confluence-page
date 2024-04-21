package confluence

type Page struct {
	Title   string  `json:"title"`
	Type    string  `json:"type"`
	Body    Body    `json:"body"`
	Version Version `json:"version"`
}
type Body struct {
	Storage `json:"storage"`
}

type Storage struct {
	Value          string `json:"value"`
	Representation string `json:"representation"`
}
type Version struct {
	Number  int8   `json:"number"`
	Message string `json:"message"`
}
