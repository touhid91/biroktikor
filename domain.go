package main

import "fmt"

// PresignInput Value object for presign input
type PresignInput struct {
	Mime string `json:"mime" validate:"required"`
	Key  string `json:"key,omitempty"`
	Meta struct {
		Title   string  `json:"title,omitempty"`
		OwnerID float64 `json:"owner_id,omitempty"`
	}
}

// PresignOutput Value object for presign output
type PresignOutput struct {
	Put string `json:"upload_url"`
	Get string `json:"access_url"`
}

// BadArgError Value object for invalid argument
type BadArgError struct {
	Arg     string
	ErrCode string
}

func (e BadArgError) Error() string {
	return fmt.Sprintf("Bad argument %s, code %s", e.Arg, e.ErrCode)
}
