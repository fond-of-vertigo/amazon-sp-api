package constants

type ProcessingStatus string

func (s ProcessingStatus) IsDone() bool {
	return s == Done
}

type MarketplaceID string
type Region string
type Endpoint string

const (
	AccessTokenHeader = "X-Amz-Access-Token"
	ServiceExecuteAPI = "execute-api"
)

const (
	Done       ProcessingStatus = "DONE"
	Cancelled  ProcessingStatus = "CANCELLED"
	Fatal      ProcessingStatus = "FATAL"
	InProgress ProcessingStatus = "IN_PROGRESS"
	InQueue    ProcessingStatus = "IN_QUEUE"

	Canada                MarketplaceID = "A2EUQ1WTGCTBG2"
	UnitedStatesOfAmerica MarketplaceID = "ATVPDKIKX0DER"
	Mexico                MarketplaceID = "A1AM78C64UM0Y8"
	Brazil                MarketplaceID = "A2Q3Y263D00KWC"
	Spain                 MarketplaceID = "A1RKKUPIHCS9HS"
	UnitedKingdom         MarketplaceID = "A1F83G8C2ARO7P"
	France                MarketplaceID = "A13V1IB3VIYZZH"
	Belgium               MarketplaceID = "AMEN7PMS3EDWL"
	Netherlands           MarketplaceID = "A1805IZSGTT6HS"
	Germany               MarketplaceID = "A1PA6795UKMFR9"
	Italy                 MarketplaceID = "APJ6JRA9NG5V4"
	Sweden                MarketplaceID = "A2NODRKZP88ZB9"
	Poland                MarketplaceID = "A1C3SOZRARQ6R3"
	Egypt                 MarketplaceID = "ARBP9OOSHTCHU"
	Turkey                MarketplaceID = "A33AVAJ2PDY3EV"
	SaudiArabia           MarketplaceID = "A17E79C6D8DWNP"
	UnitedArabEmirates    MarketplaceID = "A2VIGQ35RCS4UG"
	India                 MarketplaceID = "A21TJRUUN4KGV"
	Singapore             MarketplaceID = "A19VAU5U5O7RUS"
	Australia             MarketplaceID = "A39IBJ37TRP1C6"
	Japan                 MarketplaceID = "A1VC38T7YXB528"

	USEast Region = "us-east-1"
	EUWest Region = "eu-west-1"
	USWest Region = "us-west-2"

	NorthAmerica Endpoint = "https://sellingpartnerapi-na.amazon.com"
	Europe       Endpoint = "https://sellingpartnerapi-eu.amazon.com"
	FarEast      Endpoint = "https://sellingpartnerapi-fe.amazon.com"
)
