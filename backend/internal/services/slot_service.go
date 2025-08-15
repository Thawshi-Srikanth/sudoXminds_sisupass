package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"sisupass.com/sisupass/internal/repository"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

var (
	ErrInvalidSlotType   = errors.New("invalid slot type")
	ErrInvalidSlotAction = errors.New("invalid slot action")
	ErrSlotNotFound      = errors.New("slot not found")
	ErrSlotNotActive     = errors.New("slot is not active")
)

type SlotService struct {
	slotRepo repository.SlotRepository
}

func NewSlotService(slotRepo repository.SlotRepository) *SlotService {
	return &SlotService{
		slotRepo: slotRepo,
	}
}

// CreateSlot creates a new slot with validation
func (s *SlotService) CreateSlot(ctx context.Context, userID uuid.UUID, slot *types.Slot) (*types.Slot, error) {
	// Validate required fields
	v := validator.New()

	v.Check(slot.Title != "", "title", "must be provided")
	v.Check(len(slot.Title) <= 100, "title", "must not be more than 100 bytes long")
	v.Check(slot.SlotType != "", "slot_type", "must be provided")
	v.Check(validator.IsValidSlotType(slot.SlotType), "slot_type", "must be a valid slot type")
	v.Check(slot.Action != "", "action", "must be provided")
	v.Check(validator.IsValidSlotAction(slot.Action), "action", "must be a valid action")

	// Validate fields JSON if provided
	if slot.Fields != nil {
		if err := validator.ValidateSlotFields(slot.Fields); err != nil {
			v.Check(false, "fields", err.Error())
		}
	}

	if !v.Valid() {
		return nil, fmt.Errorf("validation failed: %v", v.Errors)
	}

	// Set user ID and timestamps
	slot.UserID = userID
	now := time.Now()
	slot.CreatedAt = now
	slot.UpdatedAt = now

	// Set default status if not provided
	if slot.Status == "" {
		slot.Status = "active"
	}

	// Create slot
	err := s.slotRepo.Create(slot)
	if err != nil {
		return nil, fmt.Errorf("failed to create slot: %w", err)
	}

	return slot, nil
}

// GetSlot retrieves a slot by ID with ownership validation
func (s *SlotService) GetSlot(ctx context.Context, userID uuid.UUID, slotID string) (*types.Slot, error) {
	slot, err := s.slotRepo.GetByID(slotID)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, ErrSlotNotFound
		}
		return nil, fmt.Errorf("failed to get slot: %w", err)
	}

	// Check ownership
	if slot.UserID != userID {
		return nil, ErrSlotNotFound // Don't reveal existence of slot to unauthorized users
	}

	return slot, nil
}

// GetUserSlots retrieves all slots for a user with optional filters
func (s *SlotService) GetUserSlots(ctx context.Context, userID uuid.UUID, slotType, status string) ([]*types.Slot, error) {
	var slots []*types.Slot
	var err error

	if slotType != "" {
		slots, err = s.slotRepo.GetByType(userID.String(), slotType)
	} else {
		slots, err = s.slotRepo.GetByUserID(userID.String())
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user slots: %w", err)
	}

	// Filter by status if provided
	if status != "" {
		filteredSlots := make([]*types.Slot, 0)
		for _, slot := range slots {
			if slot.Status == status {
				filteredSlots = append(filteredSlots, slot)
			}
		}
		return filteredSlots, nil
	}

	return slots, nil
}

// UpdateSlot updates an existing slot with validation
func (s *SlotService) UpdateSlot(ctx context.Context, userID uuid.UUID, slotID string, updates *types.Slot) (*types.Slot, error) {
	// Get existing slot and verify ownership
	existingSlot, err := s.GetSlot(ctx, userID, slotID)
	if err != nil {
		return nil, err
	}

	// Validate updates
	v := validator.New()

	if updates.Title != "" {
		v.Check(len(updates.Title) <= 100, "title", "must not be more than 100 bytes long")
		existingSlot.Title = updates.Title
	}

	if updates.Description != nil && *updates.Description != "" {
		existingSlot.Description = updates.Description
	}

	if updates.SlotType != "" {
		v.Check(validator.IsValidSlotType(updates.SlotType), "slot_type", "must be a valid slot type")
		existingSlot.SlotType = updates.SlotType
	}

	if updates.Action != "" {
		v.Check(validator.IsValidSlotAction(updates.Action), "action", "must be a valid action")
		existingSlot.Action = updates.Action
	}

	if updates.Status != "" {
		v.Check(validator.IsValidSlotStatus(updates.Status), "status", "must be a valid status")
		existingSlot.Status = updates.Status
	}

	if updates.Fields != nil {
		if err := validator.ValidateSlotFields(updates.Fields); err != nil {
			v.Check(false, "fields", err.Error())
		}
		existingSlot.Fields = updates.Fields
	}

	if !v.Valid() {
		return nil, fmt.Errorf("validation failed: %v", v.Errors)
	}

	// Update timestamp
	existingSlot.UpdatedAt = time.Now()

	// Update slot
	err = s.slotRepo.Update(existingSlot)
	if err != nil {
		return nil, fmt.Errorf("failed to update slot: %w", err)
	}

	return existingSlot, nil
}

// DeleteSlot deletes a slot with ownership validation
func (s *SlotService) DeleteSlot(ctx context.Context, userID uuid.UUID, slotID string) error {
	// Verify ownership first
	_, err := s.GetSlot(ctx, userID, slotID)
	if err != nil {
		return err
	}

	// Delete slot
	err = s.slotRepo.Delete(slotID)
	if err != nil {
		return fmt.Errorf("failed to delete slot: %w", err)
	}

	return nil
}

// SearchSlots searches for slots with given query
func (s *SlotService) SearchSlots(ctx context.Context, userID uuid.UUID, query string) ([]*types.Slot, error) {
	if query == "" {
		return s.GetUserSlots(ctx, userID, "", "")
	}

	slots, err := s.slotRepo.Search(userID.String(), query)
	if err != nil {
		return nil, fmt.Errorf("failed to search slots: %w", err)
	}

	return slots, nil
}

// ExecuteSlotAction executes the action defined in a slot
func (s *SlotService) ExecuteSlotAction(ctx context.Context, userID uuid.UUID, slotID string) (map[string]interface{}, error) {
	slot, err := s.GetSlot(ctx, userID, slotID)
	if err != nil {
		return nil, err
	}

	// Check if slot is active
	if slot.Status != "active" {
		return nil, ErrSlotNotActive
	}

	// Execute action based on type
	result := make(map[string]interface{})
	result["slot_id"] = slot.SlotID
	result["action"] = slot.Action
	result["executed_at"] = time.Now()

	switch slot.Action {
	case "trigger":
		result["message"] = "Slot triggered successfully"
		result["status"] = "triggered"
	case "authenticate":
		result["message"] = "Authentication requested"
		result["status"] = "pending_auth"
	case "authorize":
		result["message"] = "Authorization requested"
		result["status"] = "pending_auth"
	case "data_read":
		result["message"] = "Data read initiated"
		result["status"] = "reading"
	case "data_write":
		result["message"] = "Data write initiated"
		result["status"] = "writing"
	default:
		result["message"] = "Action executed"
		result["status"] = "completed"
	}

	// Include slot fields in response if available
	if slot.Fields != nil {
		result["fields"] = string(slot.Fields)
	}

	return result, nil
}
