package bitpin_client

import (
	"net/http"
	"strconv"
	"time"
)

type ClientOptions struct {
	HttpClient *http.Client

	AccessToken  string
	RefreshToken string

	ApiKey    string
	SecretKey string

	AutoAuth    bool
	AutoRefresh bool
}

type AuthResponse struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

type RefreshTokenResponse struct {
	Access string `json:"access"`
}

type Fill struct {
	Id                 int         `json:"id"`
	Symbol             string      `json:"symbol"`
	BaseAmount         string      `json:"base_amount"`
	QuoteAmount        string      `json:"quote_amount"`
	Price              string      `json:"price"`
	CreatedAt          time.Time   `json:"created_at"`
	Commission         string      `json:"commission"`
	Side               string      `json:"side"`
	CommissionCurrency string      `json:"commission_currency"`
	OrderId            int         `json:"order_id"`
	Identifier         interface{} `json:"identifier"`
}

type Fills []Fill

func (f *Fills) FilterBySymbol(symbol string) *Fills {
	var filtered Fills
	for _, fill := range *f {
		if fill.Symbol == symbol {
			filtered = append(filtered, fill)
		}
	}
	return &filtered
}

func (f *Fills) FilterBySide(side string) *Fills {
	var filtered Fills
	for _, fill := range *f {
		if fill.Side == side {
			filtered = append(filtered, fill)
		}
	}
	return &filtered
}

func (f *Fills) Before(t time.Time) *Fills {
	var filtered Fills
	for _, fill := range *f {
		if fill.CreatedAt.Before(t) {
			filtered = append(filtered, fill)
		}
	}
	return &filtered
}

func (f *Fills) After(t time.Time) *Fills {
	var filtered Fills
	for _, fill := range *f {
		if fill.CreatedAt.After(t) {
			filtered = append(filtered, fill)
		}
	}
	return &filtered
}

func (f *Fills) Buys() *Fills {
	return f.FilterBySide("buy")
}

func (f *Fills) Sells() *Fills {
	return f.FilterBySide("sell")
}

func (f *Fills) WightedAverage() float64 {
	var totalAmount float64
	var totalValue float64
	for _, fill := range *f {
		amount, _ := strconv.ParseFloat(fill.BaseAmount, 64)
		price, _ := strconv.ParseFloat(fill.Price, 64)
		totalAmount += amount
		totalValue += amount * price
	}
	return totalValue / totalAmount
}

func (f *Fills) TotalAmount() float64 {
	var totalAmount float64
	for _, fill := range *f {
		amount, _ := strconv.ParseFloat(fill.BaseAmount, 64)
		totalAmount += amount
	}
	return totalAmount
}

type WeightedAverageDetail struct {
	TotalAmount, TotalValue, WightedAverage float64
}

func (f *Fills) WightedAverageData() WeightedAverageDetail {
	var totalAmount float64
	var totalValue float64
	for _, fill := range *f {
		amount, _ := strconv.ParseFloat(fill.BaseAmount, 64)
		price, _ := strconv.ParseFloat(fill.Price, 64)
		totalAmount += amount
		totalValue += amount * price
	}
	return WeightedAverageDetail{
		TotalAmount:    totalAmount,
		TotalValue:     totalValue,
		WightedAverage: totalValue / totalAmount,
	}
}

func (f *Fills) MapBySymbol() map[string]Fills {
	separated := make(map[string]Fills)
	for _, fill := range *f {
		separated[fill.Symbol] = append(separated[fill.Symbol], fill)
	}
	return separated
}

func (f *Fills) MapBySymbolAndSide() map[string]map[string]Fills {
	separated := make(map[string]map[string]Fills)
	for _, fill := range *f {
		if _, ok := separated[fill.Symbol]; !ok {
			separated[fill.Symbol] = make(map[string]Fills)
		}
		separated[fill.Symbol][fill.Side] = append(separated[fill.Symbol][fill.Side], fill)
	}
	return separated
}

type Balance struct {
	Id      int    `json:"id"`
	Asset   string `json:"asset"`
	Balance string `json:"balance"`
	Frozen  string `json:"frozen"`
	Service string `json:"service"`
}

type GetBalancesParams struct {
	Assets  []string `json:"assets"`
	Service string   `json:"service"`
	Offset  int      `json:"offset"`
	Limit   int      `json:"limit"`
}

func (p GetBalancesParams) AsURLParams() string {
	var params string
	if len(p.Assets) > 0 {
		params += "assets=" + p.Assets[0]
		for _, asset := range p.Assets[1:] {
			params += "&assets=" + asset
		}
	}
	if p.Service != "" {
		params += "&service=" + p.Service
	}
	if p.Offset != 0 {
		params += "&offset=" + strconv.Itoa(p.Offset)
	}
	if p.Limit != 0 {
		params += "&limit=" + strconv.Itoa(p.Limit)
	}
	if params != "" {
		params = "?" + params
	}

	if params == "?" {
		params = ""
	}

	return params
}

type UserInfo struct {
	UserIdentifier       string        `json:"user_identifier"`
	EnableLimitMsg       bool          `json:"enable_limit_msg"`
	EnableMarketMsg      bool          `json:"enable_market_msg"`
	EnableOcoMsg         bool          `json:"enable_oco_msg"`
	EnableStopLimitMsg   bool          `json:"enable_stop_limit_msg"`
	IsChart              bool          `json:"is_chart"`
	IsClassic            bool          `json:"is_classic"`
	IsEnglishNumber      bool          `json:"is_english_number"`
	IsLight              bool          `json:"is_light"`
	IsVertical           bool          `json:"is_vertical"`
	Phone                string        `json:"phone"`
	Todos                []interface{} `json:"todos"`
	State                string        `json:"state"`
	IsPhoneConfirmed     bool          `json:"is_phone_confirmed"`
	IsEmailConfirmed     bool          `json:"is_email_confirmed"`
	FirstName            string        `json:"first_name"`
	LastName             string        `json:"last_name"`
	Fullname             string        `json:"fullname"`
	BirthDateText        string        `json:"birth_date_text"`
	Email                string        `json:"email"`
	TwoFactorAuthEnabled bool          `json:"two_factor_auth_enabled"`
	TwoFactorAuthType    string        `json:"two_factor_auth_type"`
	Level                struct {
		Id                          int    `json:"id"`
		Title                       string `json:"title"`
		RequiredScore               int64  `json:"required_score"`
		Order                       int    `json:"order"`
		IncomePercentPerTransaction int    `json:"income_percent_per_transaction"`
		MaxDailyWithdraw            int    `json:"max_daily_withdraw"`
	} `json:"level"`
	Type                     string      `json:"type"`
	ReviewStep               int         `json:"review_step"`
	AcceptedStep             int         `json:"accepted_step"`
	RemainingDailyWithdraw   int64       `json:"remaining_daily_withdraw"`
	RemainingMonthlyWithdraw int64       `json:"remaining_monthly_withdraw"`
	Announcement             interface{} `json:"announcement"`
	Tetherban                bool        `json:"tetherban"`
	ForgetPasswordStatus     bool        `json:"forget_password_status"`
	KycStatus                struct {
		State string `json:"state"`
		Code  int    `json:"code"`
		Error string `json:"error"`
	} `json:"kyc_status"`
}

type GetOrdersParams struct {
	Symbol string `json:"symbol"`
	Side   string `json:"side"`
	State  string `json:"state"`
}

func (p *GetOrdersParams) AsURLParams() string {
	var params string
	if p.Symbol != "" {
		params += "symbol=" + p.Symbol
	}
	if p.Side != "" {
		params += "&side=" + p.Side
	}
	if p.State != "" {
		params += "&state=" + p.State
	}
	return params
}

type OrderStatus struct {
	Id                int         `json:"id"`
	Symbol            string      `json:"symbol"`
	Type              string      `json:"type"`
	Side              string      `json:"side"`
	BaseAmount        string      `json:"base_amount"`
	QuoteAmount       string      `json:"quote_amount"`
	Price             string      `json:"price"`
	StopPrice         interface{} `json:"stop_price"`
	OcoTargetPrice    interface{} `json:"oco_target_price"`
	Identifier        string      `json:"identifier"`
	State             string      `json:"state"`
	CreatedAt         time.Time   `json:"created_at"`
	ClosedAt          interface{} `json:"closed_at"`
	DealedBaseAmount  string      `json:"dealed_base_amount"`
	DealedQuoteAmount string      `json:"dealed_quote_amount"`
	ReqToCancel       bool        `json:"req_to_cancel"`
	Commission        string      `json:"commission"`
}

type OrderBook struct {
	Asks [][]string `json:"asks"`
	Bids [][]string `json:"bids"`
}

type Ticker struct {
	Symbol           string  `json:"symbol"`
	Price            string  `json:"price"`
	DailyChangePrice float64 `json:"daily_change_price"`
	Low              string  `json:"low"`
	High             string  `json:"high"`
	Timestamp        float64 `json:"timestamp"`
}

type Trade struct {
	Id          string `json:"id"`
	Price       string `json:"price"`
	BaseAmount  string `json:"base_amount"`
	QuoteAmount string `json:"quote_amount"`
	Side        string `json:"side"`
}

type CreateOrderParams struct {
	Symbol     string `json:"symbol"`
	Type       string `json:"type"`
	Side       string `json:"side"`
	Price      string `json:"price"`
	BaseAmount string `json:"base_amount"`
}
