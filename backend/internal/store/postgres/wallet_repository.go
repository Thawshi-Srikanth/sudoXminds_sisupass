package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"sisupass.com/sisupass/internal/repository"
	"sisupass.com/sisupass/internal/types"
)

type WalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) repository.WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Create(wallet *types.Wallet) error {
	query := `INSERT INTO wallets (user_id, balance, wallet_type, parent_wallet_id, access_qr_code, access_nfc_code) 
              VALUES ($1, $2, $3, $4, $5, $6) 
              RETURNING wallet_id, created_at`

	args := []any{
		wallet.UserID,
		wallet.Balance,
		wallet.WalletType,
		wallet.ParentWalletID,
		wallet.AccessQRCode,
		wallet.AccessNFCCode,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&wallet.WalletID, &wallet.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletRepository) GetByID(walletID string) (*types.Wallet, error) {
	query := `SELECT wallet_id, user_id, balance, wallet_type, parent_wallet_id, access_qr_code, access_nfc_code, created_at 
              FROM wallets 
              WHERE wallet_id = $1`

	var wallet types.Wallet

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, walletID).Scan(
		&wallet.WalletID,
		&wallet.UserID,
		&wallet.Balance,
		&wallet.WalletType,
		&wallet.ParentWalletID,
		&wallet.AccessQRCode,
		&wallet.AccessNFCCode,
		&wallet.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, types.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &wallet, nil
}

func (r *WalletRepository) GetByUserID(userID string) ([]*types.Wallet, error) {
	query := `SELECT wallet_id, user_id, balance, wallet_type, parent_wallet_id, access_qr_code, access_nfc_code, created_at 
              FROM wallets 
              WHERE user_id = $1 
              ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []*types.Wallet

	for rows.Next() {
		var wallet types.Wallet
		err := rows.Scan(
			&wallet.WalletID,
			&wallet.UserID,
			&wallet.Balance,
			&wallet.WalletType,
			&wallet.ParentWalletID,
			&wallet.AccessQRCode,
			&wallet.AccessNFCCode,
			&wallet.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, &wallet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return wallets, nil
}

func (r *WalletRepository) Update(wallet *types.Wallet) error {
	query := `UPDATE wallets 
              SET balance = $1, wallet_type = $2, parent_wallet_id = $3, access_qr_code = $4, access_nfc_code = $5
              WHERE wallet_id = $6`

	args := []any{
		wallet.Balance,
		wallet.WalletType,
		wallet.ParentWalletID,
		wallet.AccessQRCode,
		wallet.AccessNFCCode,
		wallet.WalletID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return types.ErrRecordNotFound
	}

	return nil
}

func (r *WalletRepository) Delete(walletID string) error {
	query := `DELETE FROM wallets WHERE wallet_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, walletID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return types.ErrRecordNotFound
	}

	return nil
}

func (r *WalletRepository) GetBalance(walletID string) (float64, error) {
	query := `SELECT balance FROM wallets WHERE wallet_id = $1`

	var balance float64

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, walletID).Scan(&balance)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return 0, types.ErrRecordNotFound
		default:
			return 0, err
		}
	}

	return balance, nil
}

func (r *WalletRepository) UpdateBalance(walletID string, amount float64) error {
	query := `UPDATE wallets SET balance = balance + $1 WHERE wallet_id = $2 AND (balance + $1) >= 0`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, amount, walletID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("insufficient balance or wallet not found")
	}

	return nil
}
