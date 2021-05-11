package shortener

// RedirectRepository - Interface for all storage types. It should at least find and store redirects
type RedirectRepository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
