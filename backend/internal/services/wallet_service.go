package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"sisupass.com/sisupass/internal/data"
	"sisupass.com/sisupass/internal/types"
	"sisupass.com/sisupass/internal/validator"
)

type WalletService struct {
	models *data.Models
}

func NewWalletService(models *data.Models) *WalletService {
	return &WalletService{models: models}
}

type CreateWalletRequest struct {
	WalletType     string     `json:"wallet_type"`
	Balance        float64    `json:"balance"`
	ParentWalletID *uuid.UUID `json:"parent_wallet_id,omitempty"`
}

type WalletResponse struct {
	WalletID       uuid.UUID  `json:"wallet_id"`
	UserID         uuid.UUID  `json:"user_id"`
	Balance        float64    `json:"balance"`
	WalletType     string     `json:"wallet_type"`
	ParentWalletID *uuid.UUID `json:"parent_wallet_id,omitempty"`
	AccessQRCode   *string    `json:"access_qr_code,omitempty"`
	AccessNFCCode  *string    `json:"access_nfc_code,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type UpdateBalanceRequest struct {
	Amount float64 `json:"amount"`
}

// QRCodeData represents the structured data stored in QR codes
type QRCodeData struct {
	WalletID   string    `json:"wallet_id"`
	UserID     string    `json:"user_id"`
	WalletType string    `json:"wallet_type"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Version    string    `json:"version"`
	Signature  string    `json:"signature"`
}

// NFCCodeData represents the structured data stored in NFC chips
type NFCCodeData struct {
	WalletID   string    `json:"wallet_id"`
	UserID     string    `json:"user_id"`
	WalletType string    `json:"wallet_type"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Version    string    `json:"version"`
	AccessKey  string    `json:"access_key"`
}

func (s *WalletService) CreateWallet(ctx context.Context, userID uuid.UUID, req CreateWalletRequest) (*WalletResponse, error) {
	v := validator.New()
	validator.ValidateCreateWalletRequest(v, req.WalletType, req.Balance)

	if !v.Valid() {
		return nil, v.ToError()
	}

	// Verify user exists
	_, err := s.models.Users.GetByID(userID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Generate access codes with structured data
	qrCode, err := generateQRCode(userID, req.WalletType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	nfcCode, err := generateNFCCode(userID, req.WalletType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate NFC code: %w", err)
	}

	wallet := &types.Wallet{
		UserID:         userID,
		Balance:        req.Balance,
		WalletType:     req.WalletType,
		ParentWalletID: req.ParentWalletID,
		AccessQRCode:   &qrCode,
		AccessNFCCode:  &nfcCode,
	}

	walletValidator := validator.New()
	validator.ValidateWallet(walletValidator, wallet)

	if !walletValidator.Valid() {
		return nil, walletValidator.ToError()
	}

	err = s.models.Wallets.Create(wallet)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	// Update the QR and NFC codes with the actual wallet ID
	updatedQRCode, err := updateQRCodeWithWalletID(qrCode, wallet.WalletID.String())
	if err != nil {
		// Log error but don't fail the creation
		// In production, you might want to handle this differently
	} else {
		wallet.AccessQRCode = &updatedQRCode
	}

	updatedNFCCode, err := updateNFCCodeWithWalletID(nfcCode, wallet.WalletID.String())
	if err != nil {
		// Log error but don't fail the creation
	} else {
		wallet.AccessNFCCode = &updatedNFCCode

		// Update the wallet in database with correct codes
		_ = s.models.Wallets.Update(wallet)
	}

	return &WalletResponse{
		WalletID:       wallet.WalletID,
		UserID:         wallet.UserID,
		Balance:        wallet.Balance,
		WalletType:     wallet.WalletType,
		ParentWalletID: wallet.ParentWalletID,
		AccessQRCode:   wallet.AccessQRCode,
		AccessNFCCode:  wallet.AccessNFCCode,
		CreatedAt:      wallet.CreatedAt,
	}, nil
}

func (s *WalletService) GetWallet(ctx context.Context, walletID string) (*WalletResponse, error) {
	v := validator.New()
	validator.ValidateWalletID(v, walletID)

	if !v.Valid() {
		return nil, v.ToError()
	}

	wallet, err := s.models.Wallets.GetByID(walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &WalletResponse{
		WalletID:       wallet.WalletID,
		UserID:         wallet.UserID,
		Balance:        wallet.Balance,
		WalletType:     wallet.WalletType,
		ParentWalletID: wallet.ParentWalletID,
		AccessQRCode:   wallet.AccessQRCode,
		AccessNFCCode:  wallet.AccessNFCCode,
		CreatedAt:      wallet.CreatedAt,
	}, nil
}

func (s *WalletService) GetUserWallets(ctx context.Context, userID string) ([]*WalletResponse, error) {
	// Verify user exists
	_, err := s.models.Users.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	wallets, err := s.models.Wallets.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user wallets: %w", err)
	}

	var walletResponses []*WalletResponse
	for _, wallet := range wallets {
		walletResponses = append(walletResponses, &WalletResponse{
			WalletID:       wallet.WalletID,
			UserID:         wallet.UserID,
			Balance:        wallet.Balance,
			WalletType:     wallet.WalletType,
			ParentWalletID: wallet.ParentWalletID,
			AccessQRCode:   wallet.AccessQRCode,
			AccessNFCCode:  wallet.AccessNFCCode,
			CreatedAt:      wallet.CreatedAt,
		})
	}

	return walletResponses, nil
}

func (s *WalletService) UpdateBalance(ctx context.Context, walletID string, req UpdateBalanceRequest) (*WalletResponse, error) {
	v := validator.New()
	validator.ValidateWalletID(v, walletID)
	validator.ValidateAmount(v, req.Amount)

	if !v.Valid() {
		return nil, v.ToError()
	}

	err := s.models.Wallets.UpdateBalance(walletID, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("failed to update balance: %w", err)
	}

	// Get updated wallet
	wallet, err := s.models.Wallets.GetByID(walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated wallet: %w", err)
	}

	return &WalletResponse{
		WalletID:       wallet.WalletID,
		UserID:         wallet.UserID,
		Balance:        wallet.Balance,
		WalletType:     wallet.WalletType,
		ParentWalletID: wallet.ParentWalletID,
		AccessQRCode:   wallet.AccessQRCode,
		AccessNFCCode:  wallet.AccessNFCCode,
		CreatedAt:      wallet.CreatedAt,
	}, nil
}

func (s *WalletService) DeleteWallet(ctx context.Context, walletID string) error {
	v := validator.New()
	validator.ValidateWalletID(v, walletID)

	if !v.Valid() {
		return v.ToError()
	}

	// Check if wallet exists
	_, err := s.models.Wallets.GetByID(walletID)
	if err != nil {
		return fmt.Errorf("wallet not found: %w", err)
	}

	err = s.models.Wallets.Delete(walletID)
	if err != nil {
		return fmt.Errorf("failed to delete wallet: %w", err)
	}

	return nil
}

func generateQRCode(userID uuid.UUID, walletType string) (string, error) {
	// Generate a temporary wallet ID for the QR code structure
	tempWalletID := uuid.New()

	// Create QR code data structure
	qrData := QRCodeData{
		WalletID:   tempWalletID.String(),
		UserID:     userID.String(),
		WalletType: walletType,
		IssuedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		Version:    "1.0",
	}

	// Generate signature (in production, use proper cryptographic signing)
	signature, err := generateRandomHex(16)
	if err != nil {
		return "", err
	}
	qrData.Signature = signature

	// Convert to JSON and encode
	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal QR data: %w", err)
	}

	// In production, you might want to encrypt this data
	// For now, we'll use base64 encoding of the JSON
	return fmt.Sprintf("SISUPASS_QR_%s", hex.EncodeToString(jsonData)), nil
}

func generateNFCCode(userID uuid.UUID, walletType string) (string, error) {
	// Generate a temporary wallet ID for the NFC code structure
	tempWalletID := uuid.New()

	// Generate access key for NFC
	accessKey, err := generateRandomHex(32)
	if err != nil {
		return "", err
	}

	// Create NFC code data structure
	nfcData := NFCCodeData{
		WalletID:   tempWalletID.String(),
		UserID:     userID.String(),
		WalletType: walletType,
		IssuedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		Version:    "1.0",
		AccessKey:  accessKey,
	}

	// Convert to JSON and encode
	jsonData, err := json.Marshal(nfcData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal NFC data: %w", err)
	}

	// For NFC, we'll use a different prefix and hex encoding
	return fmt.Sprintf("SISUPASS_NFC_%s", hex.EncodeToString(jsonData)), nil
}

func generateRandomHex(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// DecodeQRCode decodes and validates a QR code string
func DecodeQRCode(qrCode string) (*QRCodeData, error) {
	// Remove the prefix
	const prefix = "SISUPASS_QR_"
	if len(qrCode) <= len(prefix) || qrCode[:len(prefix)] != prefix {
		return nil, fmt.Errorf("invalid QR code format")
	}

	// Decode hex data
	hexData := qrCode[len(prefix):]
	jsonData, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode QR code data: %w", err)
	}

	// Unmarshal JSON
	var qrData QRCodeData
	err = json.Unmarshal(jsonData, &qrData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal QR code data: %w", err)
	}

	// Validate expiration
	if time.Now().After(qrData.ExpiresAt) {
		return nil, fmt.Errorf("QR code has expired")
	}

	return &qrData, nil
}

// DecodeNFCCode decodes and validates an NFC code string
func DecodeNFCCode(nfcCode string) (*NFCCodeData, error) {
	// Remove the prefix
	const prefix = "SISUPASS_NFC_"
	if len(nfcCode) <= len(prefix) || nfcCode[:len(prefix)] != prefix {
		return nil, fmt.Errorf("invalid NFC code format")
	}

	// Decode hex data
	hexData := nfcCode[len(prefix):]
	jsonData, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode NFC code data: %w", err)
	}

	// Unmarshal JSON
	var nfcData NFCCodeData
	err = json.Unmarshal(jsonData, &nfcData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal NFC code data: %w", err)
	}

	// Validate expiration
	if time.Now().After(nfcData.ExpiresAt) {
		return nil, fmt.Errorf("NFC code has expired")
	}

	return &nfcData, nil
}

// updateQRCodeWithWalletID updates the wallet ID in an existing QR code
func updateQRCodeWithWalletID(qrCode, walletID string) (string, error) {
	qrData, err := DecodeQRCode(qrCode)
	if err != nil {
		return "", err
	}

	qrData.WalletID = walletID

	// Re-encode with updated wallet ID
	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated QR data: %w", err)
	}

	return fmt.Sprintf("SISUPASS_QR_%s", hex.EncodeToString(jsonData)), nil
}

// updateNFCCodeWithWalletID updates the wallet ID in an existing NFC code
func updateNFCCodeWithWalletID(nfcCode, walletID string) (string, error) {
	nfcData, err := DecodeNFCCode(nfcCode)
	if err != nil {
		return "", err
	}

	nfcData.WalletID = walletID

	// Re-encode with updated wallet ID
	jsonData, err := json.Marshal(nfcData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated NFC data: %w", err)
	}

	return fmt.Sprintf("SISUPASS_NFC_%s", hex.EncodeToString(jsonData)), nil
}
