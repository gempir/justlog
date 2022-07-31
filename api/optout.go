package api

import (
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// swagger:route POST /optout justlog
//
// Generates optout code to use in chat
//
//     Produces:
//     - application/json
//
//     Schemes: https
//
//     Responses:
//       200: string
func (s *Server) writeOptOutCode(w http.ResponseWriter, r *http.Request) {

	code := randomString(6)

	s.bot.OptoutCodes.Store(code, true)
	go func() {
		time.Sleep(time.Second * 60)
		s.bot.OptoutCodes.Delete(code)
	}()

	writeJSON(code, http.StatusOK, w, r)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz1234567890")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
