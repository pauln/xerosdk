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
<p><a href="/accounts"/>Accounts</p>
<p><a href="/bankTransactions"/>BankTransactions</p>
<p><a href="/refresh"/>Refresh</p>`

var contactsTemplate = `
{{range .Contacts}}
	<p>--  <b>ContactID:</b>{{.ContactID}}  |  <b>ContactNumber:</b>{{.ContactNumber}}  |  <b>AccountNumber:</b>{{.AccountNumber}}</p>
	<p>--  <b>ContactStatus:</b>{{.ContactStatus}}  |  <b>Name:</b>{{.Name}}  |  <b>FirstName:</b>{{.FirstName}}</p>
	<p>--  <b>LastName:</b>{{.LastName}}  |  <b>EmailAddress:</b>{{.EmailAddress}}  |  <b>SkypeUserName:</b>{{.SkypeUserName}}</p>
	<p>--  <b>ContactPersons:</b>{{.ContactPersons}}  |  <b>ContactNumber:</b>{{.ContactNumber}}  |  <b>AccountNumber:</b>{{.AccountNumber}}</p>
	<p>--  <b>BankAccountDetails:</b>{{.BankAccountDetails}}  |  <b>TaxNumber:</b>{{.TaxNumber}}  |  <b>AccountsReceivableTaxType:</b>{{.AccountsReceivableTaxType}}</p>
	<p>--  <b>AccountsPayableTaxType:</b>{{.AccountsPayableTaxType}}  |  <b>Addresses:</b>{{.Addresses}}  |  <b>Phones:</b>{{.Phones}}</p>
	<p>--  <b>IsSupplier:</b>{{.IsSupplier}}  |  <b>IsCustomer:</b>{{.IsCustomer}}  |  <b>DefaultCurrency:</b>{{.DefaultCurrency}}</p>
	<p>--  <b>XeroNetworkKey:</b>{{.XeroNetworkKey}}  |  <b>SalesDefaultAccountCode:</b>{{.SalesDefaultAccountCode}}  |  <b>PurchasesDefaultAccountCode:</b>{{.PurchasesDefaultAccountCode}}</p>
	<p>--  <b>SalesTrackingCategories:</b>{{.SalesTrackingCategories}}  |  <b>PurchasesTrackingCategories:</b>{{.PurchasesTrackingCategories}}  |  <b>TrackingCategoryName:</b>{{.TrackingCategoryName}}</p>
	<p>--  <b>TrackingCategoryOption:</b>{{.TrackingCategoryOption}}  |  <b>UpdatedDateUTC:</b>{{.UpdatedDateUTC}}  |  <b>ContactGroups:</b>{{.ContactGroups}}</p>
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

var accountsTemplate = `
{{range .Accounts}}
	<p><-- Code:{{.Code}} | Name:{{.Name}} | Type:{{.Type}}/p>
	<p><-- BankAccountNumber:{{.BankAccountNumber}} | Status:{{.Status}} | Description:{{.Description}}/p>
	<p><-- BankAccountType:{{.BankAccountType}} | CurrencyCode:{{.CurrencyCode}} | TaxType:{{.TaxType}}/p>
	<p><-- EnablePaymentsToAccount:{{.EnablePaymentsToAccount}} | ShowInExpenseClaims:{{.ShowInExpenseClaims}} | AccountID:{{.AccountID}}/p>
	<p><-- Class:{{.Class}} | SystemAccount:{{.SystemAccount}} | ReportingCode:{{.ReportingCode}}/p>
	<p><-- ReportingCodeName:{{.ReportingCodeName}} | HasAttachments:{{.HasAttachments}} | UpdatedDateUTC:{{.UpdatedDateUTC}}/p>
	<br>
{{end}}
`

var bankTransactionsTemplate = `
{{range .BankTransactions}}
	<p>-- Type:{{.Type}} | Contact:{{.Contact}} | LineItems:{{.LineItems}}</p>
	<p>-- IsReconciled:{{.IsReconciled}} | Date:{{.Date}} | Reference:{{.Reference}}</p>
	<p>-- CurrencyCode:{{.CurrencyCode}} | CurrencyRate:{{.CurrencyRate}} | URL:{{.URL}}</p>
	<p>-- Status:{{.Status}} | LineAmountTypes:{{.LineAmountTypes}} | SubTotal:{{.SubTotal}}</p>
	<p>-- BankAccount:{{.BankAccount}} | PrepaymentID:{{.PrepaymentID}} | OverpaymentID:{{.OverpaymentID}}</p>
	<p>-- UpdatedDateUTC:{{.UpdatedDateUTC}} | HasAttachments:{{.HasAttachments}}</p>
	<br>
{{end}}
`
