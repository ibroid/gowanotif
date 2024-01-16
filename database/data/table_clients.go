package data

type Clients struct {
	Id         int    `json:"id"`
	ClientName string `json:"client_name"`
	Jid        string `json:"jid"`
	Handler    int    `josn:"handler"`
	Status     int    `json:"status"`
	Service    string `json:"service"`
}
