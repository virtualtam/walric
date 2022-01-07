package subreddit

type Service struct {
	Repository
}

func NewService(repository Repository) *Service {
	validator := newValidator(repository)

	return &Service{
		Repository: validator,
	}
}
