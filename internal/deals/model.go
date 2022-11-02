package deals

import (
	"errors"
	"time"
)

//(вес груза, объем груза, дата+время погрузки, дата+время выгрузки, адрес погрузки, адрес выгрузки, комментарий к заказу)
type Deal struct {
	Id             int       `json:"id"`
	User           string    `json:"user"`
	Weight         int       `json:"weight"`
	Volume         int       `json:"volume"`
	LoadDateTime   time.Time `json:"load_date_time"`
	UnLoadDateTime time.Time `json:"unload_date_time"`
	LoadAdress     string    `json:"load_adress"`
	UnloadAdress   string    `json:"unload_adress"`
	Comment        string    `json:"comment"`
}

type ListObj struct {
	Deals   []Deal `json:"deals"`
	Total   int    `json:"total"`
	Page    int    `json:"page"`
	PerPage int    `json:"perPage"`
}

func (d *Deal) Validate() error {
	// Ну допустим проверим только вес объем и адреса без конекретной ошибки
	if d.Weight == 0 || d.Volume == 0 || d.LoadAdress == "" || d.UnloadAdress == "" {
		return errors.New("deal is not valid")
	}
	return nil
}
