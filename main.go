package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/prometheus/alertmanager/notify/webhook"
	"gopkg.in/yaml.v2"
)

const APIURL = "https://smsapi.free-mobile.fr/sendmsg"

var (
	cfg = config{}

	configPath = flag.String("config", "", "Path to config")

	baseUrl = &url.URL{}
)

type config struct {
	Sentry struct {
		Dsn string `yaml:"dsn"`
	}
	Server struct {
		Address string `yaml:"address"`
	} `yaml:"server"`
	Free struct {
		User string `yaml:"user"`
		Pass string `yaml:"pass"`
	} `yaml:"free"`
}

func paramBuilder(u *url.URL, alert webhook.Message) string {
	params := url.Values{
		"user": {cfg.Free.User},
		"pass": {cfg.Free.Pass},
		"msg": {
			"Alertmanager => " +
				"status : " + alert.Status +
				", alerts : " + strconv.Itoa(len(alert.Alerts))},
	}
	u.RawQuery = params.Encode()

	return u.String()
}

func sendSMS(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var alert webhook.Message
		err := json.NewDecoder(r.Body).Decode(&alert)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u := paramBuilder(baseUrl, alert)
		resp, err := sendSMS(u)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(resp.StatusCode)
		_, err = w.Write([]byte(http.StatusText(resp.StatusCode)))
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	default:
		_, err := fmt.Fprintf(w, "Only POST method is supported.")
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}
}

func newConfig(path string, config *config) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d := yaml.NewDecoder(f)
	if err := d.Decode(config); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	err := newConfig(*configPath, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("loaded config from file %s", *configPath)

	baseUrl, err = url.Parse(APIURL)
	if err != nil {
		log.Fatal(err)
	}

	err = sentry.Init(sentry.ClientOptions{
		Dsn: cfg.Sentry.Dsn,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("sentry initialized")

	sentryHandler := sentryhttp.New(sentryhttp.Options{})
	mux := http.NewServeMux()
	mux.HandleFunc("/", sentryHandler.HandleFunc(viewHandler))

	log.Printf("server is starting on %s", cfg.Server.Address)
	if err := http.ListenAndServe(cfg.Server.Address, mux); err != nil {
		log.Fatal(err)
	}
}
