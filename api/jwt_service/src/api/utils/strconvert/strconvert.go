package strconvert

import (
	"fmt"
	"unicode"
	"strings"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func ConvertVitoEn(str string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	result, _, _ := transform.String(t, str)
	fmt.Println(result)
	//xoa khoang trang
	result = strings.Replace(result, " ", "", -1)
	result = strings.ToLower(result)
	return result
}
