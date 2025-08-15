package validator

import (
	"encoding/json"
	"fmt"

	"sisupass.com/sisupass/internal/types"
)

// ValidateSlot validates the complete slot structure
func ValidateSlot(v *Validator, slot *types.Slot) {
	// Basic field validation
	v.Check(slot.Title != "", "title", "must be provided")
	v.Check(len(slot.Title) <= 100, "title", "must not be more than 100 bytes long")

	v.Check(slot.SlotType != "", "slot_type", "must be provided")
	v.Check(IsValidSlotType(slot.SlotType), "slot_type", "must be a valid slot type")

	v.Check(slot.Action != "", "action", "must be provided")
	v.Check(IsValidSlotAction(slot.Action), "action", "must be a valid action")

	// Optional description validation
	if slot.Description != nil && *slot.Description != "" {
		v.Check(len(*slot.Description) <= 500, "description", "must not be more than 500 bytes long")
	}

	// Status validation (if provided)
	if slot.Status != "" {
		v.Check(IsValidSlotStatus(slot.Status), "status", "must be a valid status")
	}

	// Fields JSON validation
	if slot.Fields != nil {
		err := ValidateSlotFields(slot.Fields)
		v.Check(err == nil, "fields", fmt.Sprintf("invalid fields: %v", err))
	}
}

// ValidateSlotUpdate validates slot update data
func ValidateSlotUpdate(v *Validator, updates *types.Slot) {
	// Only validate provided fields for updates
	if updates.Title != "" {
		v.Check(len(updates.Title) <= 100, "title", "must not be more than 100 bytes long")
	}

	if updates.Description != nil && *updates.Description != "" {
		v.Check(len(*updates.Description) <= 500, "description", "must not be more than 500 bytes long")
	}

	if updates.SlotType != "" {
		v.Check(IsValidSlotType(updates.SlotType), "slot_type", "must be a valid slot type")
	}

	if updates.Action != "" {
		v.Check(IsValidSlotAction(updates.Action), "action", "must be a valid action")
	}

	if updates.Status != "" {
		v.Check(IsValidSlotStatus(updates.Status), "status", "must be a valid status")
	}

	if updates.Fields != nil {
		err := ValidateSlotFields(updates.Fields)
		v.Check(err == nil, "fields", fmt.Sprintf("invalid fields: %v", err))
	}
}

// IsValidSlotType checks if the slot type is valid
func IsValidSlotType(slotType string) bool {
	validTypes := map[string]bool{
		"nfc":     true,
		"qr":      true,
		"auth":    true,
		"data":    true,
		"trigger": true,
		"payment": true,
		"access":  true,
	}
	return validTypes[slotType]
}

// IsValidSlotAction checks if the slot action is valid
func IsValidSlotAction(action string) bool {
	validActions := map[string]bool{
		"trigger":      true,
		"authenticate": true,
		"authorize":    true,
		"data_read":    true,
		"data_write":   true,
		"payment":      true,
		"unlock":       true,
		"lock":         true,
	}
	return validActions[action]
}

// IsValidSlotStatus checks if the slot status is valid
func IsValidSlotStatus(status string) bool {
	validStatuses := map[string]bool{
		"active":   true,
		"inactive": true,
		"expired":  true,
		"locked":   true,
	}
	return validStatuses[status]
}

// ValidateSlotFields validates the fields JSON structure
func ValidateSlotFields(fields []byte) error {
	if len(fields) == 0 {
		return nil
	}

	// Basic JSON validation
	var temp interface{}
	if err := json.Unmarshal(fields, &temp); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	// Ensure it's an object, not array or primitive
	fieldsMap, ok := temp.(map[string]interface{})
	if !ok {
		return fmt.Errorf("fields must be a JSON object")
	}

	// Additional validation based on common field types
	return validateFieldsContent(fieldsMap)
}

// validateFieldsContent validates the content of fields based on common patterns
func validateFieldsContent(fields map[string]interface{}) error {
	// Validate common field patterns
	for key, value := range fields {
		switch key {
		case "url":
			if str, ok := value.(string); ok {
				if len(str) > 2048 {
					return fmt.Errorf("url field too long")
				}
			}
		case "email":
			if str, ok := value.(string); ok {
				if !isValidEmail(str) {
					return fmt.Errorf("invalid email format")
				}
			}
		case "timeout":
			if num, ok := value.(float64); ok {
				if num < 0 || num > 3600 {
					return fmt.Errorf("timeout must be between 0 and 3600 seconds")
				}
			}
		case "max_attempts":
			if num, ok := value.(float64); ok {
				if num < 1 || num > 10 {
					return fmt.Errorf("max_attempts must be between 1 and 10")
				}
			}
		}

		// Validate nested objects recursively
		if nested, ok := value.(map[string]interface{}); ok {
			if err := validateFieldsContent(nested); err != nil {
				return fmt.Errorf("nested field %s: %w", key, err)
			}
		}
	}

	return nil
}

// isValidEmail performs basic email validation
func isValidEmail(email string) bool {
	// Basic email regex pattern
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	// Must contain @ symbol
	atCount := 0
	for _, char := range email {
		if char == '@' {
			atCount++
		}
	}

	return atCount == 1
}
