package hyperliquid

import "github.com/ethereum/go-ethereum/signer/core/apitypes"

type AssetInfo struct {
	SzDecimals int
	AssetId    int
}

type OrderRequest struct {
	Coin       string    `json:"coin"`
	IsBuy      bool      `json:"is_buy"`
	Sz         float64   `json:"sz"`
	LimitPx    float64   `json:"limit_px"`
	OrderType  OrderType `json:"order_type"`
	ReduceOnly bool      `json:"reduce_only"`
	Cloid      *string   `json:"cloid,omitempty"`
}

type OrderType struct {
	Limit   *LimitOrderType   `json:"limit" json:",omitempty" msgpack:"limit" msgpack:",omitempty"`
	Trigger *TriggerOrderType `json:"trigger" json:",omitempty"  msgpack:"trigger" msgpack:",omitempty"`
}

type LimitOrderType struct {
	Tif string `json:"tif" msgpack:"tif"`
}

type TriggerOrderType struct {
	IsMarket  bool   `json:"isMarket" msgpack:"isMarket"`
	TriggerPx string `json:"triggerPx" msgpack:"triggerPx"`
	TpSl      TpSl   `json:"tpsl" msgpack:"tpsl"`
}

type TpSl string

const TriggerTp TpSl = "tp"
const TriggerSl TpSl = "sl"

type Grouping string

const GroupingNa Grouping = "na"
const GroupingTpSl Grouping = "positionTpsl"

type RsvSignature struct {
	R string `json:"r"`
	S string `json:"s"`
	V byte   `json:"v"`
}

type ExchangeRequest struct {
	Action       any          `json:"action"`
	Nonce        int64        `json:"nonce"`
	Signature    RsvSignature `json:"signature"`
	VaultAddress *string      `json:"vaultAddress"`
}

type Message struct {
	Source       string `json:"source"`
	ConnectionId []byte `json:"connectionId"`
}

type PlaceOrderAction struct {
	Type     string      `msgpack:"type" json:"type"`
	Orders   []OrderWire `msgpack:"orders" json:"orders"`
	Grouping Grouping    `msgpack:"grouping" json:"grouping"`
}

type UpdateLeverageAction struct {
	Type     string `msgpack:"type" json:"type"`
	Asset    int    `msgpack:"asset" json:"asset"`
	IsCross  bool   `msgpack:"isCross" json:"isCross"`
	Leverage int    `msgpack:"leverage" json:"leverage"`
}

type OrderWire struct {
	Asset      int           `msgpack:"a" json:"a"`
	IsBuy      bool          `msgpack:"b" json:"b"`
	LimitPx    string        `msgpack:"p" json:"p"`
	SizePx     string        `msgpack:"s" json:"s"`
	ReduceOnly bool          `msgpack:"r" json:"r"`
	OrderType  OrderTypeWire `msgpack:"t" json:"t"`
	Cloid      *string       `msgpack:"c,omitempty" json:"c,omitempty"`
}

type OrderTypeWire struct {
	Limit   *LimitOrderType   `json:"limit,omitempty" msgpack:"limit,omitempty"`
	Trigger *TriggerOrderType `json:"trigger,omitempty" msgpack:"trigger,omitempty"`
}

type SigRequest struct {
	PrimaryType string
	DType       []apitypes.Type
	DTypeMsg    map[string]interface{}
}

type CloseRequest struct {
	Coin     string
	Px       *float64
	Sz       *float64
	Slippage *float64
	Cloid    *string
}

type TriggerRequest struct {
	Coin     string
	Px       *float64
	Sz       *float64
	Slippage *float64
	Cloid    *string
	Trigger  TriggerOrderType
}

type OpenRequest struct {
	Coin     string
	IsBuy    bool
	Px       *float64
	Sz       *float64
	Slippage *float64
	Cloid    *string
}

type UpdateLeverageRequest struct {
	Coin     string
	IsCross  bool
	Leverage int
}

type PlaceOrderResponse struct {
	Status   string         `json:"status"`
	Response *InnerResponse `json:"response"`
}

type OrderStatus string

const (
	OrderStatusFilled OrderStatus = "FILLED"
	OrderStatusOpen   OrderStatus = "OPEN"
	OrderStatusFailed OrderStatus = "FAILED"
)

func (r PlaceOrderResponse) GetAvgPrice() *string {
	for _, status := range r.Response.Data.Statuses {
		if status.Filled != nil {
			return &status.Filled.AvgPx
		} else {
			return nil
		}
	}
	return nil
}

func (r PlaceOrderResponse) GetStatus() OrderStatus {
	if r.Status != "ok" {
		return OrderStatusFailed
	}
	for _, status := range r.Response.Data.Statuses {
		if status.Error != nil {
			return OrderStatusFailed
		}
		if status.Filled != nil {
			return OrderStatusFilled
		}
		if status.Resting != nil {
			return OrderStatusOpen
		}
	}
	return OrderStatusOpen
}

type InnerResponse struct {
	Type string       `json:"type"`
	Data DataResponse `json:"data"`
}

type DataResponse struct {
	Statuses []StatusResponse `json:"statuses"`
}

type StatusResponse struct {
	Resting *RestingStatus `json:"resting"`
	Filled  *FilledStatus  `json:"filled"`
	Error   *string        `json:"error"`
}

type RestingStatus struct {
	OrderId string `json:"oid"`
}

type FilledStatus struct {
	OrderId int    `json:"oid"`
	AvgPx   string `json:"avgPx"`
	TotalSz string `json:"totalSz"`
}
