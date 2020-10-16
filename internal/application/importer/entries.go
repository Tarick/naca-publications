package importer

type Publisher struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Publication struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	LanguageCode string `json:"language_code"`
	Type         string `json:"type"`
	// Config content is different for different publication types.
	// when parsing, we decide on Type
	Config PublicationConfig `json:"config"`
}

// PublicationConfig is used to pass around different config structs
type PublicationConfig interface{}

// Entrie is one record in file, representing publisher and its publications
type Entrie struct {
	Publisher    Publisher     `json:"publisher"`
	Publications []Publication `json:"publications"`
}
