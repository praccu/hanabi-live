package sessions

import (
	"nhooyr.io/websocket"
)

type deleteData struct {
	userID   int
	username string
}

func (m *Manager) Delete(s *session) {
	if m.requestsClosed.IsSet() {
		return
	}

	m.requests <- &request{
		Type: requestTypeDelete,
		Data: &deleteData{
			userID:   s.userID,
			username: s.username,
		},
	}
}

func (m *Manager) delete(data interface{}) {
	var d *deleteData
	if v, ok := data.(*deleteData); !ok {
		m.logger.Errorf("Failed type assertion for data of type: %T", d)
		return
	} else {
		d = v
	}

	var s *session
	if v, ok := m.sessions[d.userID]; !ok {
		// The session for this user does not exist in the sessions map,
		// so we don't have to do anything
		return
	} else {
		s = v
	}

	// Close the existing connection
	// This will be a no-op if it is already closed
	// We do this in a new goroutine to avoid blocking
	go s.conn.Close(websocket.StatusNormalClosure, "")

	// Delete the entry from the sessions map
	delete(m.sessions, d.userID)
	logSession(m, d.userID, d.username, false)

	// Remove this player from any tables
	m.Dispatcher.Tables.DisconnectUser(d.userID)

	// Alert everyone that a user has logged out
	m.NotifyAllUserLeft(s.userID)
}