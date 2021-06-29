package frontend

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

type flashLevel uint8

const (
	flashLevelInfo flashLevel = iota
	flashLevelWarning
	flashLevelDanger
	flashLevelSuccess

	flashKeyMessage = "message"
)

func (fl flashLevel) String() string {
	switch fl {
	case flashLevelInfo:
		return "info"
	case flashLevelWarning:
		return "warning"
	case flashLevelDanger:
		return "danger"
	case flashLevelSuccess:
		return "success"
	}

	log.Error().Msg("could not get flashlevel string")
	return "info"
}

type flashMessage struct {
	Level   flashLevel
	Message string
}

func (w *Web) addFlash(rw http.ResponseWriter, r *http.Request, message flashMessage) {
	session, err := w.cookieStore.Get(r, flashSession)
	if err != nil {
		log.Error().Err(err).Msg("could not get flash from cookie")
		return
	}

	session.AddFlash(message, flashKeyMessage)
	if err := session.Save(r, rw); err != nil {
		log.Error().Err(err).Msg("could not save session")
		return
	}
}

func (w *Web) getFlash(rw http.ResponseWriter, r *http.Request) []flashMessage {
	session, err := w.cookieStore.Get(r, flashSession)
	if err != nil {
		log.Error().Err(err).Msg("could not get flash from cookie")
		return nil
	}

	flashes := session.Flashes(flashKeyMessage)

	var flashMessages []flashMessage
	for _, v := range flashes {
		v, ok := v.(flashMessage)
		if !ok {
			continue
		}

		flashMessages = append(flashMessages, v)
	}

	if err := session.Save(r, rw); err != nil {
		log.Error().Err(err).Msg("could not save session")
		return nil
	}

	return flashMessages
}
