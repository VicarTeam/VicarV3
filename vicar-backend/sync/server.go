package sync

import (
	"os"
	"vicar-backend/auth"
	"vicar-backend/log"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	clientsMap = make(map[*websocket.Conn]syncClient)
)

func Register(app fiber.Router) {
	cfg := websocket.Config{
		RecoverHandler: func(c *websocket.Conn) {
			if err := recover(); err != nil {
				log.Error(log.Sync, "❌", "WebSocket panic recovered: %v", err)
			}
		},
		Origins: []string{os.Getenv("CORS_ALLOW_ORIGINS")},
	}

	app.Get("/sync/:type", websocket.New(handleSync, cfg))
}

func handleSync(c *websocket.Conn) {
	typStr := c.Params("type")
	typ := parseClientType(typStr)
	if typ == clientTypeUnknown {
		log.Error(log.Sync, "❌", "Unknown client type: %s", typStr)
		c.Close()
		return
	}

	client := newSyncClient(c, typ)
	if typ == clientTypeVicar {
		user := auth.ExtractForWebSocket(c)
		if user == nil {
			log.Error(log.Sync, "❌", "Unauthorized Vicar sync client")
			c.Close()
			return
		}
		client.user = user
	}

	clientsMap[c] = client
	defer delete(clientsMap, c)

	client.listen()
}

func SyncCharacterChanges(charID string, data any) {
	for _, client := range clientsMap {
		if client.typ != clientTypeVicar {
			continue
		}

		if client.user == nil {
			continue
		}

		if _, ok := client.viewing[charID]; !ok {
			continue
		}

		if err := client.sendPacket("UPDATE_CHARACTER", map[string]any{
			"characterID": charID,
			"change":      data,
		}); err != nil {
			log.Error(log.Sync, "❌", "Failed to send character update to client: %v", err)
		}
	}
}
