package models

import "time"

type Event struct {
	ID                    int       `json:"id"`
	Title                 string    `json:"title"`
	Description           string    `json:"description"`
	EventType             string    `json:"event_type"`
	Status                string    `json:"status"`
	StartDate             time.Time `json:"start_date"`
	EndDate               time.Time `json:"end_date"`
	VenueName             string    `json:"venue_name"`
	VenueAddress          string    `json:"venue_address"`
	IsVirtual             bool      `json:"is_virtual"`
	VirtualLink           string    `json:"virtual_link,omitempty"`
	Timezone              string    `json:"timezone"`
	Capacity              int       `json:"capacity"`
	ExpectedGuests        int       `json:"expected_guests"`
	RegisteredCount       int       `json:"registered_count"`
	ActualGuests          *int      `json:"actual_guests,omitempty"`
	WaitlistEnabled       bool      `json:"waitlist_enabled"`
	AllowWalkins          bool      `json:"allow_walkins"`
	TicketPrice           float64   `json:"ticket_price"`
	EarlyBirdPrice        *float64  `json:"early_bird_price,omitempty"`
	OrganizationBudget    float64   `json:"organization_budget"`
	Expenses              float64   `json:"expenses"`
	Revenue               float64   `json:"revenue"`
	RegistrationOpenDate  time.Time `json:"registration_open_date"`
	RegistrationCloseDate time.Time `json:"registration_close_date"`
	RegistrationFormURL   string    `json:"registration_form_url,omitempty"`
	RequiresApproval      bool      `json:"requires_approval"`
	FeaturedImage         string    `json:"featured_image,omitempty"`
	GalleryImages         string    `json:"gallery_images,omitempty"` // JSONB stored as string
	VideoURL              string    `json:"video_url,omitempty"`
	LivestreamURL         string    `json:"livestream_url,omitempty"`
	OrganizerName         string    `json:"organizer_name"`
	OrganizerEmail        string    `json:"organizer_email"`
	OrganizerPhone        string    `json:"organizer_phone"`
	Speakers              string    `json:"speakers,omitempty"` // JSONB stored as string
	Sponsors              string    `json:"sponsors,omitempty"` // JSONB stored as string
	Tags                  string    `json:"tags,omitempty"`     // JSONB stored as string
	IsFeatured            bool      `json:"is_featured"`
	IsPublic              bool      `json:"is_public"`
	CreatedBy             string    `json:"created_by,omitempty"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

type CreateEventRequest struct {
	Title                 string    `json:"title" binding:"required"`
	Description           string    `json:"description" binding:"required"`
	EventType             string    `json:"event_type" binding:"required"`
	Status                string    `json:"status"`
	StartDate             time.Time `json:"start_date" binding:"required"`
	EndDate               time.Time `json:"end_date" binding:"required"`
	VenueName             string    `json:"venue_name"`
	VenueAddress          string    `json:"venue_address"`
	IsVirtual             bool      `json:"is_virtual"`
	VirtualLink           string    `json:"virtual_link,omitempty"`
	Timezone              string    `json:"timezone"`
	Capacity              int       `json:"capacity"`
	ExpectedGuests        int       `json:"expected_guests"`
	RegisteredCount       int       `json:"registered_count"`
	ActualGuests          *int      `json:"actual_guests,omitempty"`
	WaitlistEnabled       bool      `json:"waitlist_enabled"`
	AllowWalkins          bool      `json:"allow_walkins"`
	TicketPrice           float64   `json:"ticket_price"`
	EarlyBirdPrice        *float64  `json:"early_bird_price,omitempty"`
	OrganizationBudget    float64   `json:"organization_budget"`
	Expenses              float64   `json:"expenses"`
	Revenue               float64   `json:"revenue"`
	RegistrationOpenDate  time.Time `json:"registration_open_date"`
	RegistrationCloseDate time.Time `json:"registration_close_date"`
	RegistrationFormURL   string    `json:"registration_form_url,omitempty"`
	RequiresApproval      bool      `json:"requires_approval"`
	FeaturedImage         string    `json:"featured_image,omitempty"`
	GalleryImages         any       `json:"gallery_images,omitempty"` // Can be array or JSON
	VideoURL              string    `json:"video_url,omitempty"`
	LivestreamURL         string    `json:"livestream_url,omitempty"`
	OrganizerName         string    `json:"organizer_name" binding:"required"`
	OrganizerEmail        string    `json:"organizer_email" binding:"required,email"`
	OrganizerPhone        string    `json:"organizer_phone"`
	Speakers              any       `json:"speakers,omitempty"` // Can be array or JSON
	Sponsors              any       `json:"sponsors,omitempty"` // Can be array or JSON
	Tags                  any       `json:"tags,omitempty"`     // Can be array or JSON
	IsFeatured            bool      `json:"is_featured"`
	IsPublic              bool      `json:"is_public"`
	CreatedBy             string    `json:"created_by,omitempty"`
}

type UpdateEventRequest struct {
	Title                 string    `json:"title,omitempty"`
	Description           string    `json:"description,omitempty"`
	EventType             string    `json:"event_type,omitempty"`
	Status                string    `json:"status,omitempty"`
	StartDate             *time.Time `json:"start_date,omitempty"`
	EndDate               *time.Time `json:"end_date,omitempty"`
	VenueName             string    `json:"venue_name,omitempty"`
	VenueAddress          string    `json:"venue_address,omitempty"`
	IsVirtual             *bool     `json:"is_virtual,omitempty"`
	VirtualLink           string    `json:"virtual_link,omitempty"`
	Timezone              string    `json:"timezone,omitempty"`
	Capacity              *int      `json:"capacity,omitempty"`
	ExpectedGuests        *int      `json:"expected_guests,omitempty"`
	RegisteredCount       *int      `json:"registered_count,omitempty"`
	ActualGuests          *int      `json:"actual_guests,omitempty"`
	WaitlistEnabled       *bool     `json:"waitlist_enabled,omitempty"`
	AllowWalkins          *bool     `json:"allow_walkins,omitempty"`
	TicketPrice           *float64  `json:"ticket_price,omitempty"`
	EarlyBirdPrice        *float64  `json:"early_bird_price,omitempty"`
	OrganizationBudget    *float64  `json:"organization_budget,omitempty"`
	Expenses              *float64  `json:"expenses,omitempty"`
	Revenue               *float64  `json:"revenue,omitempty"`
	RegistrationOpenDate  *time.Time `json:"registration_open_date,omitempty"`
	RegistrationCloseDate *time.Time `json:"registration_close_date,omitempty"`
	RegistrationFormURL   string    `json:"registration_form_url,omitempty"`
	RequiresApproval      *bool     `json:"requires_approval,omitempty"`
	FeaturedImage         string    `json:"featured_image,omitempty"`
	GalleryImages         any       `json:"gallery_images,omitempty"`
	VideoURL              string    `json:"video_url,omitempty"`
	LivestreamURL         string    `json:"livestream_url,omitempty"`
	OrganizerName         string    `json:"organizer_name,omitempty"`
	OrganizerEmail        string    `json:"organizer_email,omitempty"`
	OrganizerPhone        string    `json:"organizer_phone,omitempty"`
	Speakers              any       `json:"speakers,omitempty"`
	Sponsors              any       `json:"sponsors,omitempty"`
	Tags                  any       `json:"tags,omitempty"`
	IsFeatured            *bool     `json:"is_featured,omitempty"`
	IsPublic              *bool     `json:"is_public,omitempty"`
}
