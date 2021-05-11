package shortener

// RedirectService is the interface that reads and saves redirects
type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
