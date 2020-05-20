package punct

//Punctuations keep punctuation symbols
var Punctuations map[rune]bool

var sepSymbols map[rune]bool

var parenthesesSymbols map[rune]bool

const (
	//StartQuote for LT
	StartQuote = '„'
	//EndQuote for LT
	EndQuote = '“'
)

func init() {
	Punctuations = make(map[rune]bool)
	for _, r := range ",./\\;:*&^%$#@!~?-+=" {
		Punctuations[r] = true
	}
	parenthesesSymbols = make(map[rune]bool)
	for _, r := range []rune("[]{}<>()") {
		parenthesesSymbols[r] = true
	}
}

//PureWord retun word without punctuation, word, <string before word>, <string after word>
func PureWord(w string) (string, string, string) {
	return TrimWord(w, IsAllSep)
}

//TrimWord return word without runes matching f, word, <string before word>, <string after word>
func TrimWord(w string, f func(rune) bool) (string, string, string) {
	rns := []rune(w)
	i1, i2 := 0, len(rns)-1
	for ; i1 < len(rns) && f(rns[i1]); i1++ {
	}
	for ; i2 > i1 && f(rns[i2]); i2-- {
	}
	return string(rns[i1 : i2+1]), string(rns[:i1]), string(rns[i2+1:])
}

//TrimQuote return word without lt quotes
func TrimQuote(w string) string {
	if w == "" {
		return w
	}
	rns := []rune(w)
	i1, i2 := 0, len(rns)-1
	if rns[i1] == StartQuote {
		i1++
	}
	if rns[i2] == EndQuote {
		i2--
	}
	return string(rns[i1 : i2+1])
}

//IsPunct return true is rune is punctuation
func IsPunct(r rune) bool {
	_, ok := Punctuations[r]
	return ok
}

//IsParentheses return true is rune is parentheses
func IsParentheses(r rune) bool {
	_, ok := parenthesesSymbols[r]
	return ok
}

//IsAllSep return true is rune is parentheses or punctuation
func IsAllSep(r rune) bool {
	return IsParentheses(r) || IsPunct(r)
}
