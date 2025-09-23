package stateUpdaterStructures

type ItemResultMain struct {
	Result ItemResult `json:"result"`
}

type ItemResult struct {
	Data Data `json:"data"`
}

type Data struct {
	Name string `json:"name"`

	Article int    `json:"uin"`
	Uinsql  string `json:"uinsql"`

	Price      int `json:"price"`
	EilatPrice int `json:"eilatPrice"`

	MinPrice      int `json:"min_price"`
	MinEilatPrice int `json:"min_eilat_price"`
}
