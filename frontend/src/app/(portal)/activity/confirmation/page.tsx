"use client";

import { useSearchParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { CheckCircle } from "lucide-react";
import Link from "next/link";

export default function ConfirmationPage() {
  const searchParams = useSearchParams();
  const transactionId = searchParams.get("transactionId");
  const amount = searchParams.get("amount");
  const newBalance = searchParams.get("newBalance");
  const type = searchParams.get("type") || "topup";

  return (
    <div className="container mx-auto p-6 max-w-2xl">
      <div className="space-y-6">
        <Card className="text-center">
          <CardHeader>
            <div className="mx-auto w-12 h-12 bg-green-100 rounded-full flex items-center justify-center mb-4">
              <CheckCircle className="h-6 w-6 text-green-600" />
            </div>
            <CardTitle className="text-2xl">
              {type === "transfer"
                ? "Transfer Successful!"
                : "Top-up Successful!"}
            </CardTitle>
            <CardDescription>
              Your transaction has been completed successfully
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <div className="flex justify-between">
                <span className="text-muted-foreground">Transaction ID:</span>
                <span className="font-mono text-sm">{transactionId}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-muted-foreground">Amount:</span>
                <span className="font-semibold">${amount}</span>
              </div>
              {newBalance && (
                <div className="flex justify-between">
                  <span className="text-muted-foreground">
                    {type === "transfer"
                      ? "Remaining Balance:"
                      : "New Balance:"}
                  </span>
                  <span className="font-semibold">${newBalance}</span>
                </div>
              )}
            </div>

            <div className="pt-4 space-y-2">
              <Link href="/">
                <Button className="w-full">Back to Home</Button>
              </Link>
              <Link href={type === "transfer" ? "/activity/send" : "/activity/topup"}>
                <Button variant="outline" className="w-full bg-transparent">
                  Make Another {type === "transfer" ? "Transfer" : "Top-up"}
                </Button>
              </Link>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
