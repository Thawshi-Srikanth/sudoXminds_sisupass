"use client";

import { useEffect, useState } from "react";
import { gql } from "@apollo/client";
import client from "@/lib/apolloClient";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Tabs, TabsList, TabsTrigger, TabsContent } from "@/components/ui/tabs";
import { Skeleton } from "@/components/ui/skeleton";
import { PlusCircle, CreditCard, Send, Copy } from "lucide-react";
import { Radar, RadarChart, PolarGrid, PolarAngleAxis } from "recharts";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";

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
  const [filter, setFilter] = useState<"all" | "spend" | "topup">("all");
  const [tabLoading, setTabLoading] = useState(false);

  useEffect(() => {
    const fetchTransactions = async () => {
      try {
        const { data } = await client.query({
          query: GET_USER_TRANSACTIONS,
          variables: { limit },
          fetchPolicy: "no-cache",
        });

        const mapped: Transaction[] = data.userTransactions.map((tx: any) => ({
          id: tx.transactionId,
          date: tx.transactionDate.split("T")[0],
          amount: tx.amount,
          type: tx.transactionType.toLowerCase(), // normalize to lowercase
        }));

        setTransactions(mapped);
        setLoading(false);
      } catch (err: any) {
        console.error(err);
        setError("Failed to fetch transactions");
        setLoading(false);
      }
    };

    fetchTransactions();
  }, [limit]);

  const TransactionSkeleton = () => (
    <div className="flex items-center w-full h-16" style={{ gap: "1rem" }}>
      <Skeleton className="h-12 w-12 rounded-md" />
      <div className="flex flex-col justify-center space-y-1 flex-grow">
        <div className="flex items-center gap-2">
          <Skeleton className="h-4 w-20" />
          <Skeleton className="h-4 w-4" />
        </div>
        <Skeleton className="h-3 w-16" />
      </div>
      <div className="flex flex-col items-end justify-center space-y-1">
        <Skeleton className="h-4 w-16" />
        <Skeleton className="h-3 w-12" />
      </div>
    </div>
  );

  const handleTabChange = async (value: string) => {
    setTabLoading(true);
    await new Promise((resolve) => setTimeout(resolve, 300));
    setFilter(value as "all" | "spend" | "topup");
    setTabLoading(false);
  };

  const transactionStats = {
    topup: transactions.filter((tx) => tx.type === "topup"),
    spend: transactions.filter((tx) => tx.type === "spend"),
    transfer: transactions.filter((tx) => tx.type === "transfer"),
  };

  const radarData = [
    {
      type: "Top ups",
      count: transactionStats.topup.length,
      amount: transactionStats.topup.reduce((sum, tx) => sum + tx.amount, 0),
    },
    {
      type: "Spend",
      count: transactionStats.spend.length,
      amount: transactionStats.spend.reduce(
        (sum, tx) => sum + Math.abs(tx.amount),
        0
      ),
    },
    {
      type: "Transfer",
      count: transactionStats.transfer.length,
      amount: transactionStats.transfer.reduce(
        (sum, tx) => sum + Math.abs(tx.amount),
        0
      ),
    },
  ];

  const chartConfig = {
    count: {
      label: "Count",
      color: "hsl(var(--chart-1))",
    },
    amount: {
      label: "Amount",
      color: "hsl(var(--chart-2))",
    },
  };

  const filteredTransactions = transactions.filter((tx) => {
    if (filter === "all") return true;
    if (filter === "spend")
      return tx.type === "spend" || tx.type === "transfer";
    if (filter === "topup") return tx.type === "topup";
    return true;
  });

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    alert("Transaction ID copied!");
  };

  const truncate = (text: string, length = 8) =>
    text.length > length ? `${text.slice(0, length)}...` : text;

  if (loading) {
    return (
      <section>
        <div className="container px-4 mx-auto gap-2 flex flex-col">
          <Card className="mb-6 border-0 shadow-none p-2">
            <CardContent>
              <Skeleton className="mx-auto aspect-square max-h-[250px] w-full" />
            </CardContent>
          </Card>

          <div className="flex justify-between">
            <Tabs value="all" className="w-full">
              <TabsList className="grid w-full grid-cols-3">
                <TabsTrigger value="all">All</TabsTrigger>
                <TabsTrigger value="spend">Spend</TabsTrigger>
                <TabsTrigger value="topup">Top ups</TabsTrigger>
              </TabsList>
              <TabsContent value="all" className="mt-4">
                <ScrollArea className="h-[350px] w-full">
                  <div className="flex flex-col divide-y w-full gap-4">
                    {Array.from({ length: 5 }).map((_, index) => (
                      <TransactionSkeleton key={index} />
                    ))}
                  </div>
                </ScrollArea>
              </TabsContent>
            </Tabs>
          </div>
        </div>
      </section>
    );
  }

  if (error) return <p>{error}</p>;

  return (
    <section>
      <div className="container px-4 mx-auto gap-2 flex flex-col">
        <Card className=" px-2 py-2 shadow-none outline-0 border-0">
          <CardContent className="px-0">
            <ChartContainer
              config={chartConfig}
              className="mx-auto w-full aspect-square" 
            >
              <RadarChart data={radarData}>
                <ChartTooltip
                  cursor={false}
                  content={<ChartTooltipContent />}
                />
                <PolarAngleAxis dataKey="type" />
                <PolarGrid />
                <Radar
                  dataKey="count"
                  fill="var--color-chart-1)"
                  fillOpacity={0.6}
                  dot={{
                    r: 4,
                    fillOpacity: 1,
                  }}
                />
                <Radar
                  dataKey="amount"
                  fill="var(--color-chart-1)"
                  fillOpacity={0.3}
                  dot={{
                    r: 4,
                    fillOpacity: 1,
                  }}
                />
              </RadarChart>
            </ChartContainer>
          </CardContent>
        </Card>

        <div className="flex justify-between">
          <Tabs
            value={filter}
            onValueChange={handleTabChange}
            className="w-full"
          >
            <TabsList className="grid w-full grid-cols-3">
              <TabsTrigger value="all">All</TabsTrigger>
              <TabsTrigger value="spend">Spend</TabsTrigger>
              <TabsTrigger value="topup">Top ups</TabsTrigger>
            </TabsList>
            <TabsContent value={filter} className="mt-4">
              <ScrollArea className="h-[350px] w-full">
                <div className="flex flex-col divide-y w-full gap-4">
                  {tabLoading ? (
                    Array.from({ length: 3 }).map((_, index) => (
                      <TransactionSkeleton key={`skeleton-${index}`} />
                    ))
                  ) : filteredTransactions.length === 0 ? (
                    <div className="flex flex-col items-center justify-center h-[300px] border-2 border-dashed border-gray-300 rounded-lg">
                      <div className="text-center space-y-2">
                        <div className="text-gray-400 text-lg">📭</div>
                        <p className="text-gray-500 font-medium">
                          No items found
                        </p>
                        <p className="text-gray-400 text-sm">
                          No transactions match the current filter
                        </p>
                      </div>
                    </div>
                  ) : (
                    filteredTransactions.map((tx) => {
                      const isExpense =
                        tx.type === "spend" || tx.type === "transfer";
                      const amountSign = tx.type === "topup" ? "+" : "-";
                      const amountColor =
                        tx.type === "topup" ? "text-green-600" : "text-red-600";

                      return (
                        <div
                          key={tx.id}
                          className="flex items-center w-full  h-16"
                          style={{ gap: "1rem" }}
                        >
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

                          <div className="flex flex-col justify-center space-y-1 flex-grow">
                            <div className="flex items-center gap-2">
                              <span className="font-mono text-sm font-semibold text-gray-800">
                                {truncate(tx.id)}
                              </span>
                              <Copy
                                className="h-4 w-4 cursor-pointer text-(--color-secondary) flex-shrink-0"
                                onClick={() => copyToClipboard(tx.id)}
                              />
                            </div>
                            <span className="text-xs text-gray-500">
                              {tx.date}
                            </span>
                          </div>

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
                    })
                  )}
                </div>
              </ScrollArea>
            </TabsContent>
          </Tabs>
        </div>
      </div>
    </section>
  );
};
