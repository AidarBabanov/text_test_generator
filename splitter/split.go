package splitter

import "strings"

type Splitter struct {
	Separator byte
	Text      string
	SubSplits []*Splitter
}

func (spl *Splitter) Size() int {
	if len(spl.SubSplits) == 0 {
		return 1
	}
	size := 0
	for _, subSpl := range spl.SubSplits {
		size += subSpl.Size()
	}
	return size
}

func (spl *Splitter) Words() []*string {
	if len(spl.SubSplits) == 0 {
		return []*string{&spl.Text}
	}
	var words []*string
	for _, subSpl := range spl.SubSplits {
		words = append(words, subSpl.Words()...)
	}
	return words
}

func (spl *Splitter) Reconfigure() string {
	if len(spl.SubSplits) == 0 {
		return spl.Text
	}
	text := ""
	for i := 0; i < len(spl.SubSplits)-1; i++ {
		text += spl.SubSplits[i].Reconfigure() + string(spl.Separator)
	}
	text += spl.SubSplits[len(spl.SubSplits)-1].Reconfigure()
	return text
}

func Split(text string, separators string) *Splitter {
	splitter := new(Splitter)
	if len(separators) == 0 {
		splitter.Text = text
		return splitter
	}
	sep := separators[0]
	splitter.Separator = sep
	sepTexts := strings.Split(text, string(sep))
	for _, sepText := range sepTexts {
		splitter.SubSplits = append(splitter.SubSplits, Split(sepText, separators[1:]))
	}
	return splitter
}

func ConvertWords(words []*string) []string {
	newWords := make([]string, len(words))
	for i := 0; i < len(words); i++ {
		newWords[i] = *words[i]
	}
	return newWords
}
