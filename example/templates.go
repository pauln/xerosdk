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
<p><a href="/organisations"/>Organisations</p>
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

var organisationsTemplate = `
{{range .Organisations}}
	<p>-- APIKey:{{.APIKey}} | Name:{{.Name}} | LegalName:{{.LegalName}}</p>
	<p>-- PaysTax:{{.PaysTax}} | Version:{{.Version}} | OrganisationType:{{.OrganisationType}}</p>
	<p>-- BaseCurrency:{{.BaseCurrency}} | CountryCode:{{.CountryCode}} | IsDemoCompany:{{.IsDemoCompany}}</p>
	<p>-- OrganisationStatus:{{.OrganisationStatus}} | RegistrationNumber:{{.RegistrationNumber}} | TaxNumber:{{.TaxNumber}}</p>
	<p>-- FinancialYearEndDay:{{.FinancialYearEndDay}} | FinancialYearEndMonth:{{.FinancialYearEndMonth}} | SalesTaxBasis:{{.SalesTaxBasis}}</p>
	<p>-- SalesTaxPeriod:{{.SalesTaxPeriod}} | DefaultSalesTax:{{.DefaultSalesTax}} | DefaultPurchasesTax:{{.DefaultPurchasesTax}}</p>
	<p>-- PeriodLockDate:{{.PeriodLockDate}} | EndOfYearLockDate:{{.EndOfYearLockDate}} | CreatedDateUTC:{{.CreatedDateUTC}}</p>
	<p>-- Timezone:{{.Timezone}} | OrganisationEntityType:{{.OrganisationEntityType}} | OrganisationID:{{.OrganisationID}}</p>
	<p>-- ShortCode:{{.ShortCode}} | LineOfBusiness:{{.LineOfBusiness}} | Addresses:{{.Addresses}}</p>
	<p>-- Phones:{{.Phones}} | ExternalLinks:{{.ExternalLinks}}</p>
	<br>
{{end}}
`
