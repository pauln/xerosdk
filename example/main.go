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
	r.HandleFunc("/invoices", XeroInvoicesHandler)
	r.HandleFunc("/refresh", XeroRefreshTokenHandler)
	r.HandleFunc("/organisations", XeroOrganisationsHandler)
	r.HandleFunc("/accounts", XeroAccountsHandler)
	r.HandleFunc("/bankTransactions", XeroBankTransactionsHandler)
	r.HandleFunc("/bankTransfers", XeroBankTransfersHandler)
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

// XeroConnectionsHandler is the handler that will show all the granted access
// tenants
func XeroConnectionsHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	tenants, err := connection.GetTenants(c.Client(&auth.Session{
		Token:    se,
		UserID:   uuid.Nil,
		TenantID: uuid.Nil,
		Repo:     repo,
	}))
	if err != nil {
		log.Panic(err)
	}
	json.NewEncoder(w).Encode(&tenants)
}

// XeroRefreshTokenHandler is the handler that will refresh the current token
// an saved it as a current one
func XeroRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	newToken, err := c.Refresh(se)
	if err != nil {
		log.Panic(err)
	}
	repo.UpdateSession(uuid.Nil, newToken)
	http.Redirect(w, r, "/", http.StatusFound)
}

//XeroContactsHandler is the handler in where we will show all the existing contacts
// with all the tenants connected
func XeroContactsHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	cl := c.Client(&auth.Session{
		Token:    se,
		UserID:   uuid.Nil,
		TenantID: uuid.Nil,
		Repo:     repo,
	})
	contacts := []accounting.Contact{}

	tenants, err := connection.GetTenants(cl)
	if err != nil {
		log.Panic(err)
	}
	for _, tenant := range tenants {
		c, err := accounting.FindContacts(c.Client(&auth.Session{
			Token:    se,
			UserID:   uuid.Nil,
			TenantID: tenant.TenantID,
			Repo:     repo,
		}))
		if err != nil {
			log.Panic(err)
		}
		contacts = append(contacts, c.Contacts...)
	}
	t, _ := template.New("contacts").Parse(contactsTemplate)
	t.Execute(w, struct {
		Contacts []accounting.Contact
	}{
		Contacts: contacts,
	})
}

// XeroContactsCreateHandler is the handler that will create a new dummy contact
func XeroContactsCreateHandler(w http.ResponseWriter, r *http.Request) {
	se, _ := repo.GetSession(uuid.Nil)
	contactID, _ := uuid.NewV4()

	contacts := accounting.Contacts{
		Contacts: []accounting.Contact{accounting.Contact{
			Name:         "Test " + contactID.String(),
			FirstName:    "Test FirstName",
			LastName:     "Test LastName",
			EmailAddress: "Test Email " + contactID.String(),
		},
		},
	}

	tenants, err := connection.GetTenants(c.Client(&auth.Session{
		Token:  se,
		UserID: uuid.Nil,
		Repo:   repo,
	}))
	if err != nil {
		log.Panic(err)
	}
	// We asume we have at least one tenant connected
	// TODO improve that to get this information from a form
	_, err = contacts.Create(c.Client(&auth.Session{
		Token:    se,
		UserID:   uuid.Nil,
		TenantID: tenants[0].TenantID,
		Repo:     repo,
	}))
	if err != nil {
		log.Panic(err)
	} else {
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// XeroInvoicesHandler is the handler that will find all the invoices
func XeroInvoicesHandler(w http.ResponseWriter, r *http.Request) {
	invoices := []accounting.Invoice{}
	se, _ := repo.GetSession(uuid.Nil)

	tenants, err := connection.GetTenants(c.Client(&auth.Session{
		Token:  se,
		UserID: uuid.Nil,
		Repo:   repo,
	}))
	if err != nil {
		log.Panic(err)
	}
	for _, tenant := range tenants {
		i, err := accounting.FindInvoices(c.Client(&auth.Session{
			Token:    se,
			UserID:   uuid.Nil,
			TenantID: tenant.TenantID,
			Repo:     repo,
		}))
		if err != nil {
			log.Panic(err)
		}
		invoices = append(invoices, i.Invoices...)
	}
	t, _ := template.New("invoices").Parse(invoicesTemplate)
	t.Execute(w, struct {
		Invoices []accounting.Invoice
	}{
		Invoices: invoices,
	})
}

// XeroOrganisationsHandler handler will ask for all the organisations linked
// to the given user and print out in a template
func XeroOrganisationsHandler(w http.ResponseWriter, r *http.Request) {
	organisations := []accounting.Organisation{}
	se, _ := repo.GetSession(uuid.Nil)

	tenants, err := connection.GetTenants(c.Client(&auth.Session{
		Token:  se,
		UserID: uuid.Nil,
		Repo:   repo,
	}))
	if err != nil {
		log.Panic(err)
	}
	for _, tenant := range tenants {
		orgs, err := accounting.FindOrganisations(c.Client(&auth.Session{
			Token:    se,
			UserID:   uuid.Nil,
			TenantID: tenant.TenantID,
			Repo:     repo,
		}))
		if err != nil {
			log.Panic(err)
		}
		organisations = append(organisations, orgs.Organisations...)
	}
	t, _ := template.New("organisations").Parse(organisationsTemplate)
	t.Execute(w, struct {
		Organisations []accounting.Organisation
	}{
		Organisations: organisations,
	})
}

// XeroAccountsHandler handler will ask for all the accounts linked to the
// given user and print out in a template
func XeroAccountsHandler(w http.ResponseWriter, r *http.Request) {
	accounts := []accounting.Account{}
	se, _ := repo.GetSession(uuid.Nil)

	tenants, err := connection.GetTenants(c.Client(&auth.Session{
		Token:  se,
		UserID: uuid.Nil,
		Repo:   repo,
	}))
	if err != nil {
		log.Panic(err)
	}
	for _, tenant := range tenants {
		accs, err := accounting.FindAccounts(c.Client(&auth.Session{
			Token:    se,
			UserID:   uuid.Nil,
			TenantID: tenant.TenantID,
			Repo:     repo,
		}), nil)
		if err != nil {
			log.Panic(err)
		}
		accounts = append(accounts, accs.Accounts...)
	}
	t, _ := template.New("accounts").Parse(accountsTemplate)
	t.Execute(w, struct {
		Accounts []accounting.Account
	}{
		Accounts: accounts,
	})
}

// XeroBankTransactionsHandler handler will ask for all the bank transactions linked to the given
// user and print out in a template
func XeroBankTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	bankTransactions := []accounting.BankTransaction{}
	se, _ := repo.GetSession(uuid.Nil)

	tenants, err := connection.GetTenants(c.Client(&auth.Session{
		Token:  se,
		UserID: uuid.Nil,
		Repo:   repo,
	}))
	if err != nil {
		log.Panic(err)
	}
	for _, tenant := range tenants {
		bankTr, err := accounting.FindBankTransactions(c.Client(&auth.Session{
			Token:    se,
			UserID:   uuid.Nil,
			TenantID: tenant.TenantID,
			Repo:     repo,
		}), nil)
		if err != nil {
			log.Panic(err)
		}
		bankTransactions = append(bankTransactions, bankTr.BankTransactions...)
	}
	t, _ := template.New("bankTransactions").Parse(bankTransactionsTemplate)
	t.Execute(w, struct {
		BankTransactions []accounting.BankTransaction
	}{
		BankTransactions: bankTransactions,
	})
}

// XeroBankTransfersHandler handler will ask for all the bank transfers linked
// to the given user and print out in a template
func XeroBankTransfersHandler(w http.ResponseWriter, r *http.Request) {
	bankTransfers := []accounting.BankTransfer{}
	se, _ := repo.GetSession(uuid.Nil)

	tenants, err := connection.GetTenants(c.Client(&auth.Session{
		Token:  se,
		UserID: uuid.Nil,
		Repo:   repo,
	}))
	if err != nil {
		log.Panic(err)
	}
	for _, tenant := range tenants {
		bankTrns, err := accounting.FindBankTransfers(c.Client(&auth.Session{
			Token:    se,
			UserID:   uuid.Nil,
			TenantID: tenant.TenantID,
			Repo:     repo,
		}), nil)
		if err != nil {
			log.Panic(err)
		}
		bankTransfers = append(bankTransfers, bankTrns.BankTransfers...)
	}
	t, _ := template.New("bankTransfers").Parse(bankTransfersTemplate)
	t.Execute(w, struct {
		BankTransfers []accounting.BankTransfer
	}{
		BankTransfers: bankTransfers,
	})
}
