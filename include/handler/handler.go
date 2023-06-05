package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Idibrov/go-final/pkg/sps"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
)

func InitRoutes() *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", handleConnection)
	return router
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	resultT := sps.ResultT{}

	simUrl := viper.GetString("simulator.addr") + ":" + viper.GetString("simulator.port")
	filePath := viper.GetString("data.path")

	resultSetT, err := sps.GetResultData("http://"+simUrl, filePath)
	if err != nil {
		resultT.Status = false
		resultT.Error = err.Error()
	} else {
		resultT.Status = true
		resultT.Data = resultSetT
		resultT.Error = ""
	}

	message, err := json.Marshal(resultT)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte{})
		if err != nil {
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(message)
	if err != nil {
		return
	}
}
