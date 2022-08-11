package translate

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	deepLKey = "5e4e98a7-332f-9230-d405-3f6f28f7ebaf:fx"
)

type deepLResponse struct {
	Translations []translatedComment `json:"translations"`
}

type translatedComment struct {
	Detected_source_language string `json:"detected_source_language"`
	Text                     string `json:"text"`
}

func DeepLTranslate(source_lang string, target_lang string, comment string) string {
	urlstr := "https://api-free.deepl.com/v2/translate?auth_key=" + deepLKey + "&text=" + url.QueryEscape(comment) + "&target_lang=" + target_lang + "&source_lang=" + source_lang

	resp, err := http.Post(urlstr, "application/x-www-form-urlencoded", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	str, err := ioutil.ReadAll(resp.Body)

	var DeepLResponse deepLResponse
	json.Unmarshal(str, &DeepLResponse)

	return DeepLResponse.Translations[0].Text
}
