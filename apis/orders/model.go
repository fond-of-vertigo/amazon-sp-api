package orders

type ResponsibleParty string
type DeemedResellerCategory string
type Model string

const (
	// MarketplaceFacilitator Tax is withheld and remitted to the taxing authority by Amazon on behalf of the seller.
	MarketplaceFacilitator Model = "MarketplaceFacilitator"

	// IOSS (Import one stop shop). The item being purchased is not held in the EU for shipment.
	IOSS DeemedResellerCategory = "IOSS"
	// UOSS (Union one stop shop). The item being purchased is held in the EU for shipment.
	UOSS DeemedResellerCategory = "UOSS"

	AmazonServicesInc ResponsibleParty = "Amazon Services, Inc."
)

// BuyerRequestedCancel provides information about whether or not a buyer requested cancellation.
type BuyerRequestedCancel struct {
	// When true, the buyer has requested cancellation.
	IsBuyerRequestedCancel *bool
	// The reason that the buyer requested cancellation.
	BuyerCancelReason *string
}

// BuyerCustomizedInfoDetail for custom orders from the Amazon Custom program.
type BuyerCustomizedInfoDetail struct {
	// The location of a zip file containing Amazon Custom data.
	CustomizedURL *string
}

// ItemBuyerInfo is a single item's buyer information.
type ItemBuyerInfo struct {
	// Buyer information for custom orders from the Amazon Custom program.
	BuyerCustomizedInfo *BuyerCustomizedInfoDetail
	// The gift wrap price of the item.
	GiftWrapPrice *Money
	// The tax on the gift wrap price.
	GiftWrapTax *Money
	// A gift message provided by the buyer.
	GiftMessageText *string
	// The gift wrap level specified by the buyer.
	GiftWrapLevel *string
}

// TaxCollection provides information about withheld taxes.
type TaxCollection struct {
	// The tax collection model applied to the item.
	Model *Model
	// The party responsible for withholding the taxes and remitting them to the taxing authority.
	ResponsibleParty *ResponsibleParty
}

// Money is the monetary value of the order.
type Money struct {
	// The three-digit currency code. In ISO 4217 format.
	CurrencyCode *string
	// The currency amount.
	Amount *string
}

// PointsGrantedDetail represents the number of Amazon Points offered with the purchase of an item, and their monetary value.
type PointsGrantedDetail struct {
	// The number of Amazon Points granted with the purchase of an item.
	PointsNumber *int
	// The monetary value of the Amazon Points granted.
	PointsMonetaryValue *Money
}

// ProductInfoDetail on the number of items.
type ProductInfoDetail struct {
	// The total number of items that are included in the ASIN.
	NumberOfItems *int
}

// OrderItem is a single order item.
type OrderItem struct {
	// The Amazon Standard Identification Number (ASIN) of the item.
	ASIN string
	// The seller stock keeping unit (SKU) of the item.
	SellerSKU *string
	// An Amazon-defined order item identifier.
	OrderItemId string
	// The name of the item.
	Title *string
	// The number of items in the order.
	QuantityOrdered int
	// The number of items shipped.
	QuantityShipped *int
	// Product information for the item.
	ProductInfo *ProductInfoDetail
	// The number and value of Amazon Points granted with the purchase of an item.
	PointsGranted *PointsGrantedDetail
	// The selling price of the order item. Note that an order item is an item and a quantity.
	// This means that the value of ItemPrice is equal to the selling price of the item multiplied by the quantity ordered.
	// Note that ItemPrice excludes ShippingPrice and GiftWrapPrice.
	ItemPrice *Money
	// The shipping price of the item.
	ShippingPrice *Money
	// The tax on the item price.
	ItemTax *Money
	// The tax on the shipping price.
	ShippingTax *Money
	// The discount on the shipping price.
	ShippingDiscount *Money
	// The tax on the discount on the shipping price.
	ShippingDiscountTax *Money
	// The total of all promotional discounts in the offer.
	PromotionDiscount *Money
	// The tax on the total of all promotional discounts in the offer.
	PromotionDiscountTax *Money
	// A list of promotion identifiers provided by the seller when the promotions were created.
	PromotionIds []string
	// The fee charged for COD service.
	CODFee *Money
	// The discount on the COD fee.
	CODFeeDiscount *Money
	// When true, the item is a gift.
	IsGift *bool
	// The condition of the item as described by the seller.
	ConditionNote *string
	// The condition of the item.
	// Possible values: New, Used, Collectible, Refurbished, Preorder, Club.
	ConditionId *string
	// The subcondition of the item.
	// Possible values: New, Mint, Very Good, Good, Acceptable, Poor, Club, OEM, Warranty, Refurbished Warranty, Refurbished, Open Box, Any, Other.
	ConditionSubtypeId *string
	// The start date of the scheduled delivery window in the time zone of the order destination. In ISO 8601 date time format.
	ScheduledDeliveryStartDate *string
	// The end date of the scheduled delivery window in the time zone of the order destination. In ISO 8601 date time format.
	ScheduledDeliveryEndDate *string
	// Indicates that the selling price is a special price that is available only for Amazon Business orders. For more information about the Amazon Business Seller Program, see the Amazon Business website.
	// Possible values: BusinessPrice - A special price that is available only for Amazon Business orders.
	PriceDesignation *string
	// Information about withheld taxes.
	TaxCollection *TaxCollection
	// When true, the product type for this item has a serial number.
	// Returned only for Amazon Easy Ship orders.
	SerialNumberRequired *bool
	// When true, transparency codes are required.
	IsTransparency *bool
	// The IOSS number for the marketplace. Sellers shipping to the European Union (EU) from outside of the EU must provide this IOSS number to their carrier when Amazon has collected the VAT on the sale.
	IOSSNumber *string
	// The store chain store identifier. Linked to a specific store in a store chain.
	StoreChainStoreId *string
	// The category of deemed reseller. This applies to selling partners that are not based in the EU and is used to help them meet the VAT Deemed Reseller tax laws in the EU and UK.
	DeemedResellerCategory *DeemedResellerCategory
	// A single item's buyer information.
	BuyerInfo *ItemBuyerInfo
	// Information about whether or not a buyer requested cancellation.
	BuyerRequestedCancel *BuyerRequestedCancel
}

// The OrderItemsList along with the order ID.
type OrderItemsList struct {
	// A list of order items.
	OrderItems []OrderItem
	// When present and not empty, pass this string token in the next request to return the next response page.
	NextToken *string
	// An Amazon-defined order identifier, in 3-7-7 format.
	AmazonOrderId string
}

// The GetOrderItemsResponse schema for the getOrderItems operation.
type GetOrderItemsResponse struct {
	OrderItemsList
}
