package model

import "fmt"

// PresignInput payload for presign request
type PresignInput struct {
	Mime string `json:"mime" validate:"required"`
	Key  string `json:"key,omitempty"`
	Meta struct {
		Title   string  `json:"title,omitempty"`
		OwnerID float64 `json:"owner_id,omitempty"`
	} `json: "meta,omitempty"`
}

// PresignOutput response for presign request
type PresignOutput struct {
	Put string `json:"access_url"`
	Get string `json:"upload_url"`
}

type BadArgError struct {
	Arg     string
	ErrCode string
}

func (e BadArgError) Error() {
	return fmt.Sprintf("Invalid argument %s", e.Arg)
}
