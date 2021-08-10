package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"phenix/store"

	log "github.com/activeshadow/libminimega/minilog"
)

type WebError struct {
	*store.Event

	Cause  error  `json:"-"`
	Status int    `json:"-"`
	URL    string `json:"url"`

	UserMetadata map[string]string `json:"metadata,omitempty"`
}

func NewWebError(cause error, format string, args ...interface{}) *WebError {
	event := store.NewErrorEvent(fmt.Errorf(format, args...))

	err := &WebError{
		Event:  event.WithMetadata("cause", cause.Error()),
		Cause:  cause,
		Status: http.StatusBadRequest,
		URL:    "/api/v1/errors/" + event.ID,
	}

	return err
}

func (this *WebError) WithMetadata(k, v string, user bool) *WebError {
	this.Event.WithMetadata(k, v)

	if user {
		if this.UserMetadata == nil {
			this.UserMetadata = make(map[string]string)
		}

		this.UserMetadata[k] = v
	}

	return this
}

func (this *WebError) SetInformational() *WebError {
	this.Event.Type = store.EventTypeInfo
	return this
}

func (this *WebError) SetStatus(status int) *WebError {
	this.Status = status
	return this
}

func (this WebError) Error() string {
	if this.Cause == nil {
		return this.Event.Message
	}

	return fmt.Sprintf("%s: %v", this.Event.Message, this.Cause)
}

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func (this ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := this(w, r); err != nil {
		web, ok := err.(*WebError)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		go store.AddEvent(*web.Event)

		web.Event.Metadata = nil

		body, _ := json.Marshal(web)
		log.Infoln(string(body))

		w.WriteHeader(web.Status)
		w.Write(body)
	}
}
