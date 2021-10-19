package channeling

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserHall struct {
	UserId int
	Hall   string
	writer http.ResponseWriter
}

type Service struct {
	ReservationChan chan *UserHall
}

func (sv *Service) ReserveSeat(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	userId, _ := r.Context().Value("userId").(int)
	hall := param.ByName("hallName")

	if len(sv.ReservationChan) == cap(sv.ReservationChan) {
		fmt.Fprintf(w, "hall is full")
		return
	}
	sv.ReservationChan <- &UserHall{userId, hall, w}
	fmt.Fprintf(w, "reservation success")
}

func Limiter(src chan *UserHall, dst chan *UserHall) {
	userHall := <-src
	if userHall != nil {
		var discard *UserHall
		src <- discard

		dst <- userHall
	}
}
