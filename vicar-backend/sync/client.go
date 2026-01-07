package sync

import (
	"vicar-backend/db"
	"vicar-backend/db/entities"
	"vicar-backend/log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
)

type clientType int

const (
	clientTypeVicar clientType = iota
	clientTypeVTT
	clientTypeUnknown
)

func parseClientType(s string) clientType {
	switch s {
	case "vicar":
		return clientTypeVicar
	case "vtt":
		return clientTypeVTT
	default:
		return clientTypeUnknown
	}
}

type syncClient struct {
	typ     clientType
	conn    *websocket.Conn
	user    *entities.User
	viewing map[string]bool // characters being viewed
}

func newSyncClient(c *websocket.Conn, typ clientType) syncClient {
	return syncClient{
		typ:  typ,
		conn: c,
	}
}

func (c *syncClient) listen() {
	for {
		mt, b, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(log.Sync, "❌", "WebSocket read error: %v", err)
			break
		}

		if mt == websocket.CloseMessage {
			log.Info(log.Sync, "ℹ️", "WebSocket closed by client")
			break
		}

		if mt == websocket.TextMessage {
			msg := string(b)
			log.Info(log.Sync, "ℹ️", "Received message: %s", msg)

			packet := map[string]any{}
			if err := json.Unmarshal(b, &packet); err != nil {
				log.Error(log.Sync, "❌", "Failed to unmarshal message: %v", err)
				continue
			}

			c.handlePacket(packet)
		}
	}
}

func (c *syncClient) sendPacket(packetType string, packet map[string]any) error {
	packet["type"] = packetType
	data, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *syncClient) handlePacket(packet map[string]any) {
	typ, ok := packet["type"].(string)
	if !ok {
		log.Error(log.Sync, "❌", "Packet missing type field")
		return
	}

	switch typ {
	case "VIEW_CHARACTER":
		charID, ok := packet["characterID"].(string)
		if !ok {
			log.Error(log.Sync, "❌", "VIEW_CHARACTER packet missing characterID field")
			return
		}

		hasAccess := false
		if c.typ == clientTypeVicar && c.user != nil {
			var char entities.V5Character
			err := db.DB.Where("id = ?", charID).First(&char).Error
			if err == nil && char.UserID == c.user.ID {
				hasAccess = true
			} else {
				var count int64
				err = db.DB.Raw("SELECT 1 FROM v5_character_viewers WHERE v5_character_id = ? AND user_id = ? LIMIT 1", charID, c.user.ID).Count(&count).Error
				if err == nil && count > 0 {
					hasAccess = true
				}
			}
		}

		if !hasAccess {
			log.Error(log.Sync, "❌", "Client does not have access to view character: %s", charID)
			return
		}

		c.viewing[charID] = true
		log.Info(log.Sync, "ℹ️", "Client is now viewing character: %s", charID)
	case "UNVIEW_CHARACTER":
		charID, ok := packet["characterID"].(string)
		if !ok {
			log.Error(log.Sync, "❌", "UNVIEW_CHARACTER packet missing characterID field")
			return
		}
		delete(c.viewing, charID)
		log.Info(log.Sync, "ℹ️", "Client stopped viewing character: %s", charID)
	default:
		log.Error(log.Sync, "❌", "Unknown packet type: %s", typ)
	}
}
