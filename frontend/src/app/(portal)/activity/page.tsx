"use client";

import Link from "next/link";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ArrowUpCircle, ArrowRightCircle } from "lucide-react";

export default function HomePage() {
  return (
    <div className="container mx-auto p-6 max-w-2xl">
      <div className="space-y-6">
        <div className="t space-y-2">
          <h1 className="text-3xl font-bold">Wallet Management</h1>
          <p className="text-muted-foreground">
            Manage your wallet transactions with ease
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <ArrowUpCircle className="h-5 w-5 text-green-600" />
                Top Up
              </CardTitle>
              <CardDescription>Add funds to your main wallet</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/activity/topup">
                <Button className="w-full">Go to Top Up</Button>
              </Link>
            </CardContent>
          </Card>

          <Card className="hover:shadow-lg transition-shadow">
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <ArrowRightCircle className="h-5 w-5 text-blue-600" />
                Send Money
              </CardTitle>
              <CardDescription>Transfer funds between wallets</CardDescription>
            </CardHeader>
            <CardContent>
              <Link href="/activity/send">
                <Button className="w-full bg-transparent" variant="outline">
                  Go to Send
                </Button>
              </Link>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
