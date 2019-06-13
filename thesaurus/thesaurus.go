package thesaurus

// Thesaurus does stuff
type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}
