package punct

//Punctuations keep punctuation symbols
var Punctuations map[rune]bool

const (
	//StartQuote for LT
	StartQuote = '„'
	//EndQuote for LT
	EndQuote = '“'
)

func init() {
	Punctuations = make(map[rune]bool)
	for _, r := range ",.!?-:;" {
		Punctuations[r] = true
	}
}

//PureWord retun word without punctuation
func PureWord(w string) string {
	rs := []rune(w)
	l := len(rs)
	for ; l > 0 && isPunct(rs[l-1]); l-- {
	}
	rs = rs[0:l]
	if len(rs) > 0 && rs[0] == StartQuote {
		rs = rs[1:]
	}
	if len(rs) > 0 && rs[len(rs)-1] == EndQuote {
		rs = rs[:len(rs)-1]
	}
	return string(rs)
}

func isPunct(r rune) bool {
	_, ok := Punctuations[r]
	return ok
}
