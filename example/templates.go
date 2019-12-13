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
<p><a href="/contacts"/>Contacts</p>`

var contactsTemplate = `<h1>Contacts from Tenant: {{.TenantOne}}</h1>
{{range .Contacts}}
	<p>--  Name:{{.Name}}  |  Email:{{.EmailAddress}}  |  ContactID:{{.ContactID}}</p>
{{end}}
<h1>Contacts from Tenant: {{.TenantTwo}}</h1>
{{range .ContactsSecond}}
	<p>--  Name:{{.Name}}  |  Email:{{.EmailAddress}}  |  ContactID:{{.ContactID}}</p>
{{end}}`
