package fetcher

import (
	"regexp"
	"strings"
)

func RmIllegalChar(s *string) {
	var re = regexp.MustCompile(`[\/:*?"<>|]`)
	*s = re.ReplaceAllString(*s, "")
}

func ReplaceIllegalChar(s *string) {
	*s = strings.ReplaceAll(*s, "\\", "、")
	*s = strings.ReplaceAll(*s, "/", "／")
	*s = strings.ReplaceAll(*s, "|", "｜")
	*s = strings.ReplaceAll(*s, "?", "？")
	*s = strings.ReplaceAll(*s, ":", "：")
	*s = strings.ReplaceAll(*s, "*", "＊")
	*s = strings.ReplaceAll(*s, "<", "《")
	*s = strings.ReplaceAll(*s, ">", "》")
	*s = strings.ReplaceAll(*s, "\"", "“")
}
