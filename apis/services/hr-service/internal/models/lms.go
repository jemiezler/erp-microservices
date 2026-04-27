package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// LMSCourse learning courses
type LMSCourse struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	CourseID    string `json:"course_id" gorm:"uniqueIndex"`
	Name        string `json:"name"`
	Description string `json:"description" gorm:"type:text"`

	Category string `json:"category"` // Technical, Management, Compliance, Soft Skills
	Level    string `json:"level"`    // Beginner, Intermediate, Advanced

	CreatedByID uint      `json:"created_by_id" gorm:"index"`
	CreatedBy   *Employee `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`

	DurationHours float64 `json:"duration_hours"`

	ThumbnailURL string `json:"thumbnail_url"`
	ContentURL   string `json:"content_url"` // Link to LMS content (Moodle, Canvas, etc)

	Prerequisites datatypes.JSONSlice[int] `json:"prerequisites" gorm:"type:jsonb"` // Course IDs

	MaxParticipants   int `json:"max_participants"`
	CurrentEnrollment int `json:"current_enrollment"` // Calculated

	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	IsActive    bool `json:"is_active" gorm:"default:true"`
	IsMandatory bool `json:"is_mandatory"`

	Enrollments []CourseEnrollment `json:"enrollments,omitempty"`
	Modules     []CourseModule     `json:"modules,omitempty"`

	gorm.Model
}

// CourseModule individual modules within a course
type CourseModule struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CourseID uint       `json:"course_id" gorm:"index"`
	Course   *LMSCourse `json:"course,omitempty"`

	ModuleNumber int    `json:"module_number"`
	Title        string `json:"title"`
	Description  string `json:"description"`

	ContentURL      string `json:"content_url"`
	DurationMinutes int    `json:"duration_minutes"`

	OrderSequence int `json:"order_sequence"`

	gorm.Model
}

// CourseEnrollment tracks course enrollments
type CourseEnrollment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	CourseID uint       `json:"course_id" gorm:"index"`
	Course   *LMSCourse `json:"course,omitempty"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	EnrollmentDate time.Time  `json:"enrollment_date"`
	CompletionDate *time.Time `json:"completion_date"`

	Status string `json:"status"` // Enrolled, InProgress, Completed, Failed, Withdrawn

	ScorePercentage float64 `json:"score_percentage"`

	Progress float64 `json:"progress"` // 0-100

	TimeSpent int `json:"time_spent"` // Minutes

	LessonProgress []LessonProgress `json:"lesson_progress,omitempty"`

	gorm.Model
}

// LessonProgress tracks progress per lesson
type LessonProgress struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	EnrollmentID uint              `json:"enrollment_id" gorm:"index"`
	Enrollment   *CourseEnrollment `json:"enrollment,omitempty"`

	ModuleID uint          `json:"module_id" gorm:"index"`
	Module   *CourseModule `json:"module,omitempty"`

	Status        string     `json:"status"` // NotStarted, InProgress, Completed
	CompletedDate *time.Time `json:"completed_date"`

	TimeSpent int `json:"time_spent"` // Minutes

	gorm.Model
}

// Quiz course assessment quizzes
type Quiz struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	CourseID uint       `json:"course_id" gorm:"index"`
	Course   *LMSCourse `json:"course,omitempty"`

	Title       string `json:"title"`
	Description string `json:"description"`

	PassingScore float64 `json:"passing_score"` // Percentage

	AttemptsAllowed int `json:"attempts_allowed"`

	Questions []QuizQuestion `json:"questions,omitempty"`

	CreatedByID uint      `json:"created_by_id"`
	CreatedBy   *Employee `json:"created_by,omitempty" gorm:"foreignKey:CreatedByID"`

	gorm.Model
}

// QuizQuestion individual quiz questions
type QuizQuestion struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	QuizID uint  `json:"quiz_id" gorm:"index"`
	Quiz   *Quiz `json:"quiz,omitempty"`

	QuestionText string `json:"question_text"`
	QuestionType string `json:"question_type"` // MultipleChoice, ShortAnswer, Essay, TrueFalse

	Options       datatypes.JSONSlice[string] `json:"options" gorm:"type:jsonb"`
	CorrectAnswer string                      `json:"correct_answer"`

	Points        float64 `json:"points"`
	OrderSequence int     `json:"order_sequence"`

	gorm.Model
}

// QuizAttempt tracks quiz attempts
type QuizAttempt struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	QuizID uint  `json:"quiz_id" gorm:"index"`
	Quiz   *Quiz `json:"quiz,omitempty"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	AttemptNumber int        `json:"attempt_number"`
	StartTime     time.Time  `json:"start_time"`
	EndTime       *time.Time `json:"end_time"`

	Score           float64 `json:"score"`
	TotalScore      float64 `json:"total_score"`
	ScorePercentage float64 `json:"score_percentage"`

	Status string `json:"status"` // InProgress, Completed, Failed, Passed

	Answers []QuizAnswer `json:"answers,omitempty"`

	gorm.Model
}

// QuizAnswer individual answers to quiz questions
type QuizAnswer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	AttemptID uint         `json:"attempt_id" gorm:"index"`
	Attempt   *QuizAttempt `json:"attempt,omitempty"`

	QuestionID uint          `json:"question_id" gorm:"index"`
	Question   *QuizQuestion `json:"question,omitempty"`

	AnswerText   string  `json:"answer_text"`
	IsCorrect    bool    `json:"is_correct"`
	PointsEarned float64 `json:"points_earned"`

	gorm.Model
}

// Certificate course completion certificates
type Certificate struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	TenantID  string    `json:"tenant_id" gorm:"index"`

	CertificateID string `json:"certificate_id" gorm:"uniqueIndex"`

	CourseID uint       `json:"course_id" gorm:"index"`
	Course   *LMSCourse `json:"course,omitempty"`

	EmployeeID uint      `json:"employee_id" gorm:"index"`
	Employee   *Employee `json:"employee,omitempty"`

	IssuedDate time.Time  `json:"issued_date"`
	ExpiryDate *time.Time `json:"expiry_date"`

	Score float64 `json:"score"`

	CertificateURL   string `json:"certificate_url"`
	VerificationCode string `json:"verification_code" gorm:"uniqueIndex"`

	Status string `json:"status"` // Valid, Expired, Revoked

	gorm.Model
}
