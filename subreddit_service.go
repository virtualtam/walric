package redwall

type SubredditService struct {
	SubredditRepository
}

func NewSubredditService(subredditRepository SubredditRepository) *SubredditService {
	validator := newSubredditValidator(subredditRepository)

	return &SubredditService{
		SubredditRepository: validator,
	}
}
