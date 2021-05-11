package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	val "gopkg.in/dealancer/validate.v2"
)

var (
	// ErrRedirectNotFound is returned when a code does not point to a existing redirect
	ErrRedirectNotFound = errors.New("Redirect Not Found")

	// ErrRedirectInvalid is returned when a url is not valid
	ErrRedirectInvalid = errors.New("Redirect Invalid")
)

type redirectService struct {
	repo RedirectRepository
}

// NewRedirectService creates a new RedirectService with the specified repo
func NewRedirectService(repo RedirectRepository) RedirectService {
	return &redirectService{
		repo,
	}
}

// Find utilizes the specified repo to locate url redirects
func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.repo.Find(code)
}

// Store utilizes the specified repo to persist urls to redirects
func (r *redirectService) Store(redirect *Redirect) error {
	if err := val.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.repo.Store(redirect)
}
