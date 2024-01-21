package job

import (
	"fmt"
	"verve-hrms/internal/schema"

	"gorm.io/gorm"
)

type Repository interface {
}

type JobRepository struct {
	client *gorm.DB
}

func NewJobRepository(client *gorm.DB) *JobRepository {
	return &JobRepository{client: client}
}

//! Job     ------------------------------------------------------

func (jr JobRepository) JobCreate(newJob *schema.Job) (*schema.Job, error) {
	result := jr.client.Create(newJob)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_create: %w", result.Error)
	}

	return newJob, nil
}

func (jr JobRepository) JobRead(JobID uint) (*schema.Job, error) {
	var job schema.Job

	result := jr.client.First(&job, JobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read: %w", result.Error)
	}

	return &job, nil
}

func (jr JobRepository) ReadAndExpand(JobID uint) (*schema.Job, error) {
	var job schema.Job

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").First(&job, JobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_and_expand: %w", result.Error)
	}

	return &job, nil
}

func (jr JobRepository) JobReadAll() ([]*schema.Job, error) {
	var jobs []*schema.Job

	result := jr.client.Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_all: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_all: %w", ErrEmptyTable)
	}

	return jobs, nil
}

func (jr JobRepository) JobReadAndExpandAll() ([]*schema.Job, error) {
	var jobs []*schema.Job

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_and_expand_all: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_and_expand_all: %w", ErrEmptyTable)
	}

	return jobs, nil
}

func (jr JobRepository) JobUpdate(JobID uint, updateData *schema.Job) (*schema.Job, error) {
	var job *schema.Job

	result := jr.client.Find(&job)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_update: %w", result.Error)
	}

	return job, nil
}

func (jr JobRepository) JobDelete(JobID uint) error {
	var job *schema.Job

	result := jr.client.Delete(&job, JobID)
	if result.Error != nil {
		return fmt.Errorf("job.r.job_delete: %w", result.Error)
	}

	return nil
}

//! AssignedJob ------------------------------------------------------

func (jr JobRepository) AssignedJobCreate(newAssignedJob *schema.AssignedJob) (*schema.AssignedJob, error) {
	result := jr.client.Create(newAssignedJob)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_create: %w", result.Error)
	}

	return newAssignedJob, nil
}

func (jr JobRepository) AssignedJobRead(AssignedJobID uint) (*schema.AssignedJob, error) {
	var assignedJob schema.AssignedJob

	result := jr.client.First(&assignedJob, AssignedJobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_read: %w", result.Error)
	}

	return &assignedJob, nil
}

func (jr JobRepository) AssignedJobReadAll() ([]*schema.AssignedJob, error) {
	var assignedJobs []*schema.AssignedJob

	result := jr.client.Find(&assignedJobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_read_all: %w", result.Error)
	}
	if len(assignedJobs) == 0 {
		return nil, fmt.Errorf("job.r.assigned_job_read_all: %w", ErrEmptyTable)
	}

	return assignedJobs, nil
}

func (jr JobRepository) AssignedJobUpdate(AssignedJobID uint, updateData *schema.AssignedJob) (*schema.AssignedJob, error) {
	var assignedJob *schema.AssignedJob

	result := jr.client.Find(&assignedJob)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_update: %w", result.Error)
	}

	return assignedJob, nil
}

func (jr JobRepository) AssignedJobDelete(AssignedJobID uint) error {
	var assignedJob *schema.AssignedJob

	result := jr.client.Delete(&assignedJob, AssignedJobID)
	if result.Error != nil {
		return fmt.Errorf("job.r.assigned_job_delete: %w", result.Error)
	}

	return nil
}
