package shortener

// RedirectSerializer - Interface for serializers. Currently only json and msgpack are supported
type RedirectSerializer interface {
	Decode(input []byte) (*Redirect, error)
	Encode(input *Redirect) ([]byte, error)
}
