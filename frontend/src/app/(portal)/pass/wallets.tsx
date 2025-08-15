"use client";
import {
  CreditCard,
  CreditCardBack,
  CreditCardBalance,
  CreditCardCvv,
  CreditCardExpiry,
  CreditCardName,
} from "@/components/ui/credit-card";
import { gql, useQuery } from "@apollo/client";
import client from "@/lib/apolloClient";

type Wallet = {
  walletId: string;
  walletType: string;
  balance: number;
  parentWallet?: Wallet | null;
  user: {
    id: string;
    email: string;
  };
};

const GET_WALLETS = gql`
  query {
    myWallets {
      walletId
      balance
      createdAt
    }
  }
`;

export default function WalletList() {
  const { data, loading } = useQuery<{ myWallets: Wallet[] }>(GET_WALLETS, {
    client,
  });
  return (
    <>
      {data?.myWallets.map((wallet) => (
        <CreditCard key={wallet.walletId}>
          <CreditCardBack className=" bg-primary text-background ">
            <div className="flex flex-col justify-between h-full p-4">
              <CreditCardName className="flex-1">Buss Pass</CreditCardName>
              <CreditCardBalance className="flex flex-col flex-1 text-background">
                <span className=" text-base tracking-wide ">Balance</span>
                {wallet.balance}
              </CreditCardBalance>

              <div className="flex justify-between gap-4 text-background">
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
