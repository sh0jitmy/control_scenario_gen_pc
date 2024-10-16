package domainsim

import (
	//"encoding/json"
	"fmt"
)

// Entity1の定義
type Entity1 struct {
	ID   string `json:"id"`
	User string `json:"user"`
}

// Entity2の定義
type Entity2 struct {
	ID       string `json:"id"`
	BuyCount uint16 `json:"buy_count"`
	TotalBuy uint32 `json:"total_buy"`
}

// Entity3の定義
type Entity3 struct {
	User     string `json:"user"`
	BuyCount uint16 `json:"buy_count"`
	TotalBuy uint32 `json:"total_buy"`
}

// EntityControlの定義
type EntityControl struct {
	Parsemap map[string]func(any)(any,error)
	Entity2 Entity2
}

func NewEntityControl() *EntityControl {
	pm := make(map[string]func(any)(any,error))
	e2 := Entity2{"ID1",1000,10000}
 
	ec := EntityControl{
		//Entity2:{ ID:"ID1",BuyCount:1000,TotalBuy:10000},
		//Parsemap: make(map[string]func(any)(any,error)),
		Parsemap: pm,
		Entity2 : e2,
	}
	ec.Parsemap["testhandle1"] = ec.CreateEntity3
	return &ec
}


func (ec *EntityControl) ProcessData(key string, data any) (any, error) {
	if funcp,ok :=ec.Parsemap[key]; ok {
		return funcp(data)
	} 
	return "Notfound",fmt.Errorf("Not Found")
}

// CreateEntity3をEntityControlのメソッドにする
func (ec *EntityControl) CreateEntity3(ia1 any) (any, error) {
	
	et1 , ok := ia1.(Entity1)
	if ! ok {
		return Entity3{}, fmt.Errorf("input cast error")
	}
	// IDが一致するかチェック
	if et1.ID != ec.Entity2.ID {
		return Entity3{}, fmt.Errorf("ID mismatch: Entity1 ID %s, Entity2 ID %s", et1.ID, ec.Entity2.ID)
	}

	// Entity3の作成
	entity3 := Entity3{
		User:     et1.User,
		BuyCount: ec.Entity2.BuyCount,
		TotalBuy: ec.Entity2.TotalBuy,
	}

	return entity3, nil
}

