package ozon

import "time"

// ItemResponse represents response from ozon api
type ItemResponse struct {
	Items []Item `json:"items"`
}

// Item represents item from ozon api
type Item struct {
	Barcodes              []string          `json:"barcodes"`
	ColorImage            []string          `json:"color_image"`
	Commissions           []Commission      `json:"commissions"`
	CreatedAt             time.Time         `json:"created_at"`
	CurrencyCode          string            `json:"currency_code"`
	DescriptionCategoryID int               `json:"description_category_id"`
	DiscountedFBOStocks   int               `json:"discounted_fbo_stocks"`
	Errors                []Error           `json:"errors"`
	HasDiscountedFBOItem  bool              `json:"has_discounted_fbo_item"`
	ID                    int               `json:"id"`
	Images                []string          `json:"images"`
	Images360             []string          `json:"images360"`
	IsArchived            bool              `json:"is_archived"`
	IsAutoarchived        bool              `json:"is_autoarchived"`
	IsDiscounted          bool              `json:"is_discounted"`
	IsKGT                 bool              `json:"is_kgt"`
	IsPrepaymentAllowed   bool              `json:"is_prepayment_allowed"`
	IsSuper               bool              `json:"is_super"`
	MarketingPrice        string            `json:"marketing_price"`
	MinPrice              string            `json:"min_price"`
	ModelInfo             ModelInfo         `json:"model_info"`
	Name                  string            `json:"name"`
	OfferID               string            `json:"offer_id"`
	OldPrice              string            `json:"old_price"`
	Price                 string            `json:"price"`
	PriceIndexes          PriceIndexes      `json:"price_indexes"`
	PrimaryImage          []string          `json:"primary_image"`
	Sources               []Source          `json:"sources"`
	Stocks                Stocks            `json:"stocks"`
	Statuses              Statuses          `json:"statuses"`
	TypeID                int               `json:"type_id"`
	UpdatedAt             time.Time         `json:"updated_at"`
	VAT                   string            `json:"vat"`
	VisibilityDetails     VisibilityDetails `json:"visibility_details"`
	VolumeWeight          float64           `json:"volume_weight"`
}

// Commission represents commission from ozon api
type Commission struct {
	DeliveryAmount int    `json:"delivery_amount"`
	Percent        int    `json:"percent"`
	ReturnAmount   int    `json:"return_amount"`
	SaleSchema     string `json:"sale_schema"`
	Value          int    `json:"value"`
}

// Error represents error from ozon api
type Error struct {
	AttributeID int    `json:"attribute_id"`
	Code        string `json:"code"`
	Field       string `json:"field"`
	Level       string `json:"level"`
	State       string `json:"state"`
	Texts       Texts  `json:"texts"`
}

// Texts represents texts from ozon api
type Texts struct {
	AttributeName string   `json:"attribute_name"`
	Description   string   `json:"description"`
	HintCode      string   `json:"hint_code"`
	Message       string   `json:"message"`
	Params        []Params `json:"params"`
	ShortDesc     string   `json:"short_description"`
}

// Params represents params from ozon api
type Params struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ModelInfo represents model info from ozon api
type ModelInfo struct {
	Count   int `json:"count"`
	ModelID int `json:"model_id"`
}

// PriceIndexes represents price indexes from ozon api
type PriceIndexes struct {
	ColorIndex                string                    `json:"color_index"`
	ExternalIndexData         ExternalIndexData         `json:"external_index_data"`
	OzonIndexData             OzonIndexData             `json:"ozon_index_data"`
	SelfMarketplacesIndexData SelfMarketplacesIndexData `json:"self_marketplaces_index_data"`
}

// ExternalIndexData represents external index data from ozon api
type ExternalIndexData struct {
	MinimalPrice         string  `json:"minimal_price"`
	MinimalPriceCurrency string  `json:"minimal_price_currency"`
	PriceIndexValue      float64 `json:"price_index_value"`
}

// OzonIndexData represents ozon index data from ozon api
type OzonIndexData struct {
	MinimalPrice         string  `json:"minimal_price"`
	MinimalPriceCurrency string  `json:"minimal_price_currency"`
	PriceIndexValue      float64 `json:"price_index_value"`
}

// SelfMarketplacesIndexData represents self marketplaces index data from ozon api
type SelfMarketplacesIndexData struct {
	MinimalPrice         string  `json:"minimal_price"`
	MinimalPriceCurrency string  `json:"minimal_price_currency"`
	PriceIndexValue      float64 `json:"price_index_value"`
}

// Source represents source from ozon api
type Source struct {
	CreatedAt    time.Time `json:"created_at"`
	QuantCode    string    `json:"quant_code"`
	ShipmentType string    `json:"shipment_type"`
	SKU          int       `json:"sku"`
	Source       string    `json:"source"`
}

// Stocks represents stocks from ozon api
type Stocks struct {
	HasStock bool    `json:"has_stock"`
	Stocks   []Stock `json:"stocks"`
}

// Stock represents stock from ozon api
type Stock struct {
	Present  int    `json:"present"`
	Reserved int    `json:"reserved"`
	SKU      int    `json:"sku"`
	Source   string `json:"source"`
}

// Statuses represents statuses from ozon api
type Statuses struct {
	IsCreated         bool      `json:"is_created"`
	ModerateStatus    string    `json:"moderate_status"`
	Status            string    `json:"status"`
	StatusDescription string    `json:"status_description"`
	StatusFailed      string    `json:"status_failed"`
	StatusName        string    `json:"status_name"`
	StatusTooltip     string    `json:"status_tooltip"`
	StatusUpdatedAt   time.Time `json:"status_updated_at"`
	ValidationStatus  string    `json:"validation_status"`
}

// VisibilityDetails represents visibility details from ozon api
type VisibilityDetails struct {
	HasPrice bool `json:"has_price"`
	HasStock bool `json:"has_stock"`
}
