package job

import (
	"fmt"
	"verve-hrms/internal/schema"
)

type Service interface {
}

type JobService struct {
	JobRepository *JobRepository
}

func NewJobService(jobRepository *JobRepository) *JobService {
	return &JobService{JobRepository: jobRepository}
}

//! Job     ------------------------------------------------------

// common
func (js *JobService) ReturnJobList() ([]*schema.Job, error) {

	existingJobs, err := js.JobRepository.JobReadAll()
	if err != nil {
		return nil, fmt.Errorf("job.s.get_job_list: %w", err)
	}

	return existingJobs, nil
}

func (js *JobService) CreateNewJobAndReturnJobList(newJob schema.Job) ([]*schema.Job, error) {

	_, err := js.JobRepository.JobCreate(&newJob)
	if err != nil {
		return nil, fmt.Errorf("job.s.create_new_job_and_return_job_list: %w", err)
	}

	existingJobs, err := js.ReturnJobList()
	if err != nil {
		return nil, fmt.Errorf("job.s.create_new_job_and_return_job_list: %w", err)
	}

	return existingJobs, nil
}

func (js *JobService) UpdateJobAndReturnJobList(jobToUpdate schema.Job) ([]*schema.Job, error) {

	jobID := jobToUpdate.ID

	_, err := js.JobRepository.JobUpdate(jobID, &jobToUpdate)
	if err != nil {
		return nil, fmt.Errorf("job.s.update_job_and_return_job_list: %w", err)
	}

	existingJobs, err := js.ReturnJobList()
	if err != nil {
		return nil, fmt.Errorf("job.s.update_job_and_return_job_list: %w", err)
	}

	return existingJobs, nil
}
