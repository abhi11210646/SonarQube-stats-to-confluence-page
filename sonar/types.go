package sonar

type Stats struct {
	Component struct {
		Name     string
		Measures []Measure `json:"measures"`
	} `json:"component"`
}
type Measure struct {
	Metric string `json:"metric"`
	Value  string `json:"value"`
}
