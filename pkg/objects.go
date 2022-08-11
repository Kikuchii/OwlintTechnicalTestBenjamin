package pkg

type Comment struct {
	Id          string `json:"id"`
	TextFR      string `json:"textfr"`
	TextEn      string `json:"texten"`
	PublishedAt string `json:"publishedat"`
	AuthorID    string `json:"authorid"`
	TargetId    string `json:"targetid"`
}

type Faultymessage struct {
	Message string `json:"message"`
	Author  string `json:"author"`
}

type deepLResponse struct {
	Translations []translatedComment `json:"translations"`
}

type translatedComment struct {
	Detected_source_language string `json:"detected_source_language"`
	Text                     string `json:"text"`
}
