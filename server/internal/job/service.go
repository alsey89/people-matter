package job

type Service interface {
}

type JobService struct {
	JobRepository JobRepository
}

func NewJobService(jobRepository JobRepository) *JobService {
	return &JobService{JobRepository: jobRepository}
}

//! Job     ------------------------------------------------------
