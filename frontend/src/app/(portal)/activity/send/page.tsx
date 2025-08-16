/* eslint-disable @typescript-eslint/no-explicit-any */
"use client";

import { gql, useMutation, useQuery } from "@apollo/client";
import client from "@/lib/apolloClient";
import {
  CreditCard,
  CreditCardBack,
  CreditCardBalance,
  CreditCardExpiry,
  CreditCardName,
} from "@/components/ui/credit-card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useState } from "react";
import { useRouter } from "next/navigation";
import { toast } from "sonner";

const MAIN_WALLET_QUERY = gql`
  query {
    myMainWallets {
      walletId
      balance
      name
      walletType
    }
  }
`;

const GET_PASS_WALLETS = gql`
  query {
    myPassWallets {
      name
      walletId
      balance
      createdAt
      walletType
      expDate
    }
  }
`;

const SEND_MONEY_MUTATION = gql`
  mutation SendMoney($toWalletId: ID, $amount: Decimal!, $note: String) {
    createTransaction(
      toWalletId: $toWalletId
      amount: $amount
      transactionType: "transfer"
      description: $note
    ) {
      transaction {
        transactionId
        amount
        status
        fromWallet {
          walletId
          balance
        }
        toWallet {
          walletId
          balance
        }
      }
    }
  }
`;

export default function SendPage() {
  const [toWalletId, setToWalletId] = useState("");
  const [amount, setAmount] = useState("");
  const [note, setNote] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const {
    data: mainWalletData,
    loading: walletLoading,
    refetch,
  } = useQuery(MAIN_WALLET_QUERY, {
    client,
  });

  const { data: passWalletsData, loading: passWalletsLoading } = useQuery(
    GET_PASS_WALLETS,
    {
      client,
    }
  );

  const [sendMoney] = useMutation(SEND_MONEY_MUTATION, {
    client,
  });

  const handleSend = async () => {
    if (!toWalletId) {
      toast.error("Please enter destination wallet");
      return;
    }
    if (!amount || Number.parseFloat(amount) <= 0) {
      toast.error("Amount is required");
      return;
    }

    setIsLoading(true);

    try {
      const { data } = await sendMoney({
        variables: {
          toWalletId,
          amount: Number.parseFloat(amount),
          note: note || null,
        },
      });

      const tx = data?.createTransaction?.transaction;

      if (tx?.status === "COMPLETED") {
        toast.success("Transfer successful!");
        await refetch();
        router.push(
          `/activity/confirmation?transactionId=${tx.transactionId}&amount=${tx.amount}&newBalance=${tx.fromWallet.balance}&type=transfer`
        );
      } else {
        toast.error("Transfer failed. Please try again.");
      }
    } catch (err: any) {
      toast.error(err.message || "Something went wrong");
    } finally {
      setIsLoading(false);
    }
  };

  const wallet = mainWalletData?.myMainWallets;
  const passWallets = passWalletsData?.myPassWallets || [];

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full h-full">
      <div className="col-span-6 flex justify-between flex-col px-2">
        <h1 className="text-2xl font-bold">Send Money</h1>
        <p className="text-base tracking-wide">
          Transfer funds between wallets
        </p>
      </div>

      <div className="col-span-6">
        <CreditCard>
          <CreditCardBack className="bg-secondary text-background">
            <div className="flex h-full flex-col justify-between p-4">
              <CreditCardName className="flex-1">
                {wallet?.name || "Main Wallet"}
              </CreditCardName>

              <CreditCardBalance className="flex flex-col flex-1">
                <span className="text-base tracking-wide opacity-90">
                  Balance
                </span>
                {walletLoading ? "Loading..." : wallet?.balance ?? "0.00"}
              </CreditCardBalance>

              <div className="flex justify-between gap-4">
                <CreditCardExpiry>
                  {wallet?.walletId || "xxxx-xxxx-xxxx-xxxx"}
                </CreditCardExpiry>
              </div>
            </div>
          </CreditCardBack>
        </CreditCard>
      </div>

      <div className="col-span-6">
        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="toWallet">To Wallet *</Label>
            <Select value={toWalletId} onValueChange={setToWalletId}>
              <SelectTrigger className=" w-full">
                <SelectValue placeholder="Select to wallet" />
              </SelectTrigger>
              <SelectContent>
                {passWalletsLoading ? (
                  <SelectItem value="loading" disabled>
                    Loading wallets...
                  </SelectItem>
                ) : (
                  passWallets.map((wallet: any) => (
                    <SelectItem key={wallet.walletId} value={wallet.walletId}>
                      {wallet.name} - Balance: {wallet.balance}
                    </SelectItem>
                  ))
                )}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="amount">Amount *</Label>
            <Input
              id="amount"
              type="number"
              placeholder="Enter amount"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              min="0"
              step="0.01"
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="note">Note (Optional)</Label>
            <Textarea
              id="note"
              placeholder="Add a note for this transfer"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              rows={3}
            />
          </div>
        </div>

        <Button
          onClick={handleSend}
          disabled={isLoading || !toWalletId || !amount}
          className="w-full mt-4"
          size="lg"
        >
          {isLoading ? "Processing..." : "Send Money"}
        </Button>
      </div>
    </div>
  );
}
