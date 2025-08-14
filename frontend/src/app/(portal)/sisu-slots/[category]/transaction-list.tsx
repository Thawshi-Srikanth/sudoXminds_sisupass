import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import React from "react";

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


export const TransactionList = ({ transactions }: TransactionListProps) => {
  return (
    <section>
      <div className="container px-2 mx-auto gap-2 flex flex-col">
        <div className="flex justify-between">
          <p className="text-large font-medium ">Up Coming Events</p>
          <Button variant="link" className=" p-0 text-secondary">
            View all
          </Button>
        </div>

        <ScrollArea className="h-[calc(100vh-600px)] w-full">
          <div className="flex flex-col  w-full  gap-4">
            {transactions.map((tx) => {
              return (
                <div
                  key={tx.id}
                  className="flex items-center w-full bg-accent/50  h-16 p-2 pl-4 rounded-md "
                  style={{ gap: "1rem" }} // optional spacing between flex items
                >
      

                  {/* Column 2: Transaction ID (top), Date (bottom) */}
                  <div className="flex flex-col justify-center space-y-1 flex-grow">
                    <span className="font-mono text-sm font-semibold">
                      {tx.id}
                    </span>
                    <span className="text-xs text-muted-foreground">{tx.date}</span>
                  </div>

                  {/* Column 3: Amount (top right), Type (bottom right) */}
                  <div className="flex flex-col items-center justify-center rounded-sm space-y-1 h-full aspect-square bg-secondary">
                    <span className="text-lg  font-mono text-center uppercase text-background  font-medium">
                      08
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
