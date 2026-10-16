package brief

import (
	"strings"
	"unicode"
)

// OpenDecision is one entry of a brief's `## Open decisions`: its text, list
// marker stripped, and who it says decides it ("captain", "mate", or "" when
// the entry names neither).
type OpenDecision struct {
	Text    string
	Decides string
}

// ParseOpenDecisions is the entries of text's `## Open decisions`, nil when
// the section is missing or is `none`. text is the Mate's file or a rendered
// brief, as SectionItems takes it.
func ParseOpenDecisions(text string) []OpenDecision {
	its, ok := SectionItems(text, OpenDecisions)
	if !ok {
		return nil
	}
	var out []OpenDecision
	for _, it := range its {
		if isNone(it) {
			continue
		}
		d := OpenDecision{Text: it}
		if m := decidesValue.FindStringSubmatch(it); m != nil {
			if who := strings.ToLower(m[1]); who == Deciders[0] || who == Deciders[1] {
				d.Decides = who
			}
		}
		out = append(out, d)
	}
	return out
}

// MatchDecision is the open decision a Crew's question is about, matched
// roughly (CON-93): the entry sharing the most words with the question,
// counting words of four letters or more and leaving out the ones every
// entry carries (`decides`, `captain`, `mate`). No shared word, or a tie
// between entries that name different deciders, is no match: the caller
// says it does not know rather than guessing whose call it is.
func MatchDecision(decisions []OpenDecision, question string) (OpenDecision, bool) {
	q := matchWords(question)
	best, score, tied := -1, 0, false
	for i, d := range decisions {
		n := 0
		for w := range matchWords(d.Text) {
			if q[w] {
				n++
			}
		}
		switch {
		case n > score:
			best, score, tied = i, n, false
		case n == score && n > 0 && decisions[best].Decides != d.Decides:
			tied = true
		}
	}
	if best < 0 || tied {
		return OpenDecision{}, false
	}
	return decisions[best], true
}

var matchStopWords = map[string]bool{
	"decides": true, "captain": true, "mate": true, "should": true, "which": true,
	"what": true, "with": true, "that": true, "this": true, "from": true,
	"into": true, "when": true, "will": true, "would": true, "there": true,
	"their": true, "them": true, "then": true, "than": true, "have": true,
}

func matchWords(s string) map[string]bool {
	out := map[string]bool{}
	for _, w := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	}) {
		if len([]rune(w)) >= 4 && !matchStopWords[w] {
			out[w] = true
		}
	}
	return out
}
