package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/bahelms/go_blueprints/thesaurus"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := &thesaurus.BigHuge{APIKey: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		syns, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("Failed when looking for synonyms for " + word + "," + err.Error())
		}
		if len(syns) == 0 {
			log.Fatalln("Couldn't find any synonyms for " + word)
		}
		for _, syn := range syns {
			fmt.Println(syn)
		}
	}
}
