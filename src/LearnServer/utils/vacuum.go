package utils

import (
	"fmt"
	"time"
)

var Manager = NewVacuum()

var itemTTL = time.Minute * 5

type TempItem struct {
	expiration time.Time
	payload    interface{}
}

type Vacuum struct {
	items map[string]TempItem
}

func NewVacuum() *Vacuum {
	v := Vacuum{}
	go v.clean()
	return &v
}

func (va *Vacuum) clean() {
	ticker := time.NewTicker(time.Minute)
	for {
		select {
		case <-ticker.C:
			now := time.Now()
			for k, v := range va.items {
				if v.expiration.After(now) {
					continue
				}
				delete(va.items, k)
			}
		}
	}
}

func (v *Vacuum) Insert(data interface{}) string {
	uid, _ := UUID()
	if v.items == nil {
		v.items = make(map[string]TempItem)
	}
	v.items[uid] = TempItem{
		expiration: time.Now().Add(itemTTL),
		payload:    data,
	}
	return uid
}

func (va *Vacuum) Delete(id string) {
	delete(va.items, id)
}

func (va *Vacuum) Get(id string) (interface{}, error) {
	for k, v := range va.items {
		if k == id {
			delete(va.items, k)
			return v.payload, nil
		}
	}
	return nil, fmt.Errorf("No such id")
}
