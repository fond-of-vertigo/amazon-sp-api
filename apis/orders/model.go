package orders

import "time"

// IncludedData represents the datasets that can be included in a GetOrder response.
type IncludedData string

const (
	IncludedDataBuyer        IncludedData = "BUYER"
	IncludedDataRecipient    IncludedData = "RECIPIENT"
	IncludedDataProceeds     IncludedData = "PROCEEDS"
	IncludedDataExpense      IncludedData = "EXPENSE"
	IncludedDataPromotion    IncludedData = "PROMOTION"
	IncludedDataCancellation IncludedData = "CANCELLATION"
	IncludedDataFulfillment  IncludedData = "FULFILLMENT"
	IncludedDataPackages     IncludedData = "PACKAGES"
)

// GetOrderResponse is the response for the GetOrder operation.
type GetOrderResponse struct {
	Order Order `json:"order"`
}

// Order contains comprehensive information about a customer order.
type Order struct {
	OrderID          string            `json:"orderId"`
	OrderAliases     []Alias           `json:"orderAliases,omitempty"`
	CreatedTime      time.Time         `json:"createdTime"`
	LastUpdatedTime  time.Time         `json:"lastUpdatedTime"`
	Programs         []string          `json:"programs,omitempty"`
	AssociatedOrders []AssociatedOrder  `json:"associatedOrders,omitempty"`
	SalesChannel     SalesChannel      `json:"salesChannel"`
	Buyer            *Buyer            `json:"buyer,omitempty"`
	Recipient        *Recipient        `json:"recipient,omitempty"`
	Proceeds         *OrderProceeds    `json:"proceeds,omitempty"`
	Fulfillment      *OrderFulfillment `json:"fulfillment,omitempty"`
	OrderItems       []OrderItem       `json:"orderItems"`
	Packages         []OrderPackage    `json:"packages,omitempty"`
}

// Alias is an alternative identifier for an order.
type Alias struct {
	AliasID   string `json:"aliasId"`
	AliasType string `json:"aliasType"`
}

// AssociatedOrder represents a related order (e.g. replacement or exchange).
type AssociatedOrder struct {
	OrderID         string `json:"orderId,omitempty"`
	AssociationType string `json:"associationType,omitempty"`
}

// SalesChannel contains information about where the order was placed.
type SalesChannel struct {
	ChannelName     string `json:"channelName"`
	MarketplaceID   string `json:"marketplaceId,omitempty"`
	MarketplaceName string `json:"marketplaceName,omitempty"`
}

// Buyer contains information about the customer who purchased the order.
type Buyer struct {
	BuyerName                string `json:"buyerName,omitempty"`
	BuyerEmail               string `json:"buyerEmail,omitempty"`
	BuyerCompanyName         string `json:"buyerCompanyName,omitempty"`
	BuyerPurchaseOrderNumber string `json:"buyerPurchaseOrderNumber,omitempty"`
}

// Recipient contains information about the delivery recipient.
type Recipient struct {
	DeliveryAddress    *CustomerAddress    `json:"deliveryAddress,omitempty"`
	DeliveryPreference *DeliveryPreference `json:"deliveryPreference,omitempty"`
}

// CustomerAddress is the physical address of the customer.
type CustomerAddress struct {
	Name             string                 `json:"name,omitempty"`
	CompanyName      string                 `json:"companyName,omitempty"`
	AddressLine1     string                 `json:"addressLine1,omitempty"`
	AddressLine2     string                 `json:"addressLine2,omitempty"`
	AddressLine3     string                 `json:"addressLine3,omitempty"`
	City             string                 `json:"city,omitempty"`
	DistrictOrCounty string                 `json:"districtOrCounty,omitempty"`
	StateOrRegion    string                 `json:"stateOrRegion,omitempty"`
	Municipality     string                 `json:"municipality,omitempty"`
	PostalCode       string                 `json:"postalCode,omitempty"`
	CountryCode      string                 `json:"countryCode,omitempty"`
	Phone            string                 `json:"phone,omitempty"`
	ExtendedFields   *AddressExtendedFields `json:"extendedFields,omitempty"`
	AddressType      string                 `json:"addressType,omitempty"`
}

// AddressExtendedFields contains additional address fields used in some countries (e.g. Brazil).
type AddressExtendedFields struct {
	StreetName   string `json:"streetName,omitempty"`
	StreetNumber string `json:"streetNumber,omitempty"`
	Complement   string `json:"complement,omitempty"`
	Neighborhood string `json:"neighborhood,omitempty"`
}

// DeliveryPreference contains delivery instructions and preferences.
type DeliveryPreference struct {
	DropOffLocation      string                        `json:"dropOffLocation,omitempty"`
	AddressInstruction   string                        `json:"addressInstruction,omitempty"`
	DeliveryTime         *DeliveryTime                 `json:"deliveryTime,omitempty"`
	DeliveryCapabilities []PreferredDeliveryCapability  `json:"deliveryCapabilities,omitempty"`
}

// PreferredDeliveryCapability is a delivery capability at the shipping address.
type PreferredDeliveryCapability = string

// DeliveryTime contains business hours and exception dates for delivery scheduling.
type DeliveryTime struct {
	BusinessHours  []BusinessHour  `json:"businessHours,omitempty"`
	ExceptionDates []ExceptionDate `json:"exceptionDates,omitempty"`
}

// BusinessHour defines operating hours for a specific day of the week.
type BusinessHour struct {
	DayOfWeek   string       `json:"dayOfWeek,omitempty"`
	TimeWindows []TimeWindow `json:"timeWindows,omitempty"`
}

// ExceptionDate represents a date when normal business hours are modified.
type ExceptionDate struct {
	ExceptionDate     string       `json:"exceptionDate,omitempty"`
	ExceptionDateType string       `json:"exceptionDateType,omitempty"`
	TimeWindows       []TimeWindow `json:"timeWindows,omitempty"`
}

// TimeWindow represents a time period within a day.
type TimeWindow struct {
	StartTime *HourMinute `json:"startTime,omitempty"`
	EndTime   *HourMinute `json:"endTime,omitempty"`
}

// HourMinute represents a time of day.
type HourMinute struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

// OrderProceeds contains financial information about the order.
type OrderProceeds struct {
	GrandTotal *Money `json:"grandTotal,omitempty"`
}

// Money represents a monetary amount with currency.
type Money struct {
	Amount       string `json:"amount"`
	CurrencyCode string `json:"currencyCode"`
}

// OrderFulfillment contains information about order processing and shipping.
type OrderFulfillment struct {
	FulfillmentStatus       string         `json:"fulfillmentStatus"`
	FulfilledBy             string         `json:"fulfilledBy,omitempty"`
	FulfillmentServiceLevel string         `json:"fulfillmentServiceLevel,omitempty"`
	ShipByWindow            *DateTimeRange `json:"shipByWindow,omitempty"`
	DeliverByWindow         *DateTimeRange `json:"deliverByWindow,omitempty"`
}

// DateTimeRange represents a time period with start and end boundaries.
type DateTimeRange struct {
	EarliestDateTime *time.Time `json:"earliestDateTime,omitempty"`
	LatestDateTime   *time.Time `json:"latestDateTime,omitempty"`
}

// OrderItem contains information about a single product within an order.
type OrderItem struct {
	OrderItemID     string           `json:"orderItemId"`
	QuantityOrdered int              `json:"quantityOrdered"`
	Measurement     *Measurement     `json:"measurement,omitempty"`
	Programs        []string         `json:"programs,omitempty"`
	Product         ItemProduct      `json:"product"`
	Proceeds        *ItemProceeds    `json:"proceeds,omitempty"`
	Expense         *ItemExpense     `json:"expense,omitempty"`
	Promotion       *ItemPromotion   `json:"promotion,omitempty"`
	Cancellation    *ItemCancellation `json:"cancellation,omitempty"`
	Fulfillment     *ItemFulfillment `json:"fulfillment,omitempty"`
}

// Measurement represents a unit of measure and value.
type Measurement struct {
	Unit  string  `json:"unit,omitempty"`
	Value float64 `json:"value,omitempty"`
}

// ItemProduct contains product information for an order item.
type ItemProduct struct {
	ASIN          string             `json:"asin,omitempty"`
	Title         string             `json:"title,omitempty"`
	SellerSKU     string             `json:"sellerSku,omitempty"`
	Condition     *ItemCondition     `json:"condition,omitempty"`
	Price         *ItemPrice         `json:"price,omitempty"`
	SerialNumbers []string           `json:"serialNumbers,omitempty"`
	Customization *ItemCustomization `json:"customization,omitempty"`
}

// ItemCondition describes the condition of an item.
type ItemCondition struct {
	ConditionType    string `json:"conditionType,omitempty"`
	ConditionSubtype string `json:"conditionSubtype,omitempty"`
	ConditionNote    string `json:"conditionNote,omitempty"`
}

// ItemPrice contains pricing information for a product.
type ItemPrice struct {
	UnitPrice        *Money `json:"unitPrice,omitempty"`
	PriceDesignation string `json:"priceDesignation,omitempty"`
}

// ItemCustomization contains customization information.
type ItemCustomization struct {
	CustomizedURL string `json:"customizedUrl,omitempty"`
}

// ItemProceeds contains financial information for an order item.
type ItemProceeds struct {
	ProceedsTotal *Money                  `json:"proceedsTotal,omitempty"`
	Breakdowns    []ItemProceedsBreakdown `json:"breakdowns,omitempty"`
}

// ItemProceedsBreakdown contains a detailed proceeds breakdown.
type ItemProceedsBreakdown struct {
	Type               string                        `json:"type,omitempty"`
	Subtotal           *Money                        `json:"subtotal,omitempty"`
	DetailedBreakdowns []ItemProceedsDetailedBreakdown `json:"detailedBreakdowns,omitempty"`
}

// ItemProceedsDetailedBreakdown is a granular breakdown of a proceeds subtotal.
type ItemProceedsDetailedBreakdown struct {
	Subtype string `json:"subtype,omitempty"`
	Value   *Money `json:"value,omitempty"`
}

// ItemExpense contains expense information for an order item.
type ItemExpense struct {
	PointsCost *ItemPointsCost `json:"pointsCost,omitempty"`
}

// ItemPointsCost contains the cost of points for an item.
type ItemPointsCost struct {
	PointsGranted *PointsGranted `json:"pointsGranted,omitempty"`
}

// PointsGranted contains information about Amazon Points awarded.
type PointsGranted struct {
	PointsNumber       int    `json:"pointsNumber,omitempty"`
	PointsMonetaryValue *Money `json:"pointsMonetaryValue,omitempty"`
}

// ItemPromotion contains promotion details for an order item.
type ItemPromotion struct {
	Breakdowns []ItemPromotionBreakdown `json:"breakdowns,omitempty"`
}

// ItemPromotionBreakdown contains details about a specific promotion.
type ItemPromotionBreakdown struct {
	PromotionID string `json:"promotionId,omitempty"`
}

// ItemCancellation contains cancellation information for an order item.
type ItemCancellation struct {
	CancellationRequest *ItemCancellationRequest `json:"cancellationRequest,omitempty"`
}

// ItemCancellationRequest details a cancellation request for an order item.
type ItemCancellationRequest struct {
	Requester    string `json:"requester,omitempty"`
	CancelReason string `json:"cancelReason,omitempty"`
}

// ItemFulfillment contains fulfillment information for an order item.
type ItemFulfillment struct {
	QuantityFulfilled   int           `json:"quantityFulfilled,omitempty"`
	QuantityUnfulfilled int           `json:"quantityUnfulfilled,omitempty"`
	Picking             *ItemPicking  `json:"picking,omitempty"`
	Packing             *ItemPacking  `json:"packing,omitempty"`
	Shipping            *ItemShipping `json:"shipping,omitempty"`
}

// ItemPicking contains warehouse picking process information.
type ItemPicking struct {
	SubstitutionPreference *ItemSubstitutionPreference `json:"substitutionPreference,omitempty"`
}

// ItemSubstitutionPreference describes substitution preferences for an item.
type ItemSubstitutionPreference struct {
	SubstitutionType    string                   `json:"substitutionType"`
	SubstitutionOptions []ItemSubstitutionOption  `json:"substitutionOptions,omitempty"`
}

// ItemSubstitutionOption is an alternative product for substitution.
type ItemSubstitutionOption struct {
	ASIN            string       `json:"asin,omitempty"`
	QuantityOrdered int          `json:"quantityOrdered,omitempty"`
	SellerSKU       string       `json:"sellerSku,omitempty"`
	Title           string       `json:"title,omitempty"`
	Measurement     *Measurement `json:"measurement,omitempty"`
}

// ItemPacking contains packaging information for an order item.
type ItemPacking struct {
	GiftOption *GiftOption `json:"giftOption,omitempty"`
}

// GiftOption contains gift wrapping and messaging options.
type GiftOption struct {
	GiftMessage   string `json:"giftMessage,omitempty"`
	GiftWrapLevel string `json:"giftWrapLevel,omitempty"`
}

// ItemShipping contains shipping and delivery information for an order item.
type ItemShipping struct {
	ScheduledDeliveryWindow *DateTimeRange              `json:"scheduledDeliveryWindow,omitempty"`
	ShippingConstraints     *ItemShippingConstraints     `json:"shippingConstraints,omitempty"`
	InternationalShipping   *ItemInternationalShipping   `json:"internationalShipping,omitempty"`
}

// ItemShippingConstraints contains shipping requirements and restrictions.
type ItemShippingConstraints struct {
	PalletDelivery                 string `json:"palletDelivery,omitempty"`
	CashOnDelivery                 string `json:"cashOnDelivery,omitempty"`
	SignatureConfirmation          string `json:"signatureConfirmation,omitempty"`
	RecipientIdentityVerification  string `json:"recipientIdentityVerification,omitempty"`
	RecipientAgeVerification       string `json:"recipientAgeVerification,omitempty"`
}

// ItemInternationalShipping contains cross-border shipping information.
type ItemInternationalShipping struct {
	IOSSNumber string `json:"iossNumber,omitempty"`
}

// OrderPackage contains information about a physical shipping package.
type OrderPackage struct {
	PackageReferenceID string         `json:"packageReferenceId"`
	CreatedTime        *time.Time     `json:"createdTime,omitempty"`
	PackageStatus      *PackageStatus `json:"packageStatus,omitempty"`
	Carrier            string         `json:"carrier,omitempty"`
	ShipTime           *time.Time     `json:"shipTime,omitempty"`
	ShippingService    string         `json:"shippingService,omitempty"`
	TrackingNumber     string         `json:"trackingNumber,omitempty"`
	ShipFromAddress    *MerchantAddress `json:"shipFromAddress,omitempty"`
	PackageItems       []PackageItem  `json:"packageItems,omitempty"`
}

// PackageStatus contains status and tracking information for a package.
type PackageStatus struct {
	Status         string `json:"status"`
	DetailedStatus string `json:"detailedStatus,omitempty"`
}

// MerchantAddress is the physical address of the merchant.
type MerchantAddress struct {
	Name             string `json:"name,omitempty"`
	AddressLine1     string `json:"addressLine1,omitempty"`
	AddressLine2     string `json:"addressLine2,omitempty"`
	AddressLine3     string `json:"addressLine3,omitempty"`
	City             string `json:"city,omitempty"`
	DistrictOrCounty string `json:"districtOrCounty,omitempty"`
	StateOrRegion    string `json:"stateOrRegion,omitempty"`
	Municipality     string `json:"municipality,omitempty"`
	PostalCode       string `json:"postalCode,omitempty"`
	CountryCode      string `json:"countryCode,omitempty"`
}

// PackageItem represents an individual order item within a shipping package.
type PackageItem struct {
	OrderItemID       string   `json:"orderItemId"`
	Quantity          int      `json:"quantity"`
	TransparencyCodes []string `json:"transparencyCodes,omitempty"`
}
