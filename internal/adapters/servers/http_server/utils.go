package http_server

import (
	"encoding/json"
	"errors"
	"github.com/assizkii/messaggio/internal/domain/interfaces"
	"net/http"
	"regexp"
)

func prepareRequestData(r *http.Request) (interfaces.Message, error) {

	var message interfaces.Message

	message.Phone = r.FormValue("phone")
	message.Text = r.FormValue("text")

	if err := phoneValidate(message.Phone); err != nil {
		return message, err
	}

	return message, nil
}

func phoneValidate(phone string) error {

	phonePattern := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	if !phonePattern.MatchString(phone) {
		return errors.New("phone number is not valid")
	}

	return nil
}

func showResponse(result *HttpResponse, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(result.status)
	w.Write(response)
}
