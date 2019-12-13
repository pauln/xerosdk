package main

import (
	"context"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/gorilla/mux"
	"github.com/quickaco/xerosdk/accounting"
	"github.com/quickaco/xerosdk/auth"
	"github.com/quickaco/xerosdk/connection"

	"github.com/joho/godotenv"
)

var (
	c    *auth.Provider
	repo auth.Repository
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	config := auth.Config{
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       strings.Split(os.Getenv("SCOPES"), ","),
		RedirectURL:  os.Getenv("REDIRECT_URL"),
	}
	c = auth.NewProvider(config)
	repo = NewRepository()
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/auth/xero", StartXeroAuthHandler)
	r.HandleFunc("/auth/xero/callback", XeroAuthCallbackHandler)
	r.HandleFunc("/connections", XeroConnectionsHandler)
	r.HandleFunc("/contacts", XeroContactsHandler)
	r.HandleFunc("/contacts/create", XeroContactsCreateHandler)
	r.HandleFunc("/refresh", XeroRefreshTokenHandler)
	http.Handle("/", r)

	srv := &http.Server{
		Addr: "0.0.0.0:3000",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}

// HomeHandler will be the base handler in where we will show information about
// token and different actions you can do
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	se, _ := repo.GetSession(uuid.Nil)
	if se != nil {
		t, _ = template.New("foo").Parse(connectedTemplate)
	} else {
		t, _ = template.New("foo").Parse(indexTemplate)
	}
	t.Execute(w, se)
}

// StartXeroAuthHandler is the handler that will start the process of Auth with
// the Xero platform
func StartXeroAuthHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, c.GetAuthURL("uniq_state"), http.StatusFound)
}

// XeroAuthCallbackHandler is the handler in where we are going to receive a
// successful callback with a code that can we use to get our user token
func XeroAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := c.GetTokenFromCode(r.FormValue("code"))
	if err != nil {
		log.Panic(err)
	}
	repo.CreateSession(uuid.Nil, token)
	t, _ := template.New("connected").Parse(connectedTemplate)
	t.Execute(w, token)
}

func XeroConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	tenants, err := connection.GetTenants(c.Client(se, repo))
	if err != nil {
		log.Panic(err)
	}
	json.NewEncoder(w).Encode(&tenants)
}

func XeroRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	log.Printf("OLD TOKEN %+v", se)
	newToken, err  := c.Refresh(se)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("REFRESHED TOKEN %+v", newToken)
	repo.UpdateSession(uuid.Nil, newToken)
}

func XeroContactsHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	cl := c.Client(se, repo)
	contacts := []accounting.Contact{}

	tenants, err := connection.GetTenants(cl)
	if err != nil {
		log.Panic(err)
	}
	for _, tenant := range tenants {
		c, err := accounting.FindContacts(cl, tenant.TenantID)
		if err != nil {
			log.Panic(err)
		}
		contacts = append(contacts, c.Contacts...)
	}
	t, _ := template.New("contacts").Parse(contactsTemplate)
	t.Execute(w, struct {
		Contacts       []accounting.Contact
	}{
		Contacts:       contacts,
	})
}

func XeroContactsCreateHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	cl := c.Client(se, repo)
	contactID, _ := uuid.NewV4()

	contacts := accounting.Contacts{
		Contacts: []accounting.Contact{accounting.Contact{
			Name: "Test " + contactID.String(),
			FirstName: "Test FirstName",
			LastName: "Test LastName",
		},
		},
	}

	tenants, err := connection.GetTenants(cl)
	if err != nil {
		log.Panic(err)
	}
	// We asume we have at least one tenant connected
	// TODO improve that to get this information from a form
	_, err = contacts.Create(cl, tenants[0].TenantID)
	if err != nil {
		log.Panic(err)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
