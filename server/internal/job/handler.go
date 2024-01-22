package job

import (
	"errors"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"verve-hrms/internal/common"
	"verve-hrms/internal/schema"
)

type JobHandler struct {
	JobService *JobService
}

func NewJobHandler(jobService *JobService) *JobHandler {
	return &JobHandler{JobService: jobService}
}

// ! Job ------------------------------------------------------------
func (jh *JobHandler) CreateJob(c echo.Context) error {
	newJob := new(schema.Job)

	err := c.Bind(newJob)
	if err != nil {
		log.Printf("job.h.create_job: error binding job data: %v", err)
		return c.JSON(http.StatusInternalServerError, common.APIResponse{
			Message: "something went wrong",
			Data:    nil,
		})
	}
	if newJob == nil {
		log.Printf("job.h.create_job: new job data is nil")
		return c.JSON(http.StatusBadRequest, common.APIResponse{
			Message: "no incoming job data",
			Data:    nil,
		})
	}

	jobData, err := jh.JobService.CreateNewJobAndReturnJobList(*newJob)
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
	jobData, err := jh.JobService.ReturnJobList()
	if err != nil {
		log.Printf("job.h.read_all_jobs: %v", err)
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
