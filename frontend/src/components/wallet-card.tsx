"use client";

import {
  CreditCard,
  CreditCardBack,
  CreditCardBalance,
  CreditCardExpiry,
  CreditCardName,
} from "@/components/ui/credit-card";

interface WalletCardProps {
  wallet?: {
    name: string;
    walletId: string;
    balance: string | number;
    walletType?: string;
  };
  isLoading?: boolean;
}

export function WalletCard({ wallet, isLoading }: WalletCardProps) {
  return (
    <CreditCard>
      <CreditCardBack className="bg-secondary text-background">
        <div className="flex h-full flex-col justify-between p-4">
          <CreditCardName className="flex-1">
            {wallet?.name || "Main Wallet"}
          </CreditCardName>

          <CreditCardBalance className="flex flex-col flex-1">
            <span className="text-base tracking-wide opacity-90">Balance</span>
            {isLoading ? "Loading..." : wallet?.balance ?? "0.00"}
          </CreditCardBalance>

          <div className="flex justify-between gap-4">
            <CreditCardExpiry>
              {wallet?.walletId || "xxxx-xxxx-xxxx-xxxx"}
            </CreditCardExpiry>
          </div>
        </div>
      </CreditCardBack>
    </CreditCard>
  );
}
