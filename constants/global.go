package constants

type ProcessingStatus string
type MarketplaceID string
type Region string
type Endpoint string

const (
	AccessTokenHeader = "X-Amz-Access-Token"
	ServiceExecuteAPI = "execute-api"
)

const (
	ProcessingStatusDone       ProcessingStatus = "DONE"
	ProcessingStatusCancelled  ProcessingStatus = "CANCELLED"
	ProcessingStatusFatal      ProcessingStatus = "FATAL"
	ProcessingStatusInProgress ProcessingStatus = "IN_PROGRESS"
	ProcessingStatusInQueue    ProcessingStatus = "IN_QUEUE"

	MarketplaceIDCanada                MarketplaceID = "A2EUQ1WTGCTBG2"
	MarketplaceIDUnitedStatesOfAmerica MarketplaceID = "ATVPDKIKX0DER"
	MarketplaceIDMexico                MarketplaceID = "A1AM78C64UM0Y8"
	MarketplaceIDBrazil                MarketplaceID = "A2Q3Y263D00KWC"
	MarketplaceIDSpain                 MarketplaceID = "A1RKKUPIHCS9HS"
	MarketplaceIDUnitedKingdom         MarketplaceID = "A1F83G8C2ARO7P"
	MarketplaceIDFrance                MarketplaceID = "A13V1IB3VIYZZH"
	MarketplaceIDBelgium               MarketplaceID = "AMEN7PMS3EDWL"
	MarketplaceIDNetherlands           MarketplaceID = "A1805IZSGTT6HS"
	MarketplaceIDGermany               MarketplaceID = "A1PA6795UKMFR9"
	MarketplaceIDItaly                 MarketplaceID = "APJ6JRA9NG5V4"
	MarketplaceIDSweden                MarketplaceID = "A2NODRKZP88ZB9"
	MarketplaceIDPoland                MarketplaceID = "A1C3SOZRARQ6R3"
	MarketplaceIDEgypt                 MarketplaceID = "ARBP9OOSHTCHU"
	MarketplaceIDTurkey                MarketplaceID = "A33AVAJ2PDY3EV"
	MarketplaceIDSaudiArabia           MarketplaceID = "A17E79C6D8DWNP"
	MarketplaceIDUnitedArabEmirates    MarketplaceID = "A2VIGQ35RCS4UG"
	MarketplaceIDIndia                 MarketplaceID = "A21TJRUUN4KGV"
	MarketplaceIDSingapore             MarketplaceID = "A19VAU5U5O7RUS"
	MarketplaceIDAustralia             MarketplaceID = "A39IBJ37TRP1C6"
	MarketplaceIDJapan                 MarketplaceID = "A1VC38T7YXB528"

	RegionUSEast Region = "us-east-1"
	RegionEUWest Region = "eu-west-1"
	RegionUSWest Region = "us-west-2"

	EndpointNorthAmerica Endpoint = "https://sellingpartnerapi-na.amazon.com"
	EndpointEurope       Endpoint = "https://sellingpartnerapi-eu.amazon.com"
	EndpointFarEast      Endpoint = "https://sellingpartnerapi-fe.amazon.com"
)
