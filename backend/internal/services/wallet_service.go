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

type QRCodeData struct {
	WalletID   string    `json:"wallet_id"`
	UserID     string    `json:"user_id"`
	WalletType string    `json:"wallet_type"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Version    string    `json:"version"`
	Signature  string    `json:"signature"`
}

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

	_, err := s.models.Users.GetByID(userID.String())
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

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

	updatedQRCode, err := updateQRCodeWithWalletID(qrCode, wallet.WalletID.String())
	if err != nil {
	} else {
		wallet.AccessQRCode = &updatedQRCode
	}

	updatedNFCCode, err := updateNFCCodeWithWalletID(nfcCode, wallet.WalletID.String())
	if err != nil {
	} else {
		wallet.AccessNFCCode = &updatedNFCCode
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
	tempWalletID := uuid.New()

	qrData := QRCodeData{
		WalletID:   tempWalletID.String(),
		UserID:     userID.String(),
		WalletType: walletType,
		IssuedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		Version:    "1.0",
	}

	signature, err := generateRandomHex(16)
	if err != nil {
		return "", err
	}
	qrData.Signature = signature

	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal QR data: %w", err)
	}

	return fmt.Sprintf("SISUPASS_QR_%s", hex.EncodeToString(jsonData)), nil
}

func generateNFCCode(userID uuid.UUID, walletType string) (string, error) {
	tempWalletID := uuid.New()

	accessKey, err := generateRandomHex(32)
	if err != nil {
		return "", err
	}

	nfcData := NFCCodeData{
		WalletID:   tempWalletID.String(),
		UserID:     userID.String(),
		WalletType: walletType,
		IssuedAt:   time.Now(),
		ExpiresAt:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity
		Version:    "1.0",
		AccessKey:  accessKey,
	}

	jsonData, err := json.Marshal(nfcData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal NFC data: %w", err)
	}

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

func DecodeQRCode(qrCode string) (*QRCodeData, error) {
	const prefix = "SISUPASS_QR_"
	if len(qrCode) <= len(prefix) || qrCode[:len(prefix)] != prefix {
		return nil, fmt.Errorf("invalid QR code format")
	}

	hexData := qrCode[len(prefix):]
	jsonData, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode QR code data: %w", err)
	}

	var qrData QRCodeData
	err = json.Unmarshal(jsonData, &qrData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal QR code data: %w", err)
	}

	if time.Now().After(qrData.ExpiresAt) {
		return nil, fmt.Errorf("QR code has expired")
	}

	return &qrData, nil
}

func DecodeNFCCode(nfcCode string) (*NFCCodeData, error) {
	const prefix = "SISUPASS_NFC_"
	if len(nfcCode) <= len(prefix) || nfcCode[:len(prefix)] != prefix {
		return nil, fmt.Errorf("invalid NFC code format")
	}

	hexData := nfcCode[len(prefix):]
	jsonData, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, fmt.Errorf("failed to decode NFC code data: %w", err)
	}

	var nfcData NFCCodeData
	err = json.Unmarshal(jsonData, &nfcData)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal NFC code data: %w", err)
	}

	if time.Now().After(nfcData.ExpiresAt) {
		return nil, fmt.Errorf("NFC code has expired")
	}

	return &nfcData, nil
}

func updateQRCodeWithWalletID(qrCode, walletID string) (string, error) {
	qrData, err := DecodeQRCode(qrCode)
	if err != nil {
		return "", err
	}

	qrData.WalletID = walletID

	jsonData, err := json.Marshal(qrData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated QR data: %w", err)
	}

	return fmt.Sprintf("SISUPASS_QR_%s", hex.EncodeToString(jsonData)), nil
}

func updateNFCCodeWithWalletID(nfcCode, walletID string) (string, error) {
	nfcData, err := DecodeNFCCode(nfcCode)
	if err != nil {
		return "", err
	}

	nfcData.WalletID = walletID

	jsonData, err := json.Marshal(nfcData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal updated NFC data: %w", err)
	}

	return fmt.Sprintf("SISUPASS_NFC_%s", hex.EncodeToString(jsonData)), nil
}
