package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func loadTransforms() (string, []string) {
	contents, err := ioutil.ReadFile("transforms.txt")
	if err != nil {
		log.Fatal(err)
	}
	transforms := strings.Split(string(contents), "\n")
	return transforms[0], transforms[1:]
}

func main() {
	token, transforms := loadTransforms()
	rand.Seed(time.Now().UTC().UnixNano())
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		t := string(transforms[rand.Intn(len(transforms))])
		fmt.Println(strings.Replace(t, token, s.Text(), -1))
	}
}
