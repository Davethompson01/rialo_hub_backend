package websocket

import (
	"github.com/Davethompson01/rialo_hub_backend/config"
)

func NewHub() *config.Hub {
	return &config.Hub{
		Clients: make(map[int]*config.Client),
	}
}
