package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func s3Handler(w http.ResponseWriter, r *http.Request) {
	if "POST" != r.Method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, `{"error":%q}`, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	command := PresignInput{}
	if err := json.NewDecoder(r.Body).Decode(&command); nil != err {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":%q}`, err)
		return
	}

	// command.Meta.OwnerID = r.Context().Value("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)

	reply, err := Presign(&command)
	if nil != err {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error":%q}`, err)
		return
	}

	mar, err := json.Marshal(reply)
	if nil != err {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(mar)
}

func main() {
	http.HandleFunc("/storage/s3", s3Handler)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
