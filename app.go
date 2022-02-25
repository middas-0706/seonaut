package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/stripe/stripe-go/v72"
)

type App struct {
	config    *Config
	datastore *datastore
	cookie    *sessions.CookieStore
	sanitizer *bluemonday.Policy
	renderer  *Renderer
}

func NewApp(configPath string) *App {
	var err error
	var app App

	app.config, err = NewConfig(configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	app.datastore, err = NewDataStore(app.config.DB)

	if err != nil {
		log.Fatalf("Error creating new datastore: %v\n", err)
	}

	app.cookie = sessions.NewCookieStore([]byte("SESSION_ID"))
	app.sanitizer = bluemonday.StrictPolicy()

	app.renderer = NewRenderer()

	return &app
}

func (app *App) Run() {
	stripe.Key = app.config.Stripe.Secret

	stripe.SetAppInfo(&stripe.AppInfo{
		Name:    "stripe-samples/checkout-single-subscription",
		Version: "0.0.1",
		URL:     "https://github.com/stripe-samples/checkout-single-subscription",
	})

	// Static
	fileServer := http.FileServer(http.Dir("./web/static"))
	http.Handle("/resources/", http.StripPrefix("/resources", fileServer))
	http.Handle("/robots.txt", fileServer)
	http.Handle("/favicon.ico", fileServer)

	// App
	http.HandleFunc("/", app.requireAuth(app.serveHome))
	http.HandleFunc("/new-project", app.requireAuth(app.serveProjectAdd))
	http.HandleFunc("/crawl", app.requireAuth(app.serveCrawl))
	http.HandleFunc("/issues", app.requireAuth(app.serveIssues))
	http.HandleFunc("/issues/view", app.requireAuth(app.serveIssuesView))
	http.HandleFunc("/download", app.requireAuth(app.serveDownloadCSV))
	http.HandleFunc("/sitemap", app.requireAuth(app.serveSitemap))
	http.HandleFunc("/resources", app.requireAuth(app.serveResourcesView))
	http.HandleFunc("/signout", app.requireAuth(app.serveSignout))
	http.HandleFunc("/signup", app.serveSignup)
	http.HandleFunc("/signin", app.serveSignin)

	// Stripe
	http.HandleFunc("/upgrade", app.requireAuth(app.upgrade))
	http.HandleFunc("/create-checkout-session", app.requireAuth(app.handleCreateCheckoutSession))
	http.HandleFunc("/checkout-session", app.requireAuth(app.handleCheckoutSession))
	http.HandleFunc("/config", app.requireAuth(app.handleConfig))
	http.HandleFunc("/manage", app.requireAuth(app.handleManageAccount))
	http.HandleFunc("/canceled", app.requireAuth(app.handleCanceled))
	http.HandleFunc("/customer-portal", app.requireAuth(app.handleCustomerPortal))
	http.HandleFunc("/webhook", app.handleWebhook)

	fmt.Printf("Starting server at %s on port %d...\n", app.config.Server, app.config.ServerPort)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", app.config.Server, app.config.ServerPort), nil)
	if err != nil {
		log.Fatal(err)
	}
}
