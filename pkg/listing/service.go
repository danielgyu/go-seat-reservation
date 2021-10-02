package listing

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Service struct {
	Conn *sql.DB
}

func (sv *Service) GetAllHalls(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "success")
}
