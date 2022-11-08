# Amazon Selling Partner API
Golang client implementation to the [Amazon SP-API](https://developer-docs.amazon.com/sp-api/).

## Installation
```shell
go get -u https://github.com/fond-of-vertigo/amazon-sp-api
```

## API-Endpoints coverage

- [ ] Authorization 
- [ ] Catalog 
- [ ] Easy Ship 
- [ ] Fulfillment by Amazon (FBA)
- [ ] Feeds
- [ ] Finances
- [ ] Fulfillment Inbound
- [ ] Fulfillment Outbound
- [ ] Listings
- [ ] Merchant Fulfillment
- [ ] Messaging
- [ ] Notifications
- [ ] Orders
- [ ] Product Fees
- [ ] Product Pricing
- [x] [Reports](https://developer-docs.amazon.com/sp-api/docs/reports-api-v2021-06-30-reference)
- [ ] Sales
- [ ] Sellers
- [ ] Service
- [ ] Shipment
- [ ] Solicitations
- [x] [Tokens](https://developer-docs.amazon.com/sp-api/docs/tokens-api-v2021-03-01-reference)
- [ ] Uploads

## CreateReport example
The following code
- Creates a new selling partner client
- Sends a request to create a new report 


```go
package main

import (
	"github.com/fond-of-vertigo/amazon-sp-api"
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"github.com/fond-of-vertigo/amazon-sp-api/apis/reports"
	"github.com/fond-of-vertigo/logger"
)
func main() {
	c := Config{
		ClientID:           "EXAMPLE_CLIENTID",
		ClientSecret:       "EXAMPLE_SECRET",
		RefreshToken:       "EXAMPLE_REFRESHTOKEN",
		IAMUserAccessKeyID: "EXAMPLE_ACCESSKEY",
		IAMUserSecretKey:   "EXAMPLE_SECRET",
		Region:             "eu-west-1",
		RoleArn:            "EXAMPLE_ROLE",
		Endpoint:           "https://sellingpartnerapi-eu.amazon.com",
		Log:                logger.New(logger.LvlDebug),
	}
	sp, err := amazon-sp-api.NewSellingPartnerClient(c)
	if err != nil {
		panic(err)
	}
	defer sp.Close()

	now := time.Now()
	from := now.Add(-24 * time.Hour * 7)
	spec := reports.CreateReportSpecification{
		ReportType:     "GET_AMAZON_FULFILLED_SHIPMENTS_DATA_INVOICING",
		DataStartTime:  apis.JsonTimeISO8601{Time: from},
		DataEndTime:    apis.JsonTimeISO8601{Time: now},
		MarketplaceIDs: []reports.MarketplaceID{reports.MarketplaceIDGermany},
	}

	resp, err := sp.Report.CreateReport(spec)
	if err != nil {
		panic(err)
	}

	println(resp.ReportID)
	
}
```