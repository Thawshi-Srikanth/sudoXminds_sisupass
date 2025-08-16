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

const TOP_UP_MUTATION = gql`
  mutation TopUpWallet($amount: Decimal!, $note: String) {
    createTransaction(
      amount: $amount
      transactionType: "topup"
      description: $note
    ) {
      transaction {
        transactionId
        amount
        status
        toWallet {
          walletId
          balance
        }
      }
    }
  }
`;

export default function TopUpPage() {
  const [amount, setAmount] = useState("");
  const [note, setNote] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const {
    data,
    loading: walletLoading,
    refetch,
  } = useQuery(MAIN_WALLET_QUERY, {
    client,
  });

  const [topUpWallet] = useMutation(TOP_UP_MUTATION, {
    client,
  });

  const handleTopUp = async () => {
    if (!amount || Number.parseFloat(amount) <= 0) {
      toast.error("Amount is required");
      return;
    }

    setIsLoading(true);

    try {
      const { data } = await topUpWallet({
        variables: {
          amount: Number.parseFloat(amount),
          note: note || null,
        },
      });

      const tx = data?.createTransaction?.transaction;

       if (tx?.status === "COMPLETED") {
        toast.success("Top-up successful!");
        await refetch();
        router.push(
          `/activity/confirmation?transactionId=${tx.transactionId}&amount=${tx.amount}&newBalance=${tx.toWallet.balance}`
        );
      } else {
        toast.error("Top-up failed. Please try again.");
      }
    } catch (err: any) {
      toast.error(err.message || "Something went wrong");
    } finally {
      setIsLoading(false);
    }
  };

  const wallet = data?.myMainWallets;

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full h-full">
      <div className="col-span-6 flex justify-between flex-col px-2">
        <h1 className="text-2xl font-bold">Top up</h1>
        <p className="text-base tracking-wide">Add funds to your main wallet</p>
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
              placeholder="Add a note for this top-up"
              value={note}
              onChange={(e) => setNote(e.target.value)}
              rows={3}
            />
          </div>
        </div>

        <Button
          onClick={handleTopUp}
          disabled={isLoading || !amount}
          className="w-full mt-4"
          size="lg"
        >
          {isLoading ? "Processing..." : "Top Up"}
        </Button>
      </div>
    </div>
  );
}
