package api

// SearchResponse represents the search API response
type SearchResponse struct {
	Type       string           `json:"type"`
	Attributes SearchAttributes `json:"attributes"`
	Items      []Product        `json:"items"`
}

// SearchAttributes contains metadata about search results
type SearchAttributes struct {
	Items        int  `json:"items"`
	Page         int  `json:"page"`
	HasMoreItems bool `json:"has_more_items"`
}

// Product represents a product in search results
type Product struct {
	ID         int               `json:"id"`
	Type       string            `json:"type"`
	Attributes ProductAttributes `json:"attributes"`
}

// ProductAttributes contains product details
type ProductAttributes struct {
	Name                      string         `json:"name"`
	FullName                  string         `json:"full_name"`
	Brand                     string         `json:"brand"`
	NameExtra                 string         `json:"name_extra"`
	GrossPrice                string         `json:"gross_price"`
	GrossUnitPrice            string         `json:"gross_unit_price"`
	UnitPriceQuantityAbbr     string         `json:"unit_price_quantity_abbreviation"`
	Currency                  string         `json:"currency"`
	Availability              Availability   `json:"availability"`
	Images                    []ProductImage `json:"images"`
}

// Availability indicates if a product is available
type Availability struct {
	IsAvailable bool   `json:"is_available"`
	Description string `json:"description"`
	Code        string `json:"code"`
}

// ProductImage represents product images
type ProductImage struct {
	Large     ImageVariant `json:"large"`
	Thumbnail ImageVariant `json:"thumbnail"`
	Variant   string       `json:"variant"`
}

// ImageVariant contains image URL and dimensions
type ImageVariant struct {
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Cart represents the shopping cart
type Cart struct {
	ID                   int            `json:"id"`
	ActiveGrouping       string         `json:"active_grouping"`
	LabelText            string         `json:"label_text"`
	ProductQuantityCount int            `json:"product_quantity_count"`
	DisplayPrice         string         `json:"display_price"`
	TotalGrossAmount     string         `json:"total_gross_amount"`
	Currency             string         `json:"currency"`
	Groups               []CartGroup    `json:"groups"`
	SummaryLines         []SummaryGroup `json:"summary_lines"`
}

// CartGroup represents a group of items in the cart
type CartGroup struct {
	Items []CartGroupItem `json:"items"`
}

// CartGroupItem represents an item in the cart
type CartGroupItem struct {
	Product      CartProduct `json:"product"`
	ItemID       int         `json:"item_id"`
	Quantity     int         `json:"quantity"`
	DisplayPrice string      `json:"display_price_total"`
}

// CartProduct contains product info within the cart
type CartProduct struct {
	ID            int    `json:"id"`
	FullName      string `json:"full_name"`
	Brand         string `json:"brand"`
	Name          string `json:"name"`
	NameExtra     string `json:"name_extra"`
	GrossPrice    string `json:"gross_price"`
	Currency      string `json:"currency"`
	AbsoluteURL   string `json:"absolute_url"`
}

// SummaryGroup represents a summary section
type SummaryGroup struct {
	ID    string        `json:"id"`
	Lines []SummaryLine `json:"lines"`
}

// SummaryLine represents a line in the cart summary
type SummaryLine struct {
	Description string `json:"description"`
	GrossAmount string `json:"gross_amount"`
	Name        string `json:"name"`
}

// CartItem represents an item to add to cart
type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
