"use client";
import {
  CreditCard,
  CreditCardBack,
  CreditCardBalance,
  CreditCardCvv,
  CreditCardExpiry,
  CreditCardName,
} from "@/components/ui/credit-card";
import { Skeleton } from "@/components/ui/skeleton";
import { gql, useQuery } from "@apollo/client";
import client from "@/lib/apolloClient";

type Wallet = {
  walletId: string;
  walletType?: string;
  balance: number;
  createdAt?: string;
};

const GET_WALLETS = gql`
  query {
    myPassWallets {
      walletId
      balance
      createdAt
    }
  }
`;

export default function WalletList() {
  const { data, loading, error } = useQuery<{ myPassWallets: Wallet[] }>(
    GET_WALLETS,
    {
      client,
    }
  );

  if (error) {
    return (
      <div className="text-sm text-destructive">
        something went wrong loading your wallets.
      </div>
    );
  }

  // skeletons while loading
  if (loading) {
    return (
      <>
        {[...Array(3)].map((_, i) => (
          <Skeleton className=" flex w-full h-[200px] rounded-2xl" key={i} />
        ))}
      </>
    );
  }

  const wallets = data?.myPassWallets ?? [];

  if (wallets.length === 0) {
    return <div className="text-sm text-muted-foreground">no wallets yet.</div>;
  }

  return (
    <>
      {wallets.map((wallet) => (
        <CreditCard key={wallet.walletId}>
          <CreditCardBack className="bg-secondary text-background">
            <div className="flex h-full flex-col justify-between p-4">
              <CreditCardName className="flex-1">Buss Pass</CreditCardName>

              <CreditCardBalance className="flex flex-col flex-1">
                <span className="text-base tracking-wide opacity-90">
                  Balance
                </span>
                {wallet.balance}
              </CreditCardBalance>

              <div className="flex justify-between gap-4">
                <CreditCardExpiry>01/24</CreditCardExpiry>
                <CreditCardCvv>123</CreditCardCvv>
              </div>
            </div>
          </CreditCardBack>
        </CreditCard>
      ))}
    </>
  );
}
