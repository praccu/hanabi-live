package tables

type getUserTablesData struct {
	userID         int
	resultsChannel chan []uint64
}

// GetUserTables requests a list of all of the tables that a user is currently playing in and
// spectating.
func (m *Manager) GetUserTables(userID int) ([]uint64, []uint64) {
	resultsChannel := make(chan []uint64)

	m.requests <- &request{
		Type: requestTypeGetUserTables,
		Data: &getUserTablesData{
			userID:         userID,
			resultsChannel: resultsChannel,
		},
	}

	playingAtTables := <-resultsChannel
	spectatingTables := <-resultsChannel

	return playingAtTables, spectatingTables
}

func (m *Manager) getUserTables(data interface{}) {
	var d *getUserTablesData
	if v, ok := data.(*getUserTablesData); !ok {
		m.logger.Errorf("Failed type assertion for data of type: %T", d)
		return
	} else {
		d = v
	}

	d.resultsChannel <- m.getUserPlaying(d.userID)
	d.resultsChannel <- m.getUserSpectating(d.userID)
}