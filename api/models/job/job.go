package job

import "time"

type JobResponse struct {
	ID             *string    `json:"id"`
	Status         *string    `json:"status,omitempty" validate:"omitempty,oneof=completed failed pending"`
	AccountId      *string    `json:"account_id,omitempty"`
	ResourceId     *string    `json:"resource_id,omitempty"`
	ResourceType   *string    `json:"resource_type,omitempty"`
	ResourceAction *string    `json:"resource_action,omitempty"`
	ErrorMessage   *string    `json:"error_message,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
}

type JobCreationResponse struct {
	ID *string `json:"job_id"`
}
