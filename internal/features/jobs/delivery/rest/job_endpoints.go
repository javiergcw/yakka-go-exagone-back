package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	"github.com/yakka-backend/internal/features/jobs/models"
	"github.com/yakka-backend/internal/features/jobs/payload"
	"github.com/yakka-backend/internal/features/jobs/usecase"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

type JobHandler struct {
	jobUsecase         usecase.JobUsecase
	builderProfileRepo builder_db.BuilderProfileRepository
}

func NewJobHandler(jobUsecase usecase.JobUsecase, builderProfileRepo builder_db.BuilderProfileRepository) *JobHandler {
	return &JobHandler{
		jobUsecase:         jobUsecase,
		builderProfileRepo: builderProfileRepo,
	}
}

// CreateJob creates a new job
func (h *JobHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Invalid user ID")
		return
	}

	// Get builder profile ID from user ID
	builderProfile, err := h.builderProfileRepo.GetByUserID(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Builder profile not found")
		return
	}

	var req payload.CreateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Set builder profile ID from JWT
	req.BuilderProfileID = builderProfile.ID

	// Create job
	job, err := h.jobUsecase.CreateJob(r.Context(), req)
	if err != nil {
		// Check for specific error types and return appropriate status codes
		if err.Error() == "builder profile with ID does not exist" {
			response.WriteError(w, http.StatusNotFound, "Builder profile not found")
			return
		}
		if err.Error() == "jobsite with ID does not exist" {
			response.WriteError(w, http.StatusBadRequest, "Jobsite not found")
			return
		}
		if err.Error() == "jobsite with ID does not belong to builder" {
			response.WriteError(w, http.StatusForbidden, "Jobsite does not belong to your builder profile")
			return
		}
		if err.Error() == "job type with ID does not exist" {
			response.WriteError(w, http.StatusBadRequest, "Job type not found")
			return
		}
		if err.Error() == "license with ID does not exist" {
			response.WriteError(w, http.StatusBadRequest, "One or more licenses not found")
			return
		}
		if err.Error() == "skill category with ID does not exist" {
			response.WriteError(w, http.StatusBadRequest, "One or more skill categories not found")
			return
		}
		if err.Error() == "skill subcategory with ID does not exist" {
			response.WriteError(w, http.StatusBadRequest, "One or more skill subcategories not found")
			return
		}
		if err.Error() == "validation failed" {
			response.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Generic error for other cases
		response.WriteError(w, http.StatusInternalServerError, "Failed to create job: "+err.Error())
		return
	}

	// Convert to response
	jobResp := h.convertToJobResponse(job)
	resp := payload.CreateJobResponse{
		Job:     jobResp,
		Message: "Job created successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

// GetJob retrieves a job by ID
func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.jobUsecase.GetJobByID(r.Context(), jobID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Job not found")
		return
	}

	jobResp := h.convertToJobResponse(job)
	resp := payload.GetJobResponse{
		Job:     jobResp,
		Message: "Job retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetJobWithRelations retrieves a job with all relations
func (h *JobHandler) GetJobWithRelations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	job, err := h.jobUsecase.GetJobWithRelations(r.Context(), jobID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Job not found")
		return
	}

	jobResp := h.convertToJobResponse(job)
	resp := payload.GetJobResponse{
		Job:     jobResp,
		Message: "Job with relations retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetJobsByBuilder retrieves jobs by builder profile ID
func (h *JobHandler) GetJobsByBuilder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	builderProfileID, err := uuid.Parse(vars["builder_profile_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid builder profile ID")
		return
	}

	// Get visibility filter from query params
	var visibility *models.JobVisibility
	if vis := r.URL.Query().Get("visibility"); vis != "" {
		v := models.JobVisibility(vis)
		visibility = &v
	}

	jobs, err := h.jobUsecase.GetJobsByBuilder(r.Context(), builderProfileID, visibility)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get jobs by builder")
		return
	}

	// Convert to response
	var jobsResp []payload.JobResponse
	for _, job := range jobs {
		jobsResp = append(jobsResp, h.convertToJobResponse(job))
	}

	resp := payload.GetJobsResponse{
		Jobs:    jobsResp,
		Message: "Jobs retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetJobsByJobsite retrieves jobs by jobsite ID
func (h *JobHandler) GetJobsByJobsite(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobsiteID, err := uuid.Parse(vars["jobsite_id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid jobsite ID")
		return
	}

	// Get visibility filter from query params
	var visibility *models.JobVisibility
	if vis := r.URL.Query().Get("visibility"); vis != "" {
		v := models.JobVisibility(vis)
		visibility = &v
	}

	jobs, err := h.jobUsecase.GetJobsByJobsite(r.Context(), jobsiteID, visibility)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get jobs by jobsite")
		return
	}

	// Convert to response
	var jobsResp []payload.JobResponse
	for _, job := range jobs {
		jobsResp = append(jobsResp, h.convertToJobResponse(job))
	}

	resp := payload.GetJobsResponse{
		Jobs:    jobsResp,
		Message: "Jobs retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetJobsByVisibility retrieves jobs by visibility
func (h *JobHandler) GetJobsByVisibility(w http.ResponseWriter, r *http.Request) {
	visibility := models.JobVisibility(r.URL.Query().Get("visibility"))
	if visibility == "" {
		response.WriteError(w, http.StatusBadRequest, "Visibility parameter is required")
		return
	}

	jobs, err := h.jobUsecase.GetJobsByVisibility(r.Context(), visibility)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get jobs by visibility")
		return
	}

	// Convert to response
	var jobsResp []payload.JobResponse
	for _, job := range jobs {
		jobsResp = append(jobsResp, h.convertToJobResponse(job))
	}

	resp := payload.GetJobsResponse{
		Jobs:    jobsResp,
		Message: "Jobs retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetAllJobs retrieves all jobs
func (h *JobHandler) GetAllJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.jobUsecase.GetAllJobs(r.Context())
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get all jobs")
		return
	}

	// Convert to response
	var jobsResp []payload.JobResponse
	for _, job := range jobs {
		jobsResp = append(jobsResp, h.convertToJobResponse(job))
	}

	resp := payload.GetJobsResponse{
		Jobs:    jobsResp,
		Message: "Jobs retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// UpdateJob updates a job
func (h *JobHandler) UpdateJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	var req payload.UpdateJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update job
	job, err := h.jobUsecase.UpdateJob(r.Context(), jobID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to update job")
		return
	}

	jobResp := h.convertToJobResponse(job)
	resp := payload.UpdateJobResponse{
		Job:     jobResp,
		Message: "Job updated successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// DeleteJob deletes a job
func (h *JobHandler) DeleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	if err := h.jobUsecase.DeleteJob(r.Context(), jobID); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to delete job")
		return
	}

	resp := payload.DeleteJobResponse{
		Message: "Job deleted successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// Helper function to convert Job model to JobResponse
func (h *JobHandler) convertToJobResponse(job *models.Job) payload.JobResponse {
	jobResp := payload.JobResponse{
		ID:                          job.ID,
		BuilderProfileID:            job.BuilderProfileID,
		JobsiteID:                   job.JobsiteID,
		JobTypeID:                   job.JobTypeID,
		ManyLabours:                 job.ManyLabours,
		OngoingWork:                 job.OngoingWork,
		WageSiteAllowance:           job.WageSiteAllowance,
		WageLeadingHandAllowance:    job.WageLeadingHandAllowance,
		WageProductivityAllowance:   job.WageProductivityAllowance,
		ExtrasOvertimeRate:          job.ExtrasOvertimeRate,
		StartDateWork:               job.StartDateWork,
		EndDateWork:                 job.EndDateWork,
		WorkSaturday:                job.WorkSaturday,
		WorkSunday:                  job.WorkSunday,
		StartTime:                   job.StartTime,
		EndTime:                     job.EndTime,
		Description:                 job.Description,
		PaymentDay:                  job.PaymentDay,
		RequiresSupervisorSignature: job.RequiresSupervisorSignature,
		SupervisorName:              job.SupervisorName,
		Visibility:                  job.Visibility,
		PaymentType:                 job.PaymentType,
		CreatedAt:                   job.CreatedAt,
		UpdatedAt:                   job.UpdatedAt,
	}

	// Convert relations if they exist (commented out due to circular imports)
	// Relations will be loaded separately using Preload() in repository
	// if job.BuilderProfile != nil {
	// 	jobResp.BuilderProfile = &payload.BuilderProfileResponse{
	// 		ID:          job.BuilderProfile.ID,
	// 		CompanyName: job.BuilderProfile.CompanyName,
	// 	}
	// }

	// if job.Jobsite != nil {
	// 	jobResp.Jobsite = &payload.JobsiteResponse{
	// 		ID:      job.Jobsite.ID,
	// 		Name:    job.Jobsite.Name,
	// 		Address: job.Jobsite.Address,
	// 	}
	// }

	// if job.JobType != nil {
	// 	jobResp.JobType = &payload.JobTypeResponse{
	// 		ID:          job.JobType.ID,
	// 		Name:        job.JobType.Name,
	// 		Description: job.JobType.Description,
	// 	}
	// }

	// Convert job licenses (commented out due to circular imports)
	// Relations will be loaded separately using Preload() in repository
	// for _, jobLicense := range job.JobLicenses {
	// 	jobLicenseResp := payload.JobLicenseResponse{
	// 		ID:        jobLicense.ID,
	// 		JobID:     jobLicense.JobID,
	// 		LicenseID: jobLicense.LicenseID,
	// 		CreatedAt: jobLicense.CreatedAt,
	// 	}
	// 	if jobLicense.License != nil {
	// 		jobLicenseResp.License = &payload.LicenseResponse{
	// 			ID:          jobLicense.License.ID,
	// 			Name:        jobLicense.License.Name,
	// 			Description: jobLicense.License.Description,
	// 		}
	// 	}
	// 	jobResp.JobLicenses = append(jobResp.JobLicenses, jobLicenseResp)
	// }

	// Convert job skills (commented out due to circular imports)
	// Relations will be loaded separately using Preload() in repository
	// for _, jobSkill := range job.JobSkills {
	// 	jobSkillResp := payload.JobSkillResponse{
	// 		ID:                 jobSkill.ID,
	// 		JobID:              jobSkill.JobID,
	// 		SkillCategoryID:    jobSkill.SkillCategoryID,
	// 		SkillSubcategoryID: jobSkill.SkillSubcategoryID,
	// 		CreatedAt:          jobSkill.CreatedAt,
	// 	}
	// 	if jobSkill.SkillCategory != nil {
	// 		jobSkillResp.SkillCategory = &payload.SkillCategoryResponse{
	// 			ID:          jobSkill.SkillCategory.ID,
	// 			Name:        jobSkill.SkillCategory.Name,
	// 			Description: jobSkill.SkillCategory.Description,
	// 		}
	// 	}
	// 	if jobSkill.SkillSubcategory != nil {
	// 		jobSkillResp.SkillSubcategory = &payload.SkillSubcategoryResponse{
	// 			ID:          jobSkill.SkillSubcategory.ID,
	// 			Name:        jobSkill.SkillSubcategory.Name,
	// 			Description: jobSkill.SkillSubcategory.Description,
	// 		}
	// 	}
	// 	jobResp.JobSkills = append(jobResp.JobSkills, jobSkillResp)
	// }

	return jobResp
}
