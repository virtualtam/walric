package subreddit

// Service handles domain operations for Subreddit management.
type Service struct {
	Repository
}

// NewService creates and initializes a Subreddit Service.
func NewService(repository Repository) *Service {
	validator := newValidator(repository)

	return &Service{
		Repository: validator,
	}
}
