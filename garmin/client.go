package garmin

import (
	connect "github.com/abrander/garmin-connect"
	"github.com/rs/zerolog"
)

func NewClient(username, password string, logger zerolog.Logger) *connect.Client {

	opt := connect.Credentials(username, password)
	return connect.NewClient(opt)
}
