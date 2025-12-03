package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aslotsu/monkreflections-form-api/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EventHandler struct {
	db *pgxpool.Pool
}

func NewEventHandler(db *pgxpool.Pool) *EventHandler {
	return &EventHandler{db: db}
}

// GetAllEvents retrieves all events
func (h *EventHandler) GetAllEvents(c *gin.Context) {
	query := `
		SELECT id, title, description, event_type, status, start_date, end_date,
		       venue_name, venue_address, is_virtual, virtual_link, timezone,
		       capacity, expected_guests, registered_count, actual_guests,
		       waitlist_enabled, allow_walkins, ticket_price, early_bird_price,
		       organization_budget, expenses, revenue,
		       registration_open_date, registration_close_date, registration_form_url,
		       requires_approval, featured_image, gallery_images, video_url, livestream_url,
		       organizer_name, organizer_email, organizer_phone,
		       speakers, sponsors, tags, is_featured, is_public, created_by,
		       created_at, updated_at
		FROM events
		ORDER BY start_date DESC
	`

	rows, err := h.db.Query(context.Background(), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
		return
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var event models.Event
		if err := rows.Scan(
			&event.ID, &event.Title, &event.Description, &event.EventType, &event.Status,
			&event.StartDate, &event.EndDate, &event.VenueName, &event.VenueAddress,
			&event.IsVirtual, &event.VirtualLink, &event.Timezone,
			&event.Capacity, &event.ExpectedGuests, &event.RegisteredCount, &event.ActualGuests,
			&event.WaitlistEnabled, &event.AllowWalkins, &event.TicketPrice, &event.EarlyBirdPrice,
			&event.OrganizationBudget, &event.Expenses, &event.Revenue,
			&event.RegistrationOpenDate, &event.RegistrationCloseDate, &event.RegistrationFormURL,
			&event.RequiresApproval, &event.FeaturedImage, &event.GalleryImages, &event.VideoURL,
			&event.LivestreamURL, &event.OrganizerName, &event.OrganizerEmail, &event.OrganizerPhone,
			&event.Speakers, &event.Sponsors, &event.Tags, &event.IsFeatured, &event.IsPublic,
			&event.CreatedBy, &event.CreatedAt, &event.UpdatedAt,
		); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan event"})
			return
		}
		events = append(events, event)
	}

	if events == nil {
		events = []models.Event{}
	}

	c.JSON(http.StatusOK, events)
}

// GetEventByID retrieves a single event by ID
func (h *EventHandler) GetEventByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	query := `
		SELECT id, title, description, event_type, status, start_date, end_date,
		       venue_name, venue_address, is_virtual, virtual_link, timezone,
		       capacity, expected_guests, registered_count, actual_guests,
		       waitlist_enabled, allow_walkins, ticket_price, early_bird_price,
		       organization_budget, expenses, revenue,
		       registration_open_date, registration_close_date, registration_form_url,
		       requires_approval, featured_image, gallery_images, video_url, livestream_url,
		       organizer_name, organizer_email, organizer_phone,
		       speakers, sponsors, tags, is_featured, is_public, created_by,
		       created_at, updated_at
		FROM events
		WHERE id = $1
	`

	var event models.Event
	err = h.db.QueryRow(context.Background(), query, id).Scan(
		&event.ID, &event.Title, &event.Description, &event.EventType, &event.Status,
		&event.StartDate, &event.EndDate, &event.VenueName, &event.VenueAddress,
		&event.IsVirtual, &event.VirtualLink, &event.Timezone,
		&event.Capacity, &event.ExpectedGuests, &event.RegisteredCount, &event.ActualGuests,
		&event.WaitlistEnabled, &event.AllowWalkins, &event.TicketPrice, &event.EarlyBirdPrice,
		&event.OrganizationBudget, &event.Expenses, &event.Revenue,
		&event.RegistrationOpenDate, &event.RegistrationCloseDate, &event.RegistrationFormURL,
		&event.RequiresApproval, &event.FeaturedImage, &event.GalleryImages, &event.VideoURL,
		&event.LivestreamURL, &event.OrganizerName, &event.OrganizerEmail, &event.OrganizerPhone,
		&event.Speakers, &event.Sponsors, &event.Tags, &event.IsFeatured, &event.IsPublic,
		&event.CreatedBy, &event.CreatedAt, &event.UpdatedAt,
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// CreateEvent creates a new event
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req models.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert JSONB fields to strings
	galleryImagesStr := marshalJSONB(req.GalleryImages)
	speakersStr := marshalJSONB(req.Speakers)
	sponsorsStr := marshalJSONB(req.Sponsors)
	tagsStr := marshalJSONB(req.Tags)

	query := `
		INSERT INTO events (
			title, description, event_type, status, start_date, end_date,
			venue_name, venue_address, is_virtual, virtual_link, timezone,
			capacity, expected_guests, registered_count, actual_guests,
			waitlist_enabled, allow_walkins, ticket_price, early_bird_price,
			organization_budget, expenses, revenue,
			registration_open_date, registration_close_date, registration_form_url,
			requires_approval, featured_image, gallery_images, video_url, livestream_url,
			organizer_name, organizer_email, organizer_phone,
			speakers, sponsors, tags, is_featured, is_public, created_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
			$33, $34, $35, $36, $37, $38, $39
		) RETURNING id
	`

	var id int
	err := h.db.QueryRow(
		context.Background(),
		query,
		req.Title, req.Description, req.EventType, req.Status, req.StartDate, req.EndDate,
		req.VenueName, req.VenueAddress, req.IsVirtual, req.VirtualLink, req.Timezone,
		req.Capacity, req.ExpectedGuests, req.RegisteredCount, req.ActualGuests,
		req.WaitlistEnabled, req.AllowWalkins, req.TicketPrice, req.EarlyBirdPrice,
		req.OrganizationBudget, req.Expenses, req.Revenue,
		req.RegistrationOpenDate, req.RegistrationCloseDate, req.RegistrationFormURL,
		req.RequiresApproval, req.FeaturedImage, galleryImagesStr, req.VideoURL, req.LivestreamURL,
		req.OrganizerName, req.OrganizerEmail, req.OrganizerPhone,
		speakersStr, sponsorsStr, tagsStr, req.IsFeatured, req.IsPublic, req.CreatedBy,
	).Scan(&id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create event", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateEvent updates an existing event
func (h *EventHandler) UpdateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var req models.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if event exists
	var exists bool
	err = h.db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM events WHERE id = $1)", id).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	// Build dynamic update query
	query := "UPDATE events SET updated_at = CURRENT_TIMESTAMP"
	args := []any{}
	argCount := 1

	if req.Title != "" {
		query += ", title = $" + strconv.Itoa(argCount)
		args = append(args, req.Title)
		argCount++
	}
	if req.Description != "" {
		query += ", description = $" + strconv.Itoa(argCount)
		args = append(args, req.Description)
		argCount++
	}
	if req.EventType != "" {
		query += ", event_type = $" + strconv.Itoa(argCount)
		args = append(args, req.EventType)
		argCount++
	}
	if req.Status != "" {
		query += ", status = $" + strconv.Itoa(argCount)
		args = append(args, req.Status)
		argCount++
	}
	if req.StartDate != nil {
		query += ", start_date = $" + strconv.Itoa(argCount)
		args = append(args, req.StartDate)
		argCount++
	}
	if req.EndDate != nil {
		query += ", end_date = $" + strconv.Itoa(argCount)
		args = append(args, req.EndDate)
		argCount++
	}
	if req.VenueName != "" {
		query += ", venue_name = $" + strconv.Itoa(argCount)
		args = append(args, req.VenueName)
		argCount++
	}
	if req.VenueAddress != "" {
		query += ", venue_address = $" + strconv.Itoa(argCount)
		args = append(args, req.VenueAddress)
		argCount++
	}
	if req.IsVirtual != nil {
		query += ", is_virtual = $" + strconv.Itoa(argCount)
		args = append(args, req.IsVirtual)
		argCount++
	}
	if req.Capacity != nil {
		query += ", capacity = $" + strconv.Itoa(argCount)
		args = append(args, req.Capacity)
		argCount++
	}
	if req.TicketPrice != nil {
		query += ", ticket_price = $" + strconv.Itoa(argCount)
		args = append(args, req.TicketPrice)
		argCount++
	}
	if req.IsFeatured != nil {
		query += ", is_featured = $" + strconv.Itoa(argCount)
		args = append(args, req.IsFeatured)
		argCount++
	}
	if req.IsPublic != nil {
		query += ", is_public = $" + strconv.Itoa(argCount)
		args = append(args, req.IsPublic)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, id)

	_, err = h.db.Exec(context.Background(), query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// DeleteEvent deletes an event
func (h *EventHandler) DeleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	result, err := h.db.Exec(context.Background(), "DELETE FROM events WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// Helper function to marshal JSONB data
func marshalJSONB(data any) string {
	if data == nil {
		return "{}"
	}
	bytes, err := json.Marshal(data)
	if err != nil {
		return "{}"
	}
	return string(bytes)
}
