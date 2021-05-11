package shortener

// The app service has to provide Find and Store methods
type RedirectService interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
