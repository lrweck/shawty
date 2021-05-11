package json

import (
	"encoding/json"

	short "github.com/lrweck/shawty/shortener"
	"github.com/pkg/errors"
)

// Redirect struct to add methods to.
type Redirect struct{}

// Decodes bytes into a Redirect struct via json
func (r *Redirect) Decode(input []byte) (*short.Redirect, error) {
	redirect := &short.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

// Encodes a Redirect struct to json bytes
func (r *Redirect) Encode(input *short.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}
