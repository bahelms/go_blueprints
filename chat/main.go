package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/bahelms/go_blueprints/chat/trace"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
	yaml "gopkg.in/yaml.v2"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		fp := filepath.Join("templates", t.filename)
		t.templ = template.Must(template.ParseFiles(fp))
	})
	data := map[string]interface{}{
		"host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

type config struct {
	SecurityKey    string `yaml:"security_key"`
	GoogleID       string `yaml:"googleID"`
	GoogleKey      string `yaml:"googleKey"`
	GoogleCallback string `yaml:"goggleCallback"`
}

func main() {
	addr := flag.String("addr", ":8080", "The port of the application.")
	flag.Parse()

	// setup gomniauth
	configData, _ := ioutil.ReadFile("config.yml")
	conf := config{}
	err := yaml.Unmarshal(configData, &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("config: %v\n", conf)
	gomniauth.SetSecurityKey(conf.SecurityKey)
	gomniauth.WithProviders(
		google.New(conf.GoogleID, conf.GoogleKey, conf.GoogleCallback),
	)

	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run() // start room

	// start server
	log.Println("Starting web server on", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
