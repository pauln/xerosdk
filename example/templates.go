package main

var indexTemplate = `<p>
		<a href="/auth/xero">
			<img src="https://developer.xero.com/static/images/documentation/connect_xero_button_blue_2x.png" alt="ConnectToXero">
		</a>
	</p>`

var connectedTemplate = `<p>AccessToken: {{.AccessToken}}</p>
<p>TokenType: {{.TokenType}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
<p>Expiry: {{.Expiry}}</p>
<p><a href="/connections"/>Connections</p>
<p><a href="/contacts"/>Contacts</p>
<p><a href="/invoices"/>Invoices</p>
<p><a href="/refresh"/>Refresh</p>`

var contactsTemplate = `
{{range .Contacts}}
	<p>--  Name:{{.Name}}  |  Email:{{.EmailAddress}}  |  ContactID:{{.ContactID}}</p>
{{end}}

<p><a href="/contacts/create">Create a dummy contact</p>`

var invoicesTemplate = `
{{range .Invoices}}
	<p>-- Type:{{.Type}} | InvoiceNumber:{{.InvoiceNumber}} | Status:{{.Status}}</p>
{{end}}`
