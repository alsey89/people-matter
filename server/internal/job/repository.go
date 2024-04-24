package job

import (
	"fmt"
	"log"

	"github.com/alsey89/hrms/internal/schema"

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

//! Position     ------------------------------------------------------

func (jr JobRepository) JobCreate(newJob *schema.Position) (*schema.Position, error) {
	log.Printf("job.r.job_create: newJob: %v", newJob)

	result := jr.client.Create(newJob)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_create: %w", result.Error)
	}

	return newJob, nil
}

func (jr JobRepository) JobRead(jobID uint) (*schema.Position, error) {
	var job schema.Position

	result := jr.client.First(&job, jobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read: %w", result.Error)
	}

	return &job, nil
}

func (jr JobRepository) JobReadAndExpand(jobID uint) (*schema.Position, error) {
	var job schema.Position

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").First(&job, jobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_and_expand: %w", result.Error)
	}

	return &job, nil
}

func (jr JobRepository) JobReadByCompany(companyID uint) ([]*schema.Position, error) {
	var jobs []*schema.Position

	result := jr.client.Where("company_id = ?", companyID).Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_by_company: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_by_company: %w", gorm.ErrRecordNotFound)
	}

	return jobs, nil
}

func (jr JobRepository) JobReadByCompanyAndExpand(companyID uint) ([]*schema.Position, error) {
	var jobs []*schema.Position

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").Where("company_id = ?", companyID).Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_by_company_and_expand: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_by_company_and_expand: %w", ErrNoRowsFound)
	}

	return jobs, nil
}

func (jr JobRepository) JobReadAll() ([]*schema.Position, error) {
	var jobs []*schema.Position

	result := jr.client.Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_all: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_all: %w", ErrEmptyTable)
	}

	return jobs, nil
}

func (jr JobRepository) JobReadAndExpandAll() ([]*schema.Position, error) {
	var jobs []*schema.Position

	result := jr.client.Preload("Subordinates").Preload("AssignedJobs").Find(&jobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_read_and_expand_all: %w", result.Error)
	}
	if len(jobs) == 0 {
		return nil, fmt.Errorf("job.r.job_read_and_expand_all: %w", ErrEmptyTable)
	}

	return jobs, nil
}

func (jr JobRepository) JobUpdate(jobID uint, updateData *schema.Position) (*schema.Position, error) {
	var job schema.Position

	updateDataMap := map[string]interface{}{
		"Title":          updateData.Title,
		"Description":    updateData.Description,
		"Duties":         updateData.Duties,
		"Qualifications": updateData.Qualifications,
		"Experience":     updateData.Experience,
		"MinSalary":      updateData.MinSalary,
		"MaxSalary":      updateData.MaxSalary,

		// "CompanyID": updateData.CompanyID, //* not part of updateData, will always initialize to 0 since it's not a pointer
		"LocationID":   updateData.LocationID,
		"DepartmentID": updateData.DepartmentID,
		"ManagerID":    updateData.ManagerID,
	}

	//*Updates only updates non-zero values, need to select nil values explicitly
	result := jr.client.Model(&job).Where("ID = ?", jobID).Updates(updateDataMap)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.job_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the job does not exist.
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("job.r.job_update: %w", gorm.ErrRecordNotFound)
	}

	return &job, nil
}

func (jr JobRepository) JobDelete(jobID uint) error {
	var job *schema.Position

	result := jr.client.Unscoped().Delete(&job, jobID)
	if result.Error != nil {
		return fmt.Errorf("job.r.job_delete: %w", result.Error)
	}

	return nil
}

//! UserPosition ------------------------------------------------------

func (jr JobRepository) AssignedJobCreate(newAssignedJob *schema.UserPosition) (*schema.UserPosition, error) {
	result := jr.client.Create(newAssignedJob)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_create: %w", result.Error)
	}

	return newAssignedJob, nil
}

func (jr JobRepository) AssignedJobRead(assignedJobID uint) (*schema.UserPosition, error) {
	var assignedJob schema.UserPosition

	result := jr.client.First(&assignedJob, assignedJobID)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_read: %w", result.Error)
	}

	return &assignedJob, nil
}

func (jr JobRepository) AssignedJobReadAll() ([]*schema.UserPosition, error) {
	var assignedJobs []*schema.UserPosition

	result := jr.client.Find(&assignedJobs)
	if result.Error != nil {
		return nil, fmt.Errorf("job.r.assigned_job_read_all: %w", result.Error)
	}
	if len(assignedJobs) == 0 {
		return nil, fmt.Errorf("job.r.assigned_job_read_all: %w", ErrEmptyTable)
	}

	return assignedJobs, nil
}

func (jr JobRepository) AssignedJobUpdate(assignedJobID uint, updateData *schema.UserPosition) error {
	var assignedJob *schema.UserPosition

	result := jr.client.Model(&assignedJob).Where("ID = ?", assignedJobID).Updates(updateData)
	if result.Error != nil {
		return fmt.Errorf("job.r.assigned_job_update: %w", result.Error)
	}

	// Check if any row was affected, if not, the assignedJob does not exist.
	if result.RowsAffected == 0 {
		return fmt.Errorf("job.r.assigned_job_update: %w", gorm.ErrRecordNotFound)
	}

	return nil
}

func (jr JobRepository) AssignedJobDelete(AssignedJobID uint) error {
	var assignedJob *schema.UserPosition

	result := jr.client.Unscoped().Delete(&assignedJob, AssignedJobID)
	if result.Error != nil {
		return fmt.Errorf("job.r.assigned_job_delete: %w", result.Error)
	}

	return nil
}
