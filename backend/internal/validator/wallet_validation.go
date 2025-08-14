package validator

import (
	"sisupass.com/sisupass/internal/types"
)

func ValidateWallet(v *Validator, wallet *types.Wallet) {
	v.Check(wallet.UserID.String() != "", "user_id", "must be provided")
	v.Check(wallet.Balance >= 0, "balance", "must be zero or positive")
	v.Check(wallet.WalletType != "", "wallet_type", "must be provided")
	v.Check(len(wallet.WalletType) <= 50, "wallet_type", "must not be more than 50 bytes long")

	// Validate wallet type is one of the allowed values
	allowedWalletTypes := []string{"main", "savings", "child", "temporary"}
	v.Check(PermittedValue(wallet.WalletType, allowedWalletTypes...), "wallet_type", "must be one of: main, savings, child, temporary")
}

func ValidateCreateWalletRequest(v *Validator, walletType string, balance float64) {
	v.Check(walletType != "", "wallet_type", "must be provided")
	v.Check(len(walletType) <= 50, "wallet_type", "must not be more than 50 bytes long")
	v.Check(balance >= 0, "balance", "must be zero or positive")

	// Validate wallet type is one of the allowed values
	allowedWalletTypes := []string{"main", "savings", "child", "temporary"}
	v.Check(PermittedValue(walletType, allowedWalletTypes...), "wallet_type", "must be one of: main, savings, child, temporary")
}

func ValidateWalletID(v *Validator, walletID string) {
	v.Check(walletID != "", "wallet_id", "must be provided")
}

func ValidateAmount(v *Validator, amount float64) {
	v.Check(amount > 0, "amount", "must be greater than zero")
}
