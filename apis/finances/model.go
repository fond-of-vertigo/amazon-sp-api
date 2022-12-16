package finances

import (
	"github.com/fond-of-vertigo/amazon-sp-api/apis"
	"time"
)

// AdjustmentEvent An adjustment to the seller's account.
type AdjustmentEvent struct {
	// The type of adjustment.  Possible values:  * FBAInventoryReimbursement - An FBA inventory reimbursement to a seller's account. This occurs if a seller's inventory is damaged.  * ReserveEvent - A reserve event that is generated at the time of a settlement period closing. This occurs when some money from a seller's account is held back.  * PostageBilling - The amount paid by a seller for shipping labels.  * PostageRefund - The reimbursement of shipping labels purchased for orders that were canceled or refunded.  * LostOrDamagedReimbursement - An Amazon Easy Ship reimbursement to a seller's account for a package that we lost or damaged.  * CanceledButPickedUpReimbursement - An Amazon Easy Ship reimbursement to a seller's account. This occurs when a package is picked up and the order is subsequently canceled. This value is used only in the India marketplace.  * ReimbursementClawback - An Amazon Easy Ship reimbursement clawback from a seller's account. This occurs when a prior reimbursement is reversed. This value is used only in the India marketplace.  * SellerRewards - An award credited to a seller's account for their participation in an offer in the Seller Rewards program. Applies only to the India marketplace.
	AdjustmentType   *string    `json:"AdjustmentType,omitempty"`
	PostedDate       *time.Time `json:"PostedDate,omitempty"`
	AdjustmentAmount *Currency  `json:"AdjustmentAmount,omitempty"`
	// A list of information about items in an adjustment to the seller's account.
	AdjustmentItemList []AdjustmentItem `json:"AdjustmentItemList,omitempty"`
}

// AdjustmentItem An item in an adjustment to the seller's account.
type AdjustmentItem struct {
	// Represents the number of units in the seller's inventory when the AdustmentType is FBAInventoryReimbursement.
	Quantity      *string   `json:"Quantity,omitempty"`
	PerUnitAmount *Currency `json:"PerUnitAmount,omitempty"`
	TotalAmount   *Currency `json:"TotalAmount,omitempty"`
	// The seller SKU of the item. The seller SKU is qualified by the seller's seller ID, which is included with every call to the Selling Partner API.
	SellerSKU *string `json:"SellerSKU,omitempty"`
	// A unique identifier assigned to products stored in and fulfilled from a fulfillment center.
	FnSKU *string `json:"FnSKU,omitempty"`
	// A short description of the item.
	ProductDescription *string `json:"ProductDescription,omitempty"`
	// The Amazon Standard Identification Number (ASIN) of the item.
	ASIN *string `json:"ASIN,omitempty"`
}

// Currency A currency type and amount.
type Currency struct {
	// The three-digit currency code in ISO 4217 format.
	CurrencyCode   *string  `json:"CurrencyCode,omitempty"`
	CurrencyAmount *float32 `json:"CurrencyAmount,omitempty"`
}

// AffordabilityExpenseEvent An expense related to an affordability promotion.
type AffordabilityExpenseEvent struct {
	// An Amazon-defined identifier for an order.
	AmazonOrderId *string    `json:"AmazonOrderId,omitempty"`
	PostedDate    *time.Time `json:"PostedDate,omitempty"`
	// An encrypted, Amazon-defined marketplace identifier.
	MarketplaceId *string `json:"MarketplaceId,omitempty"`
	// Indicates the type of transaction.   Possible values:  * Charge - For an affordability promotion expense.  * Refund - For an affordability promotion expense reversal.
	TransactionType *string   `json:"TransactionType,omitempty"`
	BaseExpense     *Currency `json:"BaseExpense,omitempty"`
	TaxTypeCGST     Currency  `json:"TaxTypeCGST"`
	TaxTypeSGST     Currency  `json:"TaxTypeSGST"`
	TaxTypeIGST     Currency  `json:"TaxTypeIGST"`
	TotalExpense    *Currency `json:"TotalExpense,omitempty"`
}

// ChargeComponent A charge on the seller's account.  Possible values:  * Principal - The selling price of the order item, equal to the selling price of the item multiplied by the quantity ordered.  * Tax - The tax collected by the seller on the Principal.  * MarketplaceFacilitatorTax-Principal - The tax withheld on the Principal.  * MarketplaceFacilitatorTax-Shipping - The tax withheld on the ShippingCharge.  * MarketplaceFacilitatorTax-Giftwrap - The tax withheld on the Giftwrap charge.  * MarketplaceFacilitatorTax-Other - The tax withheld on other miscellaneous charges.  * Discount - The promotional discount for an order item.  * TaxDiscount - The tax amount deducted for promotional rebates.  * CODItemCharge - The COD charge for an order item.  * CODItemTaxCharge - The tax collected by the seller on a CODItemCharge.  * CODOrderCharge - The COD charge for an order.  * CODOrderTaxCharge - The tax collected by the seller on a CODOrderCharge.  * CODShippingCharge - Shipping charges for a COD order.  * CODShippingTaxCharge - The tax collected by the seller on a CODShippingCharge.  * ShippingCharge - The shipping charge.  * ShippingTax - The tax collected by the seller on a ShippingCharge.  * Goodwill - The amount given to a buyer as a gesture of goodwill or to compensate for pain and suffering in the buying experience.  * Giftwrap - The gift wrap charge.  * GiftwrapTax - The tax collected by the seller on a Giftwrap charge.  * RestockingFee - The charge applied to the buyer when returning a product in certain categories.  * ReturnShipping - The amount given to the buyer to compensate for shipping the item back in the event we are at fault.  * PointsFee - The value of Amazon Points deducted from the refund if the buyer does not have enough Amazon Points to cover the deduction.  * GenericDeduction - A generic bad debt deduction.  * FreeReplacementReturnShipping - The compensation for return shipping when a buyer receives the wrong item, requests a free replacement, and returns the incorrect item.  * PaymentMethodFee - The fee collected for certain payment methods in certain marketplaces.  * ExportCharge - The export duty that is charged when an item is shipped to an international destination as part of the Amazon Global program.  * SAFE-TReimbursement - The SAFE-T claim amount for the item.  * TCS-CGST - Tax Collected at Source (TCS) for Central Goods and Services Tax (CGST).  * TCS-SGST - Tax Collected at Source for State Goods and Services Tax (SGST).  * TCS-IGST - Tax Collected at Source for Integrated Goods and Services Tax (IGST).  * TCS-UTGST - Tax Collected at Source for Union Territories Goods and Services Tax (UTGST).
type ChargeComponent struct {
	// The type of charge.
	ChargeType   *string   `json:"ChargeType,omitempty"`
	ChargeAmount *Currency `json:"ChargeAmount,omitempty"`
}

// ChargeInstrument A payment instrument.
type ChargeInstrument struct {
	// A short description of the charge instrument.
	Description *string `json:"Description,omitempty"`
	// The account tail (trailing digits) of the charge instrument.
	Tail   *string   `json:"Tail,omitempty"`
	Amount *Currency `json:"Amount,omitempty"`
}

// CouponPaymentEvent An event related to coupon payments.
type CouponPaymentEvent struct {
	PostedDate *time.Time `json:"PostedDate,omitempty"`
	// A coupon identifier.
	CouponId *string `json:"CouponId,omitempty"`
	// The description provided by the seller when they created the coupon.
	SellerCouponDescription *string `json:"SellerCouponDescription,omitempty"`
	// The number of coupon clips or redemptions.
	ClipOrRedemptionCount *int64 `json:"ClipOrRedemptionCount,omitempty"`
	// A payment event identifier.
	PaymentEventId  *string          `json:"PaymentEventId,omitempty"`
	FeeComponent    *FeeComponent    `json:"FeeComponent,omitempty"`
	ChargeComponent *ChargeComponent `json:"ChargeComponent,omitempty"`
	TotalAmount     *Currency        `json:"TotalAmount,omitempty"`
}

// DebtRecoveryEvent A debt payment or debt adjustment.
type DebtRecoveryEvent struct {
	// The debt recovery type.  Possible values:  * DebtPayment  * DebtPaymentFailure  *DebtAdjustment
	DebtRecoveryType  *string   `json:"DebtRecoveryType,omitempty"`
	RecoveryAmount    *Currency `json:"RecoveryAmount,omitempty"`
	OverPaymentCredit *Currency `json:"OverPaymentCredit,omitempty"`
	// A list of debt recovery item information.
	DebtRecoveryItemList []DebtRecoveryItem `json:"DebtRecoveryItemList,omitempty"`
	// A list of payment instruments.
	ChargeInstrumentList []ChargeInstrument `json:"ChargeInstrumentList,omitempty"`
}

// DebtRecoveryItem An item of a debt payment or debt adjustment.
type DebtRecoveryItem struct {
	RecoveryAmount *Currency  `json:"RecoveryAmount,omitempty"`
	OriginalAmount *Currency  `json:"OriginalAmount,omitempty"`
	GroupBeginDate *time.Time `json:"GroupBeginDate,omitempty"`
	GroupEndDate   *time.Time `json:"GroupEndDate,omitempty"`
}

// DirectPayment A payment made directly to a seller.
type DirectPayment struct {
	// The type of payment.  Possible values:  * StoredValueCardRevenue - The amount that is deducted from the seller's account because the seller received money through a stored value card.  * StoredValueCardRefund - The amount that Amazon returns to the seller if the order that is bought using a stored value card is refunded.  * PrivateLabelCreditCardRevenue - The amount that is deducted from the seller's account because the seller received money through a private label credit card offered by Amazon.  * PrivateLabelCreditCardRefund - The amount that Amazon returns to the seller if the order that is bought using a private label credit card offered by Amazon is refunded.  * CollectOnDeliveryRevenue - The COD amount that the seller collected directly from the buyer.  * CollectOnDeliveryRefund - The amount that Amazon refunds to the buyer if an order paid for by COD is refunded.
	DirectPaymentType   *string   `json:"DirectPaymentType,omitempty"`
	DirectPaymentAmount *Currency `json:"DirectPaymentAmount,omitempty"`
}

// FBALiquidationEvent A payment event for Fulfillment by Amazon (FBA) inventory liquidation. This event is used only in the US marketplace.
type FBALiquidationEvent struct {
	PostedDate *time.Time `json:"PostedDate,omitempty"`
	// The identifier for the original removal order.
	OriginalRemovalOrderId    *string   `json:"OriginalRemovalOrderId,omitempty"`
	LiquidationProceedsAmount *Currency `json:"LiquidationProceedsAmount,omitempty"`
	LiquidationFeeAmount      *Currency `json:"LiquidationFeeAmount,omitempty"`
}

// FeeComponent A fee associated with the event.
type FeeComponent struct {
	// The type of fee. For more information about Selling on Amazon fees, see [Selling on Amazon Fee Schedule](https://sellercentral.amazon.com/gp/help/200336920) on Seller Central. For more information about Fulfillment by Amazon fees, see [FBA features, services and fees](https://sellercentral.amazon.com/gp/help/201074400) on Seller Central.
	FeeType   *string   `json:"FeeType,omitempty"`
	FeeAmount *Currency `json:"FeeAmount,omitempty"`
}

// FinancialEventGroup Information related to a financial event group.
type FinancialEventGroup struct {
	// A unique identifier for the financial event group.
	FinancialEventGroupId *string `json:"FinancialEventGroupId,omitempty"`
	// The processing status of the financial event group indicates whether the balance of the financial event group is settled.  Possible values:  * Open  * Closed
	ProcessingStatus *string `json:"ProcessingStatus,omitempty"`
	// The status of the fund transfer.
	FundTransferStatus *string    `json:"FundTransferStatus,omitempty"`
	OriginalTotal      *Currency  `json:"OriginalTotal,omitempty"`
	ConvertedTotal     *Currency  `json:"ConvertedTotal,omitempty"`
	FundTransferDate   *time.Time `json:"FundTransferDate,omitempty"`
	// The trace identifier used by sellers to look up transactions externally.
	TraceId *string `json:"TraceId,omitempty"`
	// The account tail of the payment instrument.
	AccountTail              *string    `json:"AccountTail,omitempty"`
	BeginningBalance         *Currency  `json:"BeginningBalance,omitempty"`
	FinancialEventGroupStart *time.Time `json:"FinancialEventGroupStart,omitempty"`
	FinancialEventGroupEnd   *time.Time `json:"FinancialEventGroupEnd,omitempty"`
}

// FinancialEvents Contains all information related to a financial event.
type FinancialEvents struct {
	// A list of shipment event information.
	ShipmentEventList []ShipmentEvent `json:"ShipmentEventList,omitempty"`
	// A list of shipment event information.
	RefundEventList []ShipmentEvent `json:"RefundEventList,omitempty"`
	// A list of shipment event information.
	GuaranteeClaimEventList []ShipmentEvent `json:"GuaranteeClaimEventList,omitempty"`
	// A list of shipment event information.
	ChargebackEventList []ShipmentEvent `json:"ChargebackEventList,omitempty"`
	// A list of events related to the seller's Pay with Amazon account.
	PayWithAmazonEventList []PayWithAmazonEvent `json:"PayWithAmazonEventList,omitempty"`
	// A list of information about solution provider credits.
	ServiceProviderCreditEventList []SolutionProviderCreditEvent `json:"ServiceProviderCreditEventList,omitempty"`
	// A list of information about Retrocharge or RetrochargeReversal events.
	RetrochargeEventList []RetrochargeEvent `json:"RetrochargeEventList,omitempty"`
	// A list of rental transaction event information.
	RentalTransactionEventList []RentalTransactionEvent `json:"RentalTransactionEventList,omitempty"`
	// A list of sponsored products payment events.
	ProductAdsPaymentEventList []ProductAdsPaymentEvent `json:"ProductAdsPaymentEventList,omitempty"`
	// A list of information about service fee events.
	ServiceFeeEventList []ServiceFeeEvent `json:"ServiceFeeEventList,omitempty"`
	// A list of payment events for deal-related fees.
	SellerDealPaymentEventList []SellerDealPaymentEvent `json:"SellerDealPaymentEventList,omitempty"`
	// A list of debt recovery event information.
	DebtRecoveryEventList []DebtRecoveryEvent `json:"DebtRecoveryEventList,omitempty"`
	// A list of loan servicing events.
	LoanServicingEventList []LoanServicingEvent `json:"LoanServicingEventList,omitempty"`
	// A list of adjustment event information for the seller's account.
	AdjustmentEventList []AdjustmentEvent `json:"AdjustmentEventList,omitempty"`
	// A list of SAFETReimbursementEvents.
	SAFETReimbursementEventList []SAFETReimbursementEvent `json:"SAFETReimbursementEventList,omitempty"`
	// A list of information about fee events for the Early Reviewer Program.
	SellerReviewEnrollmentPaymentEventList []SellerReviewEnrollmentPaymentEvent `json:"SellerReviewEnrollmentPaymentEventList,omitempty"`
	// A list of FBA inventory liquidation payment events.
	FBALiquidationEventList []FBALiquidationEvent `json:"FBALiquidationEventList,omitempty"`
	// A list of coupon payment event information.
	CouponPaymentEventList []CouponPaymentEvent `json:"CouponPaymentEventList,omitempty"`
	// A list of fee events related to Amazon Imaging services.
	ImagingServicesFeeEventList []ImagingServicesFeeEvent `json:"ImagingServicesFeeEventList,omitempty"`
	// A list of network commingling transaction events.
	NetworkComminglingTransactionEventList []NetworkComminglingTransactionEvent `json:"NetworkComminglingTransactionEventList,omitempty"`
	// A list of expense information related to an affordability promotion.
	AffordabilityExpenseEventList []AffordabilityExpenseEvent `json:"AffordabilityExpenseEventList,omitempty"`
	// A list of expense information related to an affordability promotion.
	AffordabilityExpenseReversalEventList []AffordabilityExpenseEvent `json:"AffordabilityExpenseReversalEventList,omitempty"`
	// A list of information about trial shipment financial events.
	TrialShipmentEventList []TrialShipmentEvent `json:"TrialShipmentEventList,omitempty"`
	// A list of information about shipment settle financial events.
	ShipmentSettleEventList []ShipmentEvent `json:"ShipmentSettleEventList,omitempty"`
	// List of TaxWithholding events.
	TaxWithholdingEventList []TaxWithholdingEvent `json:"TaxWithholdingEventList,omitempty"`
	// A list of removal shipment event information.
	RemovalShipmentEventList []RemovalShipmentEvent `json:"RemovalShipmentEventList,omitempty"`
	// A comma-delimited list of Removal shipmentAdjustment details for FBA inventory.
	RemovalShipmentAdjustmentEventList []RemovalShipmentAdjustmentEvent `json:"RemovalShipmentAdjustmentEventList,omitempty"`
}

// ImagingServicesFeeEvent A fee event related to Amazon Imaging services.
type ImagingServicesFeeEvent struct {
	// The identifier for the imaging services request.
	ImagingRequestBillingItemID *string `json:"ImagingRequestBillingItemID,omitempty"`
	// The Amazon Standard Identification Number (ASIN) of the item for which the imaging service was requested.
	ASIN       *string    `json:"ASIN,omitempty"`
	PostedDate *time.Time `json:"PostedDate,omitempty"`
	// A list of fee component information.
	FeeList []FeeComponent `json:"FeeList,omitempty"`
}

// ListFinancialEventGroupsPayload The payload for the listFinancialEventGroups operation.
type ListFinancialEventGroupsPayload struct {
	// When present and not empty, pass this string token in the next request to return the next response page.
	NextToken *string `json:"NextToken,omitempty"`
	// A list of financial event group information.
	FinancialEventGroupList []FinancialEventGroup `json:"FinancialEventGroupList,omitempty"`
}

// ListFinancialEventGroupsResponse The response schema for the listFinancialEventGroups operation.
type ListFinancialEventGroupsResponse struct {
	Payload *ListFinancialEventGroupsPayload `json:"payload,omitempty"`
	// A list of error responses returned when a request is unsuccessful.
	Errors []apis.Error `json:"errors,omitempty"`
}

// ListFinancialEventsPayload The payload for the listFinancialEvents operation.
type ListFinancialEventsPayload struct {
	// When present and not empty, pass this string token in the next request to return the next response page.
	NextToken       *string          `json:"NextToken,omitempty"`
	FinancialEvents *FinancialEvents `json:"FinancialEvents,omitempty"`
}

// ListFinancialEventsResponse The response schema for the listFinancialEvents operation.
type ListFinancialEventsResponse struct {
	Payload *ListFinancialEventsPayload `json:"payload,omitempty"`
	// A list of error responses returned when a request is unsuccessful.
	Errors []apis.Error `json:"errors,omitempty"`
}

// LoanServicingEvent A loan advance, loan payment, or loan refund.
type LoanServicingEvent struct {
	LoanAmount *Currency `json:"LoanAmount,omitempty"`
	// The type of event.  Possible values:  * LoanAdvance  * LoanPayment  * LoanRefund
	SourceBusinessEventType *string `json:"SourceBusinessEventType,omitempty"`
}

// NetworkComminglingTransactionEvent A network commingling transaction event.
type NetworkComminglingTransactionEvent struct {
	// The type of network item swap.  Possible values:  * NetCo - A Fulfillment by Amazon inventory pooling transaction. Available only in the India marketplace.  * ComminglingVAT - A commingling VAT transaction. Available only in the UK, Spain, France, Germany, and Italy marketplaces.
	TransactionType *string    `json:"TransactionType,omitempty"`
	PostedDate      *time.Time `json:"PostedDate,omitempty"`
	// The identifier for the network item swap.
	NetCoTransactionID *string `json:"NetCoTransactionID,omitempty"`
	// The reason for the network item swap.
	SwapReason *string `json:"SwapReason,omitempty"`
	// The Amazon Standard Identification Number (ASIN) of the swapped item.
	ASIN *string `json:"ASIN,omitempty"`
	// The marketplace in which the event took place.
	MarketplaceId      *string   `json:"MarketplaceId,omitempty"`
	TaxExclusiveAmount *Currency `json:"TaxExclusiveAmount,omitempty"`
	TaxAmount          *Currency `json:"TaxAmount,omitempty"`
}

// PayWithAmazonEvent An event related to the seller's Pay with Amazon account.
type PayWithAmazonEvent struct {
	// An order identifier that is specified by the seller.
	SellerOrderId         *string    `json:"SellerOrderId,omitempty"`
	TransactionPostedDate *time.Time `json:"TransactionPostedDate,omitempty"`
	// The type of business object.
	BusinessObjectType *string `json:"BusinessObjectType,omitempty"`
	// The sales channel for the transaction.
	SalesChannel *string          `json:"SalesChannel,omitempty"`
	Charge       *ChargeComponent `json:"Charge,omitempty"`
	// A list of fee component information.
	FeeList []FeeComponent `json:"FeeList,omitempty"`
	// The type of payment.  Possible values:  * Sales
	PaymentAmountType *string `json:"PaymentAmountType,omitempty"`
	// A short description of this payment event.
	AmountDescription *string `json:"AmountDescription,omitempty"`
	// The fulfillment channel.  Possible values:  * AFN - Amazon Fulfillment Network (Fulfillment by Amazon)  * MFN - Merchant Fulfillment Network (self-fulfilled)
	FulfillmentChannel *string `json:"FulfillmentChannel,omitempty"`
	// The store name where the event occurred.
	StoreName *string `json:"StoreName,omitempty"`
}

// ProductAdsPaymentEvent A Sponsored Products payment event.
type ProductAdsPaymentEvent struct {
	PostedDate *time.Time `json:"postedDate,omitempty"`
	// Indicates if the transaction is for a charge or a refund.  Possible values:  * charge - Charge  * refund - Refund
	TransactionType *string `json:"transactionType,omitempty"`
	// Identifier for the invoice that the transaction appears in.
	InvoiceId        *string   `json:"invoiceId,omitempty"`
	BaseValue        *Currency `json:"baseValue,omitempty"`
	TaxValue         *Currency `json:"taxValue,omitempty"`
	TransactionValue *Currency `json:"transactionValue,omitempty"`
}

// Promotion A promotion applied to an item.
type Promotion struct {
	// The type of promotion.
	PromotionType *string `json:"PromotionType,omitempty"`
	// The seller-specified identifier for the promotion.
	PromotionId     *string   `json:"PromotionId,omitempty"`
	PromotionAmount *Currency `json:"PromotionAmount,omitempty"`
}

// RemovalShipmentAdjustmentEvent A financial adjustment event for FBA liquidated inventory. A positive value indicates money owed to Amazon by the buyer (for example, when the charge was incorrectly calculated as less than it should be). A negative value indicates a full or partial refund owed to the buyer (for example, when the buyer receives damaged items or fewer items than ordered).
type RemovalShipmentAdjustmentEvent struct {
	PostedDate *time.Time `json:"PostedDate,omitempty"`
	// The unique identifier for the adjustment event.
	AdjustmentEventId *string `json:"AdjustmentEventId,omitempty"`
	// The merchant removal orderId.
	MerchantOrderId *string `json:"MerchantOrderId,omitempty"`
	// The orderId for shipping inventory.
	OrderId *string `json:"OrderId,omitempty"`
	// The type of removal order.  Possible values:  * WHOLESALE_LIQUIDATION.
	TransactionType *string `json:"TransactionType,omitempty"`
	// A comma-delimited list of Removal shipmentItemAdjustment details for FBA inventory.
	RemovalShipmentItemAdjustmentList []RemovalShipmentItemAdjustment `json:"RemovalShipmentItemAdjustmentList,omitempty"`
}

// RemovalShipmentEvent A removal shipment event for a removal order.
type RemovalShipmentEvent struct {
	PostedDate *time.Time `json:"PostedDate,omitempty"`
	// The merchant removal orderId.
	MerchantOrderId *string `json:"MerchantOrderId,omitempty"`
	// The identifier for the removal shipment order.
	OrderId *string `json:"OrderId,omitempty"`
	// The type of removal order.  Possible values:  * WHOLESALE_LIQUIDATION
	TransactionType *string `json:"TransactionType,omitempty"`
	// A list of information about removal shipment items.
	RemovalShipmentItemList []RemovalShipmentItem `json:"RemovalShipmentItemList,omitempty"`
}

// RemovalShipmentItem Item-level information for a removal shipment.
type RemovalShipmentItem struct {
	// An identifier for an item in a removal shipment.
	RemovalShipmentItemId *string `json:"RemovalShipmentItemId,omitempty"`
	// The tax collection model applied to the item.  Possible values:  * MarketplaceFacilitator - Tax is withheld and remitted to the taxing authority by Amazon on behalf of the seller.  * Standard - Tax is paid to the seller and not remitted to the taxing authority by Amazon.
	TaxCollectionModel *string `json:"TaxCollectionModel,omitempty"`
	// The Amazon fulfillment network SKU for the item.
	FulfillmentNetworkSKU *string `json:"FulfillmentNetworkSKU,omitempty"`
	// The quantity of the item.
	Quantity    *int32    `json:"Quantity,omitempty"`
	Revenue     *Currency `json:"Revenue,omitempty"`
	FeeAmount   *Currency `json:"FeeAmount,omitempty"`
	TaxAmount   *Currency `json:"TaxAmount,omitempty"`
	TaxWithheld *Currency `json:"TaxWithheld,omitempty"`
}

// RemovalShipmentItemAdjustment Item-level information for a removal shipment item adjustment.
type RemovalShipmentItemAdjustment struct {
	// An identifier for an item in a removal shipment.
	RemovalShipmentItemId *string `json:"RemovalShipmentItemId,omitempty"`
	// The tax collection model applied to the item.  Possible values:  * MarketplaceFacilitator - Tax is withheld and remitted to the taxing authority by Amazon on behalf of the seller.  * Standard - Tax is paid to the seller and not remitted to the taxing authority by Amazon.
	TaxCollectionModel *string `json:"TaxCollectionModel,omitempty"`
	// The Amazon fulfillment network SKU for the item.
	FulfillmentNetworkSKU *string `json:"FulfillmentNetworkSKU,omitempty"`
	// Adjusted quantity of removal shipmentItemAdjustment items.
	AdjustedQuantity      *int32    `json:"AdjustedQuantity,omitempty"`
	RevenueAdjustment     *Currency `json:"RevenueAdjustment,omitempty"`
	TaxAmountAdjustment   *Currency `json:"TaxAmountAdjustment,omitempty"`
	TaxWithheldAdjustment *Currency `json:"TaxWithheldAdjustment,omitempty"`
}

// RentalTransactionEvent An event related to a rental transaction.
type RentalTransactionEvent struct {
	// An Amazon-defined identifier for an order.
	AmazonOrderId *string `json:"AmazonOrderId,omitempty"`
	// The type of rental event.  Possible values:  * RentalCustomerPayment-Buyout - Transaction type that represents when the customer wants to buy out a rented item.  * RentalCustomerPayment-Extension - Transaction type that represents when the customer wants to extend the rental period.  * RentalCustomerRefund-Buyout - Transaction type that represents when the customer requests a refund for the buyout of the rented item.  * RentalCustomerRefund-Extension - Transaction type that represents when the customer requests a refund over the extension on the rented item.  * RentalHandlingFee - Transaction type that represents the fee that Amazon charges sellers who rent through Amazon.  * RentalChargeFailureReimbursement - Transaction type that represents when Amazon sends money to the seller to compensate for a failed charge.  * RentalLostItemReimbursement - Transaction type that represents when Amazon sends money to the seller to compensate for a lost item.
	RentalEventType *string `json:"RentalEventType,omitempty"`
	// The number of days that the buyer extended an already rented item. This value is only returned for RentalCustomerPayment-Extension and RentalCustomerRefund-Extension events.
	ExtensionLength *int32     `json:"ExtensionLength,omitempty"`
	PostedDate      *time.Time `json:"PostedDate,omitempty"`
	// A list of charge information on the seller's account.
	RentalChargeList []ChargeComponent `json:"RentalChargeList,omitempty"`
	// A list of fee component information.
	RentalFeeList []FeeComponent `json:"RentalFeeList,omitempty"`
	// The name of the marketplace.
	MarketplaceName     *string   `json:"MarketplaceName,omitempty"`
	RentalInitialValue  *Currency `json:"RentalInitialValue,omitempty"`
	RentalReimbursement *Currency `json:"RentalReimbursement,omitempty"`
	// A list of information about taxes withheld.
	RentalTaxWithheldList []TaxWithheldComponent `json:"RentalTaxWithheldList,omitempty"`
}

// RetrochargeEvent A retrocharge or retrocharge reversal.
type RetrochargeEvent struct {
	// The type of event.  Possible values:  * Retrocharge  * RetrochargeReversal
	RetrochargeEventType *string `json:"RetrochargeEventType,omitempty"`
	// An Amazon-defined identifier for an order.
	AmazonOrderId *string    `json:"AmazonOrderId,omitempty"`
	PostedDate    *time.Time `json:"PostedDate,omitempty"`
	BaseTax       *Currency  `json:"BaseTax,omitempty"`
	ShippingTax   *Currency  `json:"ShippingTax,omitempty"`
	// The name of the marketplace where the retrocharge event occurred.
	MarketplaceName *string `json:"MarketplaceName,omitempty"`
	// A list of information about taxes withheld.
	RetrochargeTaxWithheldList []TaxWithheldComponent `json:"RetrochargeTaxWithheldList,omitempty"`
}

// SAFETReimbursementEvent A SAFE-T claim reimbursement on the seller's account.
type SAFETReimbursementEvent struct {
	PostedDate *time.Time `json:"PostedDate,omitempty"`
	// A SAFE-T claim identifier.
	SAFETClaimId     *string   `json:"SAFETClaimId,omitempty"`
	ReimbursedAmount *Currency `json:"ReimbursedAmount,omitempty"`
	// Indicates why the seller was reimbursed.
	ReasonCode *string `json:"ReasonCode,omitempty"`
	// A list of SAFETReimbursementItems.
	SAFETReimbursementItemList []SAFETReimbursementItem `json:"SAFETReimbursementItemList,omitempty"`
}

// SAFETReimbursementItem An item from a SAFE-T claim reimbursement.
type SAFETReimbursementItem struct {
	// A list of charge information on the seller's account.
	ItemChargeList []ChargeComponent `json:"itemChargeList,omitempty"`
	// The description of the item as shown on the product detail page on the retail website.
	ProductDescription *string `json:"productDescription,omitempty"`
	// The number of units of the item being reimbursed.
	Quantity *string `json:"quantity,omitempty"`
}

// SellerDealPaymentEvent An event linked to the payment of a fee related to the specified deal.
type SellerDealPaymentEvent struct {
	PostedDate *time.Time `json:"postedDate,omitempty"`
	// The unique identifier of the deal.
	DealId *string `json:"dealId,omitempty"`
	// The internal description of the deal.
	DealDescription *string `json:"dealDescription,omitempty"`
	// The type of event: SellerDealComplete.
	EventType *string `json:"eventType,omitempty"`
	// The type of fee: RunLightningDealFee.
	FeeType     *string   `json:"feeType,omitempty"`
	FeeAmount   *Currency `json:"feeAmount,omitempty"`
	TaxAmount   *Currency `json:"taxAmount,omitempty"`
	TotalAmount *Currency `json:"totalAmount,omitempty"`
}

// SellerReviewEnrollmentPaymentEvent A fee payment event for the Early Reviewer Program.
type SellerReviewEnrollmentPaymentEvent struct {
	PostedDate *time.Time `json:"PostedDate,omitempty"`
	// An enrollment identifier.
	EnrollmentId *string `json:"EnrollmentId,omitempty"`
	// The Amazon Standard Identification Number (ASIN) of the item that was enrolled in the Early Reviewer Program.
	ParentASIN      *string          `json:"ParentASIN,omitempty"`
	FeeComponent    *FeeComponent    `json:"FeeComponent,omitempty"`
	ChargeComponent *ChargeComponent `json:"ChargeComponent,omitempty"`
	TotalAmount     *Currency        `json:"TotalAmount,omitempty"`
}

// ServiceFeeEvent A service fee on the seller's account.
type ServiceFeeEvent struct {
	// An Amazon-defined identifier for an order.
	AmazonOrderId *string `json:"AmazonOrderId,omitempty"`
	// A short description of the service fee reason.
	FeeReason *string `json:"FeeReason,omitempty"`
	// A list of fee component information.
	FeeList []FeeComponent `json:"FeeList,omitempty"`
	// The seller SKU of the item. The seller SKU is qualified by the seller's seller ID, which is included with every call to the Selling Partner API.
	SellerSKU *string `json:"SellerSKU,omitempty"`
	// A unique identifier assigned by Amazon to products stored in and fulfilled from an Amazon fulfillment center.
	FnSKU *string `json:"FnSKU,omitempty"`
	// A short description of the service fee event.
	FeeDescription *string `json:"FeeDescription,omitempty"`
	// The Amazon Standard Identification Number (ASIN) of the item.
	ASIN *string `json:"ASIN,omitempty"`
}

// ShipmentEvent A shipment, refund, guarantee claim, or chargeback.
type ShipmentEvent struct {
	// An Amazon-defined identifier for an order.
	AmazonOrderId *string `json:"AmazonOrderId,omitempty"`
	// A seller-defined identifier for an order.
	SellerOrderId *string `json:"SellerOrderId,omitempty"`
	// The name of the marketplace where the event occurred.
	MarketplaceName *string `json:"MarketplaceName,omitempty"`
	// A list of charge information on the seller's account.
	OrderChargeList []ChargeComponent `json:"OrderChargeList,omitempty"`
	// A list of charge information on the seller's account.
	OrderChargeAdjustmentList []ChargeComponent `json:"OrderChargeAdjustmentList,omitempty"`
	// A list of fee component information.
	ShipmentFeeList []FeeComponent `json:"ShipmentFeeList,omitempty"`
	// A list of fee component information.
	ShipmentFeeAdjustmentList []FeeComponent `json:"ShipmentFeeAdjustmentList,omitempty"`
	// A list of fee component information.
	OrderFeeList []FeeComponent `json:"OrderFeeList,omitempty"`
	// A list of fee component information.
	OrderFeeAdjustmentList []FeeComponent `json:"OrderFeeAdjustmentList,omitempty"`
	// A list of direct payment information.
	DirectPaymentList []DirectPayment `json:"DirectPaymentList,omitempty"`
	PostedDate        *time.Time      `json:"PostedDate,omitempty"`
	// A list of shipment items.
	ShipmentItemList []ShipmentItem `json:"ShipmentItemList,omitempty"`
	// A list of shipment items.
	ShipmentItemAdjustmentList []ShipmentItem `json:"ShipmentItemAdjustmentList,omitempty"`
}

// ShipmentItem An item of a shipment, refund, guarantee claim, or chargeback.
type ShipmentItem struct {
	// The seller SKU of the item. The seller SKU is qualified by the seller's seller ID, which is included with every call to the Selling Partner API.
	SellerSKU *string `json:"SellerSKU,omitempty"`
	// An Amazon-defined order item identifier.
	OrderItemId *string `json:"OrderItemId,omitempty"`
	// An Amazon-defined order adjustment identifier defined for refunds, guarantee claims, and chargeback events.
	OrderAdjustmentItemId *string `json:"OrderAdjustmentItemId,omitempty"`
	// The number of items shipped.
	QuantityShipped *int32 `json:"QuantityShipped,omitempty"`
	// A list of charge information on the seller's account.
	ItemChargeList []ChargeComponent `json:"ItemChargeList,omitempty"`
	// A list of charge information on the seller's account.
	ItemChargeAdjustmentList []ChargeComponent `json:"ItemChargeAdjustmentList,omitempty"`
	// A list of fee component information.
	ItemFeeList []FeeComponent `json:"ItemFeeList,omitempty"`
	// A list of fee component information.
	ItemFeeAdjustmentList []FeeComponent `json:"ItemFeeAdjustmentList,omitempty"`
	// A list of information about taxes withheld.
	ItemTaxWithheldList []TaxWithheldComponent `json:"ItemTaxWithheldList,omitempty"`
	// A list of promotions.
	PromotionList []Promotion `json:"PromotionList,omitempty"`
	// A list of promotions.
	PromotionAdjustmentList []Promotion `json:"PromotionAdjustmentList,omitempty"`
	CostOfPointsGranted     *Currency   `json:"CostOfPointsGranted,omitempty"`
	CostOfPointsReturned    *Currency   `json:"CostOfPointsReturned,omitempty"`
}

// SolutionProviderCreditEvent A credit given to a solution provider.
type SolutionProviderCreditEvent struct {
	// The transaction type.
	ProviderTransactionType *string `json:"ProviderTransactionType,omitempty"`
	// A seller-defined identifier for an order.
	SellerOrderId *string `json:"SellerOrderId,omitempty"`
	// The identifier of the marketplace where the order was placed.
	MarketplaceId *string `json:"MarketplaceId,omitempty"`
	// The two-letter country code of the country associated with the marketplace where the order was placed.
	MarketplaceCountryCode *string `json:"MarketplaceCountryCode,omitempty"`
	// The Amazon-defined identifier of the seller.
	SellerId *string `json:"SellerId,omitempty"`
	// The store name where the payment event occurred.
	SellerStoreName *string `json:"SellerStoreName,omitempty"`
	// The Amazon-defined identifier of the solution provider.
	ProviderId *string `json:"ProviderId,omitempty"`
	// The store name where the payment event occurred.
	ProviderStoreName       *string    `json:"ProviderStoreName,omitempty"`
	TransactionAmount       *Currency  `json:"TransactionAmount,omitempty"`
	TransactionCreationDate *time.Time `json:"TransactionCreationDate,omitempty"`
}

// TaxWithheldComponent Information about the taxes withheld.
type TaxWithheldComponent struct {
	// The tax collection model applied to the item.  Possible values:  * MarketplaceFacilitator - Tax is withheld and remitted to the taxing authority by Amazon on behalf of the seller.  * Standard - Tax is paid to the seller and not remitted to the taxing authority by Amazon.
	TaxCollectionModel *string `json:"TaxCollectionModel,omitempty"`
	// A list of charge information on the seller's account.
	TaxesWithheld []ChargeComponent `json:"TaxesWithheld,omitempty"`
}

// TaxWithholdingEvent A TaxWithholding event on seller's account.
type TaxWithholdingEvent struct {
	PostedDate           *time.Time            `json:"PostedDate,omitempty"`
	BaseAmount           *Currency             `json:"BaseAmount,omitempty"`
	WithheldAmount       *Currency             `json:"WithheldAmount,omitempty"`
	TaxWithholdingPeriod *TaxWithholdingPeriod `json:"TaxWithholdingPeriod,omitempty"`
}

// TaxWithholdingPeriod Period which taxwithholding on seller's account is calculated.
type TaxWithholdingPeriod struct {
	StartDate *time.Time `json:"StartDate,omitempty"`
	EndDate   *time.Time `json:"EndDate,omitempty"`
}

// TrialShipmentEvent An event related to a trial shipment.
type TrialShipmentEvent struct {
	// An Amazon-defined identifier for an order.
	AmazonOrderId *string `json:"AmazonOrderId,omitempty"`
	// The identifier of the financial event group.
	FinancialEventGroupId *string    `json:"FinancialEventGroupId,omitempty"`
	PostedDate            *time.Time `json:"PostedDate,omitempty"`
	// The seller SKU of the item. The seller SKU is qualified by the seller's seller ID, which is included with every call to the Selling Partner API.
	SKU *string `json:"SKU,omitempty"`
	// A list of fee component information.
	FeeList []FeeComponent `json:"FeeList,omitempty"`
}
