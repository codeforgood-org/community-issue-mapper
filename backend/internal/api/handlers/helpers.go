package handlers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func getIDFromPath(r *http.Request, param string) (int64, error) {
	vars := mux.Vars(r)
	return strconv.ParseInt(vars[param], 10, 64)
}
