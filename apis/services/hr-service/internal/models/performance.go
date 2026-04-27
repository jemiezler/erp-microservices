package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// PerformanceGoal defines goals for employees
type PerformanceGoal struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	GoalTitle   string `json:"goal_title"`
	Description string `json:"description" gorm:"type:text"`

	Category string `json:"category"` // Sales, Operations, Quality, Innovation, Learning
	Type     string `json:"type"`     // Individual, Team, Department

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	TargetValue  float64 `json:"target_value"`
	CurrentValue float64 `json:"current_value"`
	Unit         string  `json:"unit"` // Percentage, Count, Amount, Rating

	Priority string  `json:"priority"` // Low, Medium, High, Critical
	Weight   float64 `json:"weight"`   // % weight in overall performance (all should sum to 100)

	Status string `json:"status"` // Draft, Active, Completed, Cancelled

	Alignments datatypes.JSONSlice[int] `json:"alignments" gorm:"type:jsonb"` // Linked to organizational goals

	Progress []GoalProgress `json:"progress,omitempty"`

	gorm.Model
}

// GoalProgress tracks goal achievement progress
type GoalProgress struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	GoalID uint             `json:"goal_id" gorm:"index"`
	Goal   *PerformanceGoal `json:"goal,omitempty"`

	CheckinDate     time.Time `json:"checkin_date"`
	ProgressValue   float64   `json:"progress_value"`
	ProgressPercent float64   `json:"progress_percent"` // Calculated

	Comments       string     `json:"comments"`
	LastReviewedBy uint       `json:"last_reviewed_by"`
	LastReviewDate *time.Time `json:"last_review_date"`

	gorm.Model
}

// PerformanceRating annual performance ratings
type PerformanceRating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	RatingNumber string    `json:"rating_number" gorm:"uniqueIndex"` // PR-2024-001
	EmployeeID   uint      `json:"employee_id" gorm:"index"`
	Employee     *Employee `json:"employee,omitempty"`

	RatePeriodStart time.Time `json:"rate_period_start"`
	RatePeriodEnd   time.Time `json:"rate_period_end"`

	ManagerID uint      `json:"manager_id" gorm:"index"`
	Manager   *Employee `json:"manager,omitempty" gorm:"foreignKey:ManagerID"`

	OverallRating float64 `json:"overall_rating"` // 1-5 or 1-10

	// Competency-based ratings
	Competencies []CompetencyRating `json:"competencies,omitempty"`

	// KPI achievements
	KPIScore float64 `json:"kpi_score"`

	// Goal achievement
	GoalAchievement float64 `json:"goal_achievement"`

	Strengths       string `json:"strengths" gorm:"type:text"`
	Improvements    string `json:"improvements" gorm:"type:text"`
	DevelopmentPlan string `json:"development_plan" gorm:"type:text"`

	Status string `json:"status"` // Draft, Submitted, Under Review, Finalized, Approved

	SubmittedDate *time.Time `json:"submitted_date"`
	FinalizedDate *time.Time `json:"finalized_date"`

	EmployeeComments string     `json:"employee_comments" gorm:"type:text"`
	EmployeeAckDate  *time.Time `json:"employee_ack_date"`

	GradeCode string `json:"grade_code"` // A, B, C, D, E

	gorm.Model
}

// CompetencyRating rates employee against competencies
type CompetencyRating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	PerformanceRatingID uint               `json:"performance_rating_id" gorm:"index"`
	PerformanceRating   *PerformanceRating `json:"performance_rating,omitempty"`

	CompetencyID uint        `json:"competency_id"`
	Competency   *Competency `json:"competency,omitempty"`

	Rating   float64 `json:"rating"` // 1-5
	Comments string  `json:"comments"`

	gorm.Model
}

// Competency defines skills/competencies
type Competency struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	Name        string `json:"name"` // Leadership, Communication, Technical Skills
	Description string `json:"description"`
	Category    string `json:"category"` // Core, Technical, Leadership, Behavioral

	Levels datatypes.JSONSlice[string] `json:"levels" gorm:"type:jsonb"` // [Beginner, Intermediate, Advanced, Expert]

	IsActive bool `json:"is_active" gorm:"default:true"`

	gorm.Model
}

// EmployeeReview 360-degree feedback reviews
type EmployeeReview struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	ReviewNumber string    `json:"review_number" gorm:"uniqueIndex"`
	EmployeeID   uint      `json:"employee_id" gorm:"index"`
	Employee     *Employee `json:"employee,omitempty"`

	ReviewType string `json:"review_type"` // 360, Manager, Peer, Subordinate, Self

	ReviewerID uint      `json:"reviewer_id" gorm:"index"`
	Reviewer   *Employee `json:"reviewer,omitempty" gorm:"foreignKey:ReviewerID"`

	ReviewPeriod string `json:"review_period"`

	RatingScale   string  `json:"rating_scale"` // 1-5, 1-10
	OverallRating float64 `json:"overall_rating"`

	Questions []ReviewQuestion `json:"questions,omitempty"`

	Feedback string `json:"feedback" gorm:"type:text"`

	Status         string     `json:"status"` // Draft, Completed, Submitted, Anonymous
	CompletionDate *time.Time `json:"completion_date"`

	IsAnonymous bool `json:"is_anonymous"`

	gorm.Model
}

// ReviewQuestion individual review questions and answers
type ReviewQuestion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ReviewID uint            `json:"review_id" gorm:"index"`
	Review   *EmployeeReview `json:"review,omitempty"`

	QuestionNumber int    `json:"question_number"`
	Question       string `json:"question"`
	QuestionType   string `json:"question_type"` // Rating, Text, MultipleChoice

	Rating     *float64 `json:"rating"`
	TextAnswer string   `json:"text_answer"`

	gorm.Model
}

// TrainingRequest training and development requests
type TrainingRequest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	TrainingID string    `json:"training_id" gorm:"uniqueIndex"`
	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	TrainingTitle string `json:"training_title"`
	Provider      string `json:"provider"`
	Category      string `json:"category"` // Technical, Leadership, Soft Skills

	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`

	Cost     float64 `json:"cost"`
	Currency string  `json:"currency"`

	Justification string `json:"justification" gorm:"type:text"`

	Status string `json:"status"` // Pending, Approved, Rejected, Completed

	ApprovedByID *uint      `json:"approved_by_id"`
	ApprovedDate *time.Time `json:"approved_date"`

	Feedback string   `json:"feedback" gorm:"type:text"`
	Rating   *float64 `json:"rating"`

	Certificate string `json:"certificate"` // URL to certificate

	gorm.Model
}
