package pageHandler

import (
	"github.com/gorilla/mux"
	"golang.captainalm.com/cityuni-webserver/conf"
	"net/http"
)

var theRouter *mux.Router

func GetRouter(config conf.ConfigYaml) http.Handler {
	if theRouter == nil {
		theRouter = mux.NewRouter()
		//Mux routing stuff
	}
	return theRouter
}
