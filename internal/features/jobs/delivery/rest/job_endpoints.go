package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	builder_db "github.com/yakka-backend/internal/features/builder_profiles/entity/database"
	builder_models "github.com/yakka-backend/internal/features/builder_profiles/models"
	"github.com/yakka-backend/internal/features/jobs/models"
	"github.com/yakka-backend/internal/features/jobs/payload"
	"github.com/yakka-backend/internal/features/jobs/usecase"
	jobsite_db "github.com/yakka-backend/internal/features/jobsites/entity/database"
	jobsite_models "github.com/yakka-backend/internal/features/jobsites/models"
	job_requirement_db "github.com/yakka-backend/internal/features/masters/job_requirements/entity/database"
	job_type_db "github.com/yakka-backend/internal/features/masters/job_types/entity/database"
	job_type_models "github.com/yakka-backend/internal/features/masters/job_types/models"
	license_db "github.com/yakka-backend/internal/features/masters/licenses/entity/database"
	skill_category_db "github.com/yakka-backend/internal/features/masters/skills/entity/database"
	"github.com/yakka-backend/internal/infrastructure/http/middleware"
	"github.com/yakka-backend/internal/shared/response"
	"github.com/yakka-backend/internal/shared/validation"
)

type JobHandler struct {
	jobUsecase           usecase.JobUsecase
	builderProfileRepo   builder_db.BuilderProfileRepository
	jobsiteRepo          jobsite_db.JobsiteRepository
	jobTypeRepo          job_type_db.JobTypeRepository
	licenseRepo          license_db.LicenseRepository
	jobRequirementRepo   job_requirement_db.JobRequirementRepository
	skillCategoryRepo    skill_category_db.SkillCategoryRepository
	skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository
}

func NewJobHandler(jobUsecase usecase.JobUsecase, builderProfileRepo builder_db.BuilderProfileRepository, jobsiteRepo jobsite_db.JobsiteRepository, jobTypeRepo job_type_db.JobTypeRepository, licenseRepo license_db.LicenseRepository, jobRequirementRepo job_requirement_db.JobRequirementRepository, skillCategoryRepo skill_category_db.SkillCategoryRepository, skillSubcategoryRepo skill_category_db.SkillSubcategoryRepository) *JobHandler {
	return &JobHandler{
		jobUsecase:           jobUsecase,
		builderProfileRepo:   builderProfileRepo,
		jobsiteRepo:          jobsiteRepo,
		jobTypeRepo:          jobTypeRepo,
		licenseRepo:          licenseRepo,
		jobRequirementRepo:   jobRequirementRepo,
		skillCategoryRepo:    skillCategoryRepo,
		skillSubcategoryRepo: skillSubcategoryRepo,
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

	// Get job with all relations for response
	log.Printf("🔍 CreateJob - Getting job with relations for ID: %s", job.ID)
	jobWithRelations, err := h.jobUsecase.GetJobWithRelations(r.Context(), job.ID)
	if err != nil {
		log.Printf("🚫 CreateJob - Failed to get job with relations: %v", err)
		response.WriteError(w, http.StatusInternalServerError, "Failed to get job details")
		return
	}
	log.Printf("🔍 CreateJob - Job with relations loaded, JobSkills count: %d", len(jobWithRelations.JobSkills))

	// Get additional relation data
	jobsite, err := h.jobsiteRepo.GetByID(r.Context(), job.JobsiteID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get jobsite details")
		return
	}

	jobType, err := h.jobTypeRepo.GetByID(r.Context(), job.JobTypeID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get job type details")
		return
	}

	// Convert to response with relations using usecase method
	jobResp := h.convertToJobResponseWithRelations(r.Context(), jobWithRelations, builderProfile, jobsite, jobType)
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

// GetMyJobs retrieves jobs for the authenticated builder
func (h *JobHandler) GetMyJobs(w http.ResponseWriter, r *http.Request) {
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

	// Get visibility filter from query params
	var visibility *models.JobVisibility
	if vis := r.URL.Query().Get("visibility"); vis != "" {
		v := models.JobVisibility(vis)
		visibility = &v
	}

	jobs, err := h.jobUsecase.GetJobsByBuilder(r.Context(), builderProfile.ID, visibility)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get jobs")
		return
	}

	// Convert to response with full relations
	var jobsResp []payload.JobResponse
	for _, job := range jobs {
		// Get job with all relations for each job
		jobWithRelations, err := h.jobUsecase.GetJobWithRelations(r.Context(), job.ID)
		if err != nil {
			log.Printf("🚫 GetMyJobs - Failed to get job with relations for ID %s: %v", job.ID, err)
			// Fallback to basic response if relations fail
			jobsResp = append(jobsResp, h.convertToJobResponse(job))
			continue
		}

		// Get additional relation data
		jobsite, err := h.jobsiteRepo.GetByID(r.Context(), job.JobsiteID)
		if err != nil {
			log.Printf("🚫 GetMyJobs - Failed to get jobsite for job %s: %v", job.ID, err)
			jobsite = nil
		}

		jobType, err := h.jobTypeRepo.GetByID(r.Context(), job.JobTypeID)
		if err != nil {
			log.Printf("🚫 GetMyJobs - Failed to get job type for job %s: %v", job.ID, err)
			jobType = nil
		}

		// Convert to response with full relations
		jobResp := h.convertToJobResponseWithRelations(r.Context(), jobWithRelations, builderProfile, jobsite, jobType)
		jobsResp = append(jobsResp, jobResp)
	}

	resp := payload.GetJobsResponse{
		Jobs:    jobsResp,
		Message: "Jobs retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// GetBuilderApplicants retrieves all applicants for the authenticated builder's jobs
func (h *JobHandler) GetBuilderApplicants(w http.ResponseWriter, r *http.Request) {
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

	// Get applicants for builder's jobs grouped by jobsite
	jobsitesWithJobs, err := h.jobUsecase.GetBuilderApplicantsByJobsite(r.Context(), builderProfile.ID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get applicants")
		return
	}

	// Calculate total jobs
	totalJobs := 0
	for _, jobsite := range jobsitesWithJobs {
		totalJobs += len(jobsite.Jobs)
	}

	resp := payload.BuilderApplicantsByJobsiteResponse{
		Jobsites: jobsitesWithJobs,
		Total:    totalJobs,
		Message:  "Applicants retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// ProcessApplicantDecision processes hiring or rejection of an applicant
func (h *JobHandler) ProcessApplicantDecision(w http.ResponseWriter, r *http.Request) {
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

	var req payload.BuilderApplicantDecisionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Process the decision
	result, err := h.jobUsecase.ProcessApplicantDecision(r.Context(), builderProfile.ID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to process applicant decision")
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
}

// GetLabourJobs retrieves all jobs with application status for a labour user
func (h *JobHandler) GetLabourJobs(w http.ResponseWriter, r *http.Request) {
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

	// Get jobs with application status for this labour user
	jobs, err := h.jobUsecase.GetLabourJobs(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get jobs")
		return
	}

	resp := payload.LabourJobsResponse{
		Jobs:    jobs,
		Total:   len(jobs),
		Message: "Jobs retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// ApplyToJob allows a labour user to apply for a job
func (h *JobHandler) ApplyToJob(w http.ResponseWriter, r *http.Request) {
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

	var req payload.LabourApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(req); err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Apply to job
	result, err := h.jobUsecase.ApplyToJob(r.Context(), userID, req)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to apply to job")
		return
	}

	response.WriteJSON(w, http.StatusCreated, result)
}

// GetLabourApplicants retrieves all applications for the authenticated labour user
func (h *JobHandler) GetLabourApplicants(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by LabourMiddleware)
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

	// Get applications for this labour user
	result, err := h.jobUsecase.GetLabourApplicants(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "Failed to get applications")
		return
	}

	response.WriteJSON(w, http.StatusOK, result)
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
		ManyLabours:                 job.ManyLabours,
		OngoingWork:                 job.OngoingWork,
		WageSiteAllowance:           job.WageSiteAllowance,
		WageLeadingHandAllowance:    job.WageLeadingHandAllowance,
		WageProductivityAllowance:   job.WageProductivityAllowance,
		ExtrasOvertimeRate:          job.ExtrasOvertimeRate,
		WageHourlyRate:              job.WageHourlyRate,
		TravelAllowance:             job.TravelAllowance,
		GST:                         job.GST,
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

// GetLabourJobDetail retrieves a job detail for a labour user with application info
func (h *JobHandler) GetLabourJobDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	// Get labour user ID from context (set by AuthMiddleware)
	log.Printf("🔍 Labour Handler - Checking context for user_id")
	labourUserIDStr, ok := r.Context().Value(middleware.UserIDKey).(string)
	if !ok {
		log.Printf("🚫 Labour Handler - User ID not found in context")
		response.WriteError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}
	log.Printf("🔍 Labour Handler - User ID found: %s", labourUserIDStr)

	// Parse string to UUID
	labourUserID, err := uuid.Parse(labourUserIDStr)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Invalid user ID format")
		return
	}

	// Get job detail with application info
	jobDetail, err := h.jobUsecase.GetLabourJobDetail(r.Context(), jobID, labourUserID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Job not found")
		return
	}

	response.WriteJSON(w, http.StatusOK, jobDetail)
}

// GetBuilderJobDetail retrieves a job detail for a builder (only their own jobs)
func (h *JobHandler) GetBuilderJobDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	// Get builder profile ID from context (set by BuilderMiddleware)
	log.Printf("🔍 Handler - Checking context for builder_profile_id")
	builderProfileIDStr, ok := r.Context().Value(middleware.BuilderProfileIDKey).(string)
	if !ok {
		log.Printf("🚫 Handler - Builder profile ID not found in context")
		response.WriteError(w, http.StatusUnauthorized, "Builder profile ID not found in context")
		return
	}
	log.Printf("🔍 Handler - Builder profile ID found: %s", builderProfileIDStr)

	// Parse string to UUID
	builderProfileID, err := uuid.Parse(builderProfileIDStr)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Invalid builder profile ID format")
		return
	}

	// Get job with all relations
	job, err := h.jobUsecase.GetJobWithRelations(r.Context(), jobID)
	if err != nil {
		response.WriteError(w, http.StatusNotFound, "Job not found")
		return
	}

	// Verify that the job belongs to the builder
	if job.BuilderProfileID != builderProfileID {
		response.WriteError(w, http.StatusForbidden, "Invalid job - not owned by builder")
		return
	}

	// Get additional relation data
	builderProfile, err := h.builderProfileRepo.GetByID(r.Context(), job.BuilderProfileID)
	if err != nil {
		log.Printf("🚫 GetBuilderJobDetail - Failed to get builder profile: %v", err)
	}

	jobsite, err := h.jobsiteRepo.GetByID(r.Context(), job.JobsiteID)
	if err != nil {
		log.Printf("🚫 GetBuilderJobDetail - Failed to get jobsite: %v", err)
	}

	jobType, err := h.jobTypeRepo.GetByID(r.Context(), job.JobTypeID)
	if err != nil {
		log.Printf("🚫 GetBuilderJobDetail - Failed to get job type: %v", err)
	}

	// Convert to detail response (shows null values)
	jobResp := h.convertToJobDetailResponse(r.Context(), job, builderProfile, jobsite, jobType)

	resp := payload.GetBuilderJobDetailResponse{
		Job:     jobResp,
		Message: "Job detail retrieved successfully",
	}

	response.WriteJSON(w, http.StatusOK, resp)
}

// UpdateJobVisibility updates the visibility of a job (only for the owner builder)
func (h *JobHandler) UpdateJobVisibility(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID, err := uuid.Parse(vars["id"])
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid job ID")
		return
	}

	log.Printf("🔍 UpdateJobVisibility - Job ID: %s", jobID)

	// Get builder profile ID from context (set by BuilderMiddleware)
	builderProfileIDStr, ok := r.Context().Value(middleware.BuilderProfileIDKey).(string)
	if !ok {
		log.Printf("🚫 UpdateJobVisibility - Builder profile ID not found in context")
		response.WriteError(w, http.StatusUnauthorized, "Builder profile ID not found in context")
		return
	}
	log.Printf("🔍 UpdateJobVisibility - Builder Profile ID: %s", builderProfileIDStr)

	// Parse string to UUID
	builderProfileID, err := uuid.Parse(builderProfileIDStr)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, "Invalid builder profile ID format")
		return
	}

	// Parse request body
	var req payload.UpdateJobVisibilityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update job visibility (only if owned by builder)
	jobDetail, err := h.jobUsecase.UpdateJobVisibility(r.Context(), jobID, builderProfileID, req)
	if err != nil {
		if err.Error() == "job not found" {
			response.WriteError(w, http.StatusNotFound, "Job not found")
			return
		}
		if err.Error() == "invalid job - not owned by builder" {
			response.WriteError(w, http.StatusForbidden, "Invalid job - not owned by builder")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "Failed to update job visibility")
		return
	}

	response.WriteJSON(w, http.StatusOK, jobDetail)
}

// convertToJobResponseWithRelations converts a Job model to JobResponse with full relations
func (h *JobHandler) convertToJobResponseWithRelations(ctx context.Context, job *models.Job, builderProfile *builder_models.BuilderProfile, jobsite *jobsite_models.Jobsite, jobType *job_type_models.JobType) payload.JobResponse {
	jobResp := payload.JobResponse{
		ID:                          job.ID,
		ManyLabours:                 job.ManyLabours,
		OngoingWork:                 job.OngoingWork,
		WageSiteAllowance:           job.WageSiteAllowance,
		WageLeadingHandAllowance:    job.WageLeadingHandAllowance,
		WageProductivityAllowance:   job.WageProductivityAllowance,
		ExtrasOvertimeRate:          job.ExtrasOvertimeRate,
		WageHourlyRate:              job.WageHourlyRate,
		TravelAllowance:             job.TravelAllowance,
		GST:                         job.GST,
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

	// Calculate total wage
	totalWage := calculateTotalWage(job)
	jobResp.TotalWage = &totalWage

	// Add builder profile information
	if builderProfile != nil {
		jobResp.BuilderProfile = &payload.BuilderProfileResponse{
			ID:          builderProfile.ID,
			CompanyName: getCompanyName(builderProfile.Company),
			DisplayName: builderProfile.DisplayName,
			Location:    builderProfile.Location,
			Phone:       nil,
			Email:       nil,
			CreatedAt:   builderProfile.CreatedAt,
			UpdatedAt:   builderProfile.UpdatedAt,
		}
	}

	// Add jobsite information
	if jobsite != nil {
		jobResp.Jobsite = &payload.JobsiteResponse{
			ID:          jobsite.ID,
			Name:        getStringValue(jobsite.Description),
			Address:     jobsite.Address,
			City:        jobsite.City,
			Suburb:      jobsite.Suburb,
			Description: jobsite.Description,
			Latitude:    jobsite.Latitude,
			Longitude:   jobsite.Longitude,
			Phone:       jobsite.Phone,
			CreatedAt:   jobsite.CreatedAt,
			UpdatedAt:   jobsite.UpdatedAt,
		}
	}

	// Add job type information
	if jobType != nil {
		jobResp.JobType = &payload.JobTypeResponse{
			ID:          jobType.ID,
			Name:        jobType.Name,
			Description: jobType.Description,
			CreatedAt:   jobType.CreatedAt,
			UpdatedAt:   jobType.UpdatedAt,
		}
	}

	// Convert job licenses with full details
	for _, jobLicense := range job.JobLicenses {
		// Get license details
		licenseDetails, err := h.licenseRepo.GetByID(ctx, jobLicense.LicenseID)
		if err != nil {
			log.Printf("🚫 Failed to get license details for ID %s: %v", jobLicense.LicenseID, err)
			continue
		}

		jobLicenseResp := payload.JobLicenseResponse{
			ID:        jobLicense.ID,
			JobID:     jobLicense.JobID,
			LicenseID: jobLicense.LicenseID,
			License: &payload.LicenseResponse{
				ID:          licenseDetails.ID,
				Name:        licenseDetails.Name,
				Description: &licenseDetails.Description,
			},
			CreatedAt: jobLicense.CreatedAt,
		}
		jobResp.JobLicenses = append(jobResp.JobLicenses, jobLicenseResp)
	}

	// Convert job skills with full details
	log.Printf("🔍 convertToJobResponseWithRelations - JobSkills count: %d", len(job.JobSkills))
	for i, jobSkill := range job.JobSkills {
		log.Printf("🔍 convertToJobResponseWithRelations - JobSkill %d: ID=%s, CategoryID=%v, SubcategoryID=%v", i, jobSkill.ID, jobSkill.SkillCategoryID, jobSkill.SkillSubcategoryID)

		jobSkillResp := payload.JobSkillResponse{
			ID:                 jobSkill.ID,
			JobID:              jobSkill.JobID,
			SkillCategoryID:    jobSkill.SkillCategoryID,
			SkillSubcategoryID: jobSkill.SkillSubcategoryID,
			CreatedAt:          jobSkill.CreatedAt,
		}

		// Get skill category details if available
		if jobSkill.SkillCategoryID != nil {
			skillCategory, err := h.skillCategoryRepo.GetByID(ctx, *jobSkill.SkillCategoryID)
			if err != nil {
				log.Printf("🚫 Failed to get skill category details for ID %s: %v", *jobSkill.SkillCategoryID, err)
			} else {
				jobSkillResp.SkillCategory = &payload.SkillCategoryResponse{
					ID:          skillCategory.ID,
					Name:        skillCategory.Name,
					Description: &skillCategory.Description,
				}
				log.Printf("🔍 convertToJobResponseWithRelations - SkillCategory loaded: %s", skillCategory.Name)
			}
		}

		// Get skill subcategory details if available
		if jobSkill.SkillSubcategoryID != nil {
			skillSubcategory, err := h.skillSubcategoryRepo.GetByID(ctx, *jobSkill.SkillSubcategoryID)
			if err != nil {
				log.Printf("🚫 Failed to get skill subcategory details for ID %s: %v", *jobSkill.SkillSubcategoryID, err)
			} else {
				jobSkillResp.SkillSubcategory = &payload.SkillSubcategoryResponse{
					ID:          skillSubcategory.ID,
					Name:        skillSubcategory.Name,
					Description: &skillSubcategory.Description,
				}
				log.Printf("🔍 convertToJobResponseWithRelations - SkillSubcategory loaded: %s", skillSubcategory.Name)
			}
		}

		jobResp.JobSkills = append(jobResp.JobSkills, jobSkillResp)
	}
	log.Printf("🔍 convertToJobResponseWithRelations - Final JobSkills count: %d", len(jobResp.JobSkills))

	// Convert job requirements with full details
	for _, jobRequirement := range job.JobRequirements {
		// Get job requirement details
		requirementDetails, err := h.jobRequirementRepo.GetByID(ctx, jobRequirement.JobRequirementID)
		if err != nil {
			log.Printf("🚫 Failed to get job requirement details for ID %s: %v", jobRequirement.JobRequirementID, err)
			continue
		}

		jobRequirementResp := payload.JobRequirementResponse{
			ID:          requirementDetails.ID,
			Name:        requirementDetails.Name,
			Description: requirementDetails.Description,
			IsActive:    requirementDetails.IsActive,
			CreatedAt:   requirementDetails.CreatedAt,
			UpdatedAt:   requirementDetails.UpdatedAt,
		}
		jobResp.JobRequirements = append(jobResp.JobRequirements, jobRequirementResp)
	}

	return jobResp
}

// getStringValue safely dereferences a string pointer
func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// getCompanyName safely gets company name from Company model
func getCompanyName(company *builder_models.Company) string {
	if company == nil {
		return ""
	}
	return company.Name
}

// calculateTotalWage calculates the total wage by summing all wage components
func calculateTotalWage(job *models.Job) float64 {
	var total float64

	if job.WageSiteAllowance != nil {
		total += *job.WageSiteAllowance
	}
	if job.WageLeadingHandAllowance != nil {
		total += *job.WageLeadingHandAllowance
	}
	if job.WageProductivityAllowance != nil {
		total += *job.WageProductivityAllowance
	}
	if job.WageHourlyRate != nil {
		total += *job.WageHourlyRate
	}
	if job.TravelAllowance != nil {
		total += *job.TravelAllowance
	}
	if job.ExtrasOvertimeRate != nil {
		total += *job.ExtrasOvertimeRate
	}

	return total
}

// convertToJobDetailResponse converts a Job model to JobDetailResponse with full relations (shows null values)
func (h *JobHandler) convertToJobDetailResponse(ctx context.Context, job *models.Job, builderProfile *builder_models.BuilderProfile, jobsite *jobsite_models.Jobsite, jobType *job_type_models.JobType) payload.JobDetailResponse {
	jobResp := payload.JobDetailResponse{
		ID:                          job.ID,
		ManyLabours:                 job.ManyLabours,
		OngoingWork:                 job.OngoingWork,
		WageSiteAllowance:           job.WageSiteAllowance,
		WageLeadingHandAllowance:    job.WageLeadingHandAllowance,
		WageProductivityAllowance:   job.WageProductivityAllowance,
		ExtrasOvertimeRate:          job.ExtrasOvertimeRate,
		WageHourlyRate:              job.WageHourlyRate,
		TravelAllowance:             job.TravelAllowance,
		GST:                         job.GST,
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

	// Calculate total wage
	totalWage := calculateTotalWage(job)
	jobResp.TotalWage = &totalWage

	// Add builder profile information
	if builderProfile != nil {
		jobResp.BuilderProfile = &payload.BuilderProfileResponse{
			ID:          builderProfile.ID,
			CompanyName: getCompanyName(builderProfile.Company),
			DisplayName: builderProfile.DisplayName,
			Location:    builderProfile.Location,
			Phone:       nil,
			Email:       nil,
			CreatedAt:   builderProfile.CreatedAt,
			UpdatedAt:   builderProfile.UpdatedAt,
		}
	}

	// Add jobsite information
	if jobsite != nil {
		jobResp.Jobsite = &payload.JobsiteResponse{
			ID:          jobsite.ID,
			Name:        getStringValue(jobsite.Description),
			Address:     jobsite.Address,
			City:        jobsite.City,
			Suburb:      jobsite.Suburb,
			Description: jobsite.Description,
			Latitude:    jobsite.Latitude,
			Longitude:   jobsite.Longitude,
			Phone:       jobsite.Phone,
			CreatedAt:   jobsite.CreatedAt,
			UpdatedAt:   jobsite.UpdatedAt,
		}
	}

	// Add job type information
	if jobType != nil {
		jobResp.JobType = &payload.JobTypeResponse{
			ID:          jobType.ID,
			Name:        jobType.Name,
			Description: jobType.Description,
			CreatedAt:   jobType.CreatedAt,
			UpdatedAt:   jobType.UpdatedAt,
		}
	}

	// Convert job licenses with full details
	for _, jobLicense := range job.JobLicenses {
		// Get license details
		licenseDetails, err := h.licenseRepo.GetByID(ctx, jobLicense.LicenseID)
		if err != nil {
			log.Printf("🚫 Failed to get license details for ID %s: %v", jobLicense.LicenseID, err)
			continue
		}

		jobLicenseResp := payload.JobLicenseResponse{
			ID:        jobLicense.ID,
			JobID:     jobLicense.JobID,
			LicenseID: jobLicense.LicenseID,
			License: &payload.LicenseResponse{
				ID:          licenseDetails.ID,
				Name:        licenseDetails.Name,
				Description: &licenseDetails.Description,
			},
			CreatedAt: jobLicense.CreatedAt,
		}
		jobResp.JobLicenses = append(jobResp.JobLicenses, jobLicenseResp)
	}

	// Convert job skills with full details
	for _, jobSkill := range job.JobSkills {
		jobSkillResp := payload.JobSkillResponse{
			ID:                 jobSkill.ID,
			JobID:              jobSkill.JobID,
			SkillCategoryID:    jobSkill.SkillCategoryID,
			SkillSubcategoryID: jobSkill.SkillSubcategoryID,
			CreatedAt:          jobSkill.CreatedAt,
		}

		// Get skill category details if available
		if jobSkill.SkillCategoryID != nil {
			skillCategory, err := h.skillCategoryRepo.GetByID(ctx, *jobSkill.SkillCategoryID)
			if err != nil {
				log.Printf("🚫 Failed to get skill category details for ID %s: %v", *jobSkill.SkillCategoryID, err)
			} else {
				jobSkillResp.SkillCategory = &payload.SkillCategoryResponse{
					ID:          skillCategory.ID,
					Name:        skillCategory.Name,
					Description: &skillCategory.Description,
				}
			}
		}

		// Get skill subcategory details if available
		if jobSkill.SkillSubcategoryID != nil {
			skillSubcategory, err := h.skillSubcategoryRepo.GetByID(ctx, *jobSkill.SkillSubcategoryID)
			if err != nil {
				log.Printf("🚫 Failed to get skill subcategory details for ID %s: %v", *jobSkill.SkillSubcategoryID, err)
			} else {
				jobSkillResp.SkillSubcategory = &payload.SkillSubcategoryResponse{
					ID:          skillSubcategory.ID,
					Name:        skillSubcategory.Name,
					Description: &skillSubcategory.Description,
				}
			}
		}

		jobResp.JobSkills = append(jobResp.JobSkills, jobSkillResp)
	}

	// Convert job requirements with full details
	for _, jobRequirement := range job.JobRequirements {
		// Get job requirement details
		requirementDetails, err := h.jobRequirementRepo.GetByID(ctx, jobRequirement.JobRequirementID)
		if err != nil {
			log.Printf("🚫 Failed to get job requirement details for ID %s: %v", jobRequirement.JobRequirementID, err)
			continue
		}

		jobRequirementResp := payload.JobRequirementResponse{
			ID:          requirementDetails.ID,
			Name:        requirementDetails.Name,
			Description: requirementDetails.Description,
			IsActive:    requirementDetails.IsActive,
			CreatedAt:   requirementDetails.CreatedAt,
			UpdatedAt:   requirementDetails.UpdatedAt,
		}
		jobResp.JobRequirements = append(jobResp.JobRequirements, jobRequirementResp)
	}

	return jobResp
}
