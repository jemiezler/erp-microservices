package models

import (
	"time"

	"gorm.io/datatypes"
)

// JobPosting represents a job opening
type JobPosting struct {
	ID              uint                        `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time                   `json:"created_at"`
	UpdatedAt       time.Time                   `json:"updated_at"`
	TenantID        string                      `json:"tenant_id" gorm:"index"`
	DepartmentID    uint                        `json:"department_id" gorm:"index"`
	Department      *Department                 `json:"department,omitempty" gorm:"foreignKey:DepartmentID"`
	JobTitleID      uint                        `json:"job_title_id" gorm:"index"`
	JobTitle        *JobTitle                   `json:"job_title,omitempty" gorm:"foreignKey:JobTitleID"`
	Title           string                      `json:"title"`
	Description     string                      `json:"description" gorm:"type:text"`
	NoOfPositions   int                         `json:"no_of_positions"`
	SalaryMin       float64                     `json:"salary_min"`
	SalaryMax       float64                     `json:"salary_max"`
	Currency        string                      `json:"currency" gorm:"default:USD"`
	Location        string                      `json:"location"`
	EmploymentType  string                      `json:"employment_type"` // Full-time, Part-time, Contract
	ExperienceMin   int                         `json:"experience_min"`  // Years
	ExperienceMax   int                         `json:"experience_max"`  // Years
	RequiredSkills  datatypes.JSONSlice[string] `json:"required_skills" gorm:"type:jsonb"`
	PreferredSkills datatypes.JSONSlice[string] `json:"preferred_skills" gorm:"type:jsonb"`
	Qualifications  string                      `json:"qualifications" gorm:"type:text"`
	PostedDate      time.Time                   `json:"posted_date"`
	ClosingDate     time.Time                   `json:"closing_date"`
	Status          string                      `json:"status" gorm:"default:open"` // open, closed, draft
	CreatedBy       uint                        `json:"created_by"`
	ApprovedBy      *uint                       `json:"approved_by"`
	Candidates      []Candidate                 `json:"candidates,omitempty" gorm:"foreignKey:JobPostingID"`
	IsActive        bool                        `json:"is_active" gorm:"default:true"`
}

// Candidate represents job applicant
type Candidate struct {
	ID                 uint        `gorm:"primaryKey" json:"id"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	TenantID           string      `json:"tenant_id" gorm:"index"`
	JobPostingID       uint        `json:"job_posting_id" gorm:"index"`
	JobPosting         *JobPosting `json:"job_posting,omitempty" gorm:"foreignKey:JobPostingID"`
	FirstName          string      `json:"first_name"`
	LastName           string      `json:"last_name"`
	Email              string      `json:"email" gorm:"index"`
	PhoneNumber        string      `json:"phone_number"`
	CurrentCompany     string      `json:"current_company"`
	CurrentDesignation string      `json:"current_designation"`
	YearsOfExperience  int         `json:"years_of_experience"`
	ResumeURL          string      `json:"resume_url"`
	LinkedInURL        string      `json:"linkedin_url"`
	Source             string      `json:"source"`                        // LinkedIn, Indeed, Internal, Referral, Career Portal
	SkillsMatch        float64     `json:"skills_match"`                  // 0-100 percentage
	Status             string      `json:"status" gorm:"default:applied"` // applied, screening, interview, offer, rejected, hired
	Rating             float64     `json:"rating"`                        // 1-5
	Comments           string      `json:"comments" gorm:"type:text"`
	Interviews         []Interview `json:"interviews,omitempty" gorm:"foreignKey:CandidateID"`
	Offers             []JobOffer  `json:"offers,omitempty" gorm:"foreignKey:CandidateID"`
	IsActive           bool        `json:"is_active" gorm:"default:true"`
}

// Interview represents interview rounds
type Interview struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	TenantID      string     `json:"tenant_id" gorm:"index"`
	CandidateID   uint       `json:"candidate_id" gorm:"index"`
	Candidate     *Candidate `json:"candidate,omitempty" gorm:"foreignKey:CandidateID"`
	InterviewType string     `json:"interview_type"` // phone, video, in-person, technical, hr
	InterviewDate time.Time  `json:"interview_date"`
	InterviewerID uint       `json:"interviewer_id"`
	Interviewer   *Employee  `json:"interviewer,omitempty" gorm:"foreignKey:InterviewerID"`
	Location      string     `json:"location"`
	MeetingLink   string     `json:"meeting_link"`
	Rating        float64    `json:"rating"` // 1-5
	Comments      string     `json:"comments" gorm:"type:text"`
	Result        string     `json:"result"` // pass, fail, pending
	Feedback      string     `json:"feedback" gorm:"type:text"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`
}

// JobOffer represents job offer extended to candidate
type JobOffer struct {
	ID           uint        `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	TenantID     string      `json:"tenant_id" gorm:"index"`
	CandidateID  uint        `json:"candidate_id" gorm:"index"`
	Candidate    *Candidate  `json:"candidate,omitempty" gorm:"foreignKey:CandidateID"`
	JobPostingID uint        `json:"job_posting_id" gorm:"index"`
	JobPosting   *JobPosting `json:"job_posting,omitempty" gorm:"foreignKey:JobPostingID"`
	OfferNumber  string      `json:"offer_number" gorm:"uniqueIndex"`
	Position     string      `json:"position"`
	Department   string      `json:"department"`
	CtcOffered   float64     `json:"ctc_offered"`
	Currency     string      `json:"currency" gorm:"default:USD"`
	JoiningDate  time.Time   `json:"joining_date"`
	ValidTill    time.Time   `json:"valid_till"`
	Terms        string      `json:"terms" gorm:"type:text"`
	Status       string      `json:"status" gorm:"default:pending"` // pending, accepted, rejected, withdrawn
	AcceptedDate *time.Time  `json:"accepted_date"`
	CreatedBy    uint        `json:"created_by"`
	ApprovedBy   *uint       `json:"approved_by"`
	ApprovedDate *time.Time  `json:"approved_date"`
	IsActive     bool        `json:"is_active" gorm:"default:true"`
}

// OnboardingTask represents onboarding activities
type OnboardingTask struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	TenantID      string     `json:"tenant_id" gorm:"index"`
	EmployeeID    uint       `json:"employee_id" gorm:"index"`
	Employee      *Employee  `json:"employee,omitempty" gorm:"foreignKey:EmployeeID"`
	TaskCategory  string     `json:"task_category"` // IT, HR, Manager, Security, Facilities
	TaskName      string     `json:"task_name"`
	Description   string     `json:"description" gorm:"type:text"`
	AssignedTo    uint       `json:"assigned_to"`
	AssignedUser  *Employee  `json:"assigned_user,omitempty" gorm:"foreignKey:AssignedTo"`
	DueDate       time.Time  `json:"due_date"`
	CompletedDate *time.Time `json:"completed_date"`
	Status        string     `json:"status" gorm:"default:pending"` // pending, in_progress, completed, skipped
	IsActive      bool       `json:"is_active" gorm:"default:true"`
}
