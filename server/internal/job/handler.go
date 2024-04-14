package job

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"github.com/alsey89/hrms/internal/common"
	"github.com/alsey89/hrms/internal/schema"
)

type JobHandler struct {
	JobService *JobService
}

func NewJobHandler(jobService *JobService) *JobHandler {
	return &JobHandler{JobService: jobService}
}

// ! Position ------------------------------------------------------------
func (jh *JobHandler) CreateJob(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("job.h.create_job: company_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("job.h.create_job: error converting company_id to uint: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	newJob := new(schema.Position)
	err = c.Bind(newJob)
	if err != nil {
		log.Printf("job.h.create_job: error binding job data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	jobData, err := jh.JobService.CreateNewJobAndReturnJobList(uintCompanyID, newJob)
	if err != nil {
		log.Printf("job.h.create_job: %v", err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return c.JSON(http.StatusNotFound, common.APIResponse{
				Message: "job already exists",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error creating job data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "job data has been created",
		Data:    jobData,
	})
}

func (jh *JobHandler) GetAllJobs(c echo.Context) error {
	stringCompanyID := c.Param("company_id")
	if stringCompanyID == "" {
		log.Printf("job.h.read_all_jobs: company_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyID)
	if err != nil {
		log.Printf("job.h.read_all_jobs: error converting company_id to uint: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	jobData, err := jh.JobService.ReturnJobListForCompany(uintCompanyID)
	if err != nil {
		log.Printf("job.h.read_all_jobs: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNoContent, common.APIResponse{
				Message: "no job data found",
				Data:    nil,
			})
		}
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error reading all job data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "job list has been fetched",
		Data:    jobData,
	})
}

func (jh *JobHandler) UpdateJob(c echo.Context) error {
	jobToUpdate := new(schema.Position)
	err := c.Bind(jobToUpdate)
	if err != nil {
		log.Printf("job.h.update_job: error binding job data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}
	if jobToUpdate == nil {
		log.Printf("job.h.update_job: error binding job data: %v", err)
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid update data",
			Data:    nil,
		})
	}

	stringJobId := c.Param("job_id")
	if stringJobId == "" {
		log.Printf("job.h.update_job: job_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid job id",
			Data:    nil,
		})
	}

	uintJobID, err := common.ConvertStringOfNumbersToUint(stringJobId)
	if err != nil {
		log.Printf("job.h.delete_job: error converting job_id to uint: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	stringCompanyId := c.Param("company_id")
	if stringCompanyId == "" {
		log.Printf("job.h.update_job: company_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid company id",
			Data:    nil,
		})
	}

	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyId)
	if err != nil {
		log.Printf("job.h.delete_job: error converting company_id to uint: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	jobData, err := jh.JobService.UpdateJobAndReturnJobList(uintCompanyID, uintJobID, *jobToUpdate)

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "job data has been updated",
		Data:    jobData,
	})
}

func (jh *JobHandler) DeleteJob(c echo.Context) error {
	stringJobId := c.Param("job_id")
	if stringJobId == "" {
		log.Printf("job.h.delete_job: job_id is empty")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "invalid job id",
			Data:    nil,
		})
	}

	uintJobID, err := common.ConvertStringOfNumbersToUint(stringJobId)
	if err != nil {
		log.Printf("job.h.delete_job: error converting job_id to uint: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	stringCompanyId := c.Param("company_id")
	uintCompanyID, err := common.ConvertStringOfNumbersToUint(stringCompanyId)
	if err != nil {
		log.Printf("job.h.delete_job: error converting company_id to uint: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}

	jobData, err := jh.JobService.DeleteJobAndReturnJobList(uintCompanyID, uintJobID)
	if err != nil {
		log.Printf("job.h.delete_job: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "error deleting job data",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, common.APIResponse{
		Message: "job data has been deleted",
		Data:    jobData,
	})
}
