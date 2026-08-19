package entity

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type TagNamesByLang map[string][]string

func (t *TagNamesByLang) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		*t = nil
		return nil
	}
	switch data[0] {
	case '{':
		var m map[string][]string
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		*t = TagNamesByLang(m)
		return nil
	case '[':
		var names []string
		if err := json.Unmarshal(data, &names); err != nil {
			return fmt.Errorf("tags: expected object map or string array: %w", err)
		}
		*t = TagNamesByLang{"uz": names}
		return nil
	default:
		return fmt.Errorf("tags must be an object or array")
	}
}

func TagNamesByLangFromEntities(tags []*TagEntity) TagNamesByLang {
	if len(tags) == 0 {
		return nil
	}
	out := make(TagNamesByLang)
	for _, tag := range tags {
		if tag == nil || tag.Name == "" {
			continue
		}
		lang := tag.Lang
		if lang == "" {
			lang = "uz"
		}
		out[lang] = append(out[lang], tag.Name)
	}
	return out
}
