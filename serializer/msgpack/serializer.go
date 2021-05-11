package msgpack

import (
	short "github.com/lrweck/shawty/shortener"
	"github.com/pkg/errors"
	msg "github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*short.Redirect, error) {
	redirect := &short.Redirect{}
	if err := msg.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *short.Redirect) ([]byte, error) {
	rawMsg, err := msg.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}
