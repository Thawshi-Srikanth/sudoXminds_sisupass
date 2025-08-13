import { PlusCircle, CreditCard, Send } from "lucide-react";
import React from "react";
import { Button } from "../ui/button";
import { ScrollArea } from "../ui/scroll-area";

interface Transaction {
  id: string;
  date: string;
  amount: number;
  type: "topup" | "spending" | "sending";
  description?: string;
}

interface TransactionListProps {
  transactions: Transaction[];
}

const iconByType = {
  topup: <PlusCircle className="h-6 w-6 text-green-600" />,
  spending: <CreditCard className="h-6 w-6 text-red-600" />,
  sending: <Send className="h-6 w-6 text-yellow-400" />,
};

const bgByType = {
  topup: "bg-green-100",
  spending: "bg-red-100",
  sending: "bg-yellow-100",
};

export const TransactionList = ({ transactions }: TransactionListProps) => {
  return (
    <section>
      <div className="container px-4 mx-auto gap-2 flex flex-col">
        <div className="flex justify-between">
          <p className="text-large font-medium ">Recent Transactions</p>
          <Button variant="link" className=" p-0 text-secondary">
            View all
          </Button>
        </div>

        <ScrollArea className="h-[calc(100vh-600px)] w-full">
          <div className="flex flex-col divide-y w-full  gap-4">
            {transactions.map((tx) => {
              const isExpense = tx.type === "spending";
              const amountSign = isExpense ? "-" : "+";
              const amountColor = isExpense ? "text-red-600" : "text-green-600";

              return (
                <div
                  key={tx.id}
                  className="flex items-center w-full pb-4"
                  style={{ gap: "1rem" }} // optional spacing between flex items
                >
                  {/* Column 1: Icon with background */}
                  <div
                    className={` h-12 w-12 items-center justify-center rounded-md bg-accent/40 hidden sm:flex`}
                  >
                    {iconByType[tx.type]}
                  </div>

                  {/* Column 2: Transaction ID (top), Date (bottom) */}
                  <div className="flex flex-col justify-center space-y-1 flex-grow">
                    <span className="font-mono text-sm font-semibold text-gray-800">
                      {tx.id}
                    </span>
                    <span className="text-xs text-gray-500">{tx.date}</span>
                  </div>

                  {/* Column 3: Amount (top right), Type (bottom right) */}
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
