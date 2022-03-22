package subreddit

// Service handles domain operations for Subreddit management.
type Service struct {
	*validator
}

// NewService creates and initializes a Subreddit Service.
func NewService(repository Repository) *Service {
	validator := newValidator(repository)

	return &Service{
		validator: validator,
	}
}
