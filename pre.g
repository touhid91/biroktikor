package storage

import (
	"app/core"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// PresignHandler GET presign
func PresignHandler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	command := model.PresignInput{}

	if err := json.NewDecoder(r.Body).Decode(&command); nil != err {
		core.HttpErrResponse(&w, err, http.StatusBadRequest)
		return
	}

	if len(command.Mime) < 1 {
		// TODO add support for accepted mimes only
		core.HttpErrResponse(&w, "invalid request parameter: mime", http.StatusBadRequest)
		return
	}

	/* Get context user */
	context := r.Context().Value("user")
	ctxUser := context.(*jwt.Token).Claims.(jwt.MapClaims)
	command.Meta.OwnerID = uint(ctxUser["id"].(float64))

	url, err := Presign(&command)
	if nil != err {
		fmt.Print(err)
	}

	io.WriteString(w, url)
}
