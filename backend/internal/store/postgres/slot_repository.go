package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sisupass.com/sisupass/internal/repository"
	"sisupass.com/sisupass/internal/types"
)

type SlotRepository struct {
	db *sql.DB
}

func NewSlotRepository(db *sql.DB) repository.SlotRepository {
	return &SlotRepository{db: db}
}

func (r *SlotRepository) Create(slot *types.Slot) error {
	query := `INSERT INTO slots (user_id, slot_type, title, description, action, status, fields) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) 
              RETURNING slot_id, created_at, updated_at`

	args := []any{
		slot.UserID,
		slot.SlotType,
		slot.Title,
		slot.Description,
		slot.Action,
		slot.Status,
		slot.Fields,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&slot.SlotID, &slot.CreatedAt, &slot.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *SlotRepository) GetByID(slotID string) (*types.Slot, error) {
	query := `SELECT slot_id, user_id, slot_type, title, description, action, status, fields, created_at, updated_at 
              FROM slots 
              WHERE slot_id = $1`

	var slot types.Slot

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, slotID).Scan(
		&slot.SlotID,
		&slot.UserID,
		&slot.SlotType,
		&slot.Title,
		&slot.Description,
		&slot.Action,
		&slot.Status,
		&slot.Fields,
		&slot.CreatedAt,
		&slot.UpdatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errors.New("record not found")
		default:
			return nil, err
		}
	}

	return &slot, nil
}

func (r *SlotRepository) GetByUserID(userID string) ([]*types.Slot, error) {
	query := `SELECT slot_id, user_id, slot_type, title, description, action, status, fields, created_at, updated_at 
              FROM slots 
              WHERE user_id = $1 
              ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*types.Slot

	for rows.Next() {
		var slot types.Slot
		err := rows.Scan(
			&slot.SlotID,
			&slot.UserID,
			&slot.SlotType,
			&slot.Title,
			&slot.Description,
			&slot.Action,
			&slot.Status,
			&slot.Fields,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, &slot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}

func (r *SlotRepository) GetByType(userID, slotType string) ([]*types.Slot, error) {
	query := `SELECT slot_id, user_id, slot_type, title, description, action, status, fields, created_at, updated_at 
              FROM slots 
              WHERE user_id = $1 AND slot_type = $2 
              ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userID, slotType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*types.Slot

	for rows.Next() {
		var slot types.Slot
		err := rows.Scan(
			&slot.SlotID,
			&slot.UserID,
			&slot.SlotType,
			&slot.Title,
			&slot.Description,
			&slot.Action,
			&slot.Status,
			&slot.Fields,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, &slot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}

func (r *SlotRepository) Update(slot *types.Slot) error {
	query := `UPDATE slots 
              SET slot_type = $2, title = $3, description = $4, action = $5, status = $6, fields = $7, updated_at = CURRENT_TIMESTAMP 
              WHERE slot_id = $1 
              RETURNING updated_at`

	args := []any{
		slot.SlotID,
		slot.SlotType,
		slot.Title,
		slot.Description,
		slot.Action,
		slot.Status,
		slot.Fields,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&slot.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return errors.New("record not found")
		default:
			return err
		}
	}

	return nil
}

func (r *SlotRepository) Delete(slotID string) error {
	query := `DELETE FROM slots WHERE slot_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, slotID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return nil
}

func (r *SlotRepository) Search(userID, query string) ([]*types.Slot, error) {
	sqlQuery := `SELECT slot_id, user_id, slot_type, title, description, action, status, fields, created_at, updated_at 
                 FROM slots 
                 WHERE user_id = $1 AND (
                     title ILIKE '%' || $2 || '%' OR 
                     slot_type ILIKE '%' || $2 || '%' OR 
                     action ILIKE '%' || $2 || '%' OR
                     (description IS NOT NULL AND description ILIKE '%' || $2 || '%')
                 )
                 ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, sqlQuery, userID, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var slots []*types.Slot

	for rows.Next() {
		var slot types.Slot
		err := rows.Scan(
			&slot.SlotID,
			&slot.UserID,
			&slot.SlotType,
			&slot.Title,
			&slot.Description,
			&slot.Action,
			&slot.Status,
			&slot.Fields,
			&slot.CreatedAt,
			&slot.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		slots = append(slots, &slot)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return slots, nil
}
