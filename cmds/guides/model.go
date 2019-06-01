package guides

type Guide struct {
	Code     string
	Name     string
	URL      string `json:"url"`
	Type     string
	ECTS     string `json:"ects"`
	Semester string `json:"semestre"`
}

type GuideList map[string][]*Guide
