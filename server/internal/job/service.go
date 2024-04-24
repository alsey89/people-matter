package job

import (
	"fmt"

	"github.com/alsey89/hrms/internal/schema"
)

type Service interface {
}

type JobService struct {
	JobRepository *JobRepository
}

func NewJobService(jobRepository *JobRepository) *JobService {
	return &JobService{JobRepository: jobRepository}
}

//! Position     ------------------------------------------------------

// common
func (js *JobService) ReturnJobListForCompany(companyID uint) ([]*schema.Position, error) {

	existingJobs, err := js.JobRepository.JobReadByCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("job.s.get_job_list: %w", err)
	}

	return existingJobs, nil
}

func (js *JobService) CreateNewJobAndReturnJobList(companyID uint, newJob *schema.Position) ([]*schema.Position, error) {

	_, err := js.JobRepository.JobCreate(newJob)
	if err != nil {
		return nil, fmt.Errorf("job.s.create_new_job_and_return_job_list: %w", err)
	}

	existingJobs, err := js.ReturnJobListForCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("job.s.create_new_job_and_return_job_list: %w", err)
	}

	return existingJobs, nil
}

func (js *JobService) UpdateJobAndReturnJobList(companyID uint, jobID uint, jobToUpdate schema.Position) ([]*schema.Position, error) {

	_, err := js.JobRepository.JobUpdate(jobID, &jobToUpdate)
	if err != nil {
		return nil, fmt.Errorf("job.s.update_job_and_return_job_list: %w", err)
	}

	existingJobs, err := js.ReturnJobListForCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("job.s.update_job_and_return_job_list: %w", err)
	}

	return existingJobs, nil
}

func (js *JobService) DeleteJobAndReturnJobList(companyID uint, jobID uint) ([]*schema.Position, error) {

	err := js.JobRepository.JobDelete(jobID)
	if err != nil {
		return nil, fmt.Errorf("job.s.delete_job_and_return_job_list: %w", err)
	}

	existingJobs, err := js.ReturnJobListForCompany(companyID)
	if err != nil {
		return nil, fmt.Errorf("job.s.delete_job_and_return_job_list: %w", err)
	}

	return existingJobs, nil
}
