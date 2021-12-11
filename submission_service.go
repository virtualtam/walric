package redwall

var _ SubmissionRepository = &SubmissionService{}

type SubmissionService struct {
	SubmissionRepository
}

func NewSubmissionService(submissionRepository SubmissionRepository) *SubmissionService {
	validator := newSubmissionValidator(submissionRepository)

	return &SubmissionService{
		SubmissionRepository: validator,
	}
}
