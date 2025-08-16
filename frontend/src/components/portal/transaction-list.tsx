"use client";

import React, { useEffect, useState } from "react";
import { gql } from "@apollo/client";
import client from "@/lib/apolloClient";
import { Button } from "../ui/button";
import { ScrollArea } from "../ui/scroll-area";
import { PlusCircle, CreditCard, Send, Copy } from "lucide-react";
import Link from "next/link";

// --- Types ---
interface Transaction {
  id: string;
  date: string;
  amount: number;
  type: "topup" | "spend" | "transfer";
}

// --- GraphQL Query: only necessary fields ---
const GET_USER_TRANSACTIONS = gql`
  query GetUserTransactions($limit: Int) {
    userTransactions(limit: $limit) {
      transactionId
      transactionDate
      amount
      transactionType
    }
  }
`;

// --- Component ---
export const TransactionList = ({ limit }: { limit?: number }) => {
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const { data } = await client.query({
          query: GET_USER_TRANSACTIONS,
          variables: { limit },
          fetchPolicy: "no-cache",
        });

        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const mapped: Transaction[] = data.userTransactions.map((tx: any) => ({
          id: tx.transactionId,
          date: tx.transactionDate.split("T")[0],
          amount: tx.amount,
          type: tx.transactionType.toLowerCase(), // normalize to lowercase
        }));

        setTransactions(mapped);
        setLoading(false);
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
      } catch (err: any) {
        console.error(err);
        setError("Failed to fetch transactions");
        setLoading(false);
      }
    };

    fetchTransactions();
  }, [limit]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>{error}</p>;

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    alert("Transaction ID copied!");
  };

  const truncate = (text: string, length: number = 8) =>
    text.length > length ? `${text.slice(0, length)}...` : text;

  return (
    <section>
      <div className="container px-4 mx-auto gap-2 flex flex-col">
        <div className="flex justify-between">
          <p className="text-large font-medium">Recent Transactions</p>
          <Link href="/transactions">
            <Button variant="link" className="p-0 text-secondary">
              View all
            </Button>
          </Link>
        </div>

        <ScrollArea className="h-[350px] w-full">
          <div className="flex flex-col divide-y w-full gap-4">
            {transactions.map((tx) => {
              const amountSign = tx.type === "topup" ? "+" : "-";
              const amountColor =
                tx.type === "topup" ? "text-green-600" : "text-red-600";

              return (
                <div
                  key={tx.id}
                  className="flex items-center w-full  h-16"
                  style={{ gap: "1rem" }}
                >
                  {/* Icon */}
                  <div className="h-12 w-12 flex items-center justify-center rounded-md bg-accent/40">
                    {tx.type === "topup" && (
                      <PlusCircle
                        className="h-6 w-6 flex-shrink-0 stroke-emerald-400"
                        color="currentColor"
                      />
                    )}
                    {tx.type === "spend" && (
                      <CreditCard
                        className="h-6 w-6 flex-shrink-0 stroke-red-600"
                        color="currentColor"
                      />
                    )}
                    {tx.type === "transfer" && (
                      <Send className="h-6 w-6 flex-shrink-0 stroke-yellow-500 " />
                    )}
                  </div>

                  {/* Transaction ID + Copy */}
                  <div className="flex flex-col justify-center space-y-1 flex-grow">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-sm font-semibold text-gray-800">
                        {truncate(tx.id)}
                      </span>
                      <Copy
                        className="h-4 w-4 cursor-pointer text-[--color-secondary] flex-shrink-0"
                        onClick={() => copyToClipboard(tx.id)}
                      />
                    </div>
                    <span className="text-xs text-gray-500">{tx.date}</span>
                  </div>

                  {/* Amount + Type */}
                  <div className="flex flex-col items-end justify-center space-y-1">
                    <span
                      className={`font-semibold text-base ${amountColor} font-mono`}
                    >
                      {amountSign}${Math.abs(tx.amount).toFixed(2)}
                    </span>
                    <span className="text-xs uppercase text-gray-500 font-medium">
                      {tx.type}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        </ScrollArea>
      </div>
    </section>
  );
};
