"use client";

import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import React, { useState } from "react";
import { QRCode } from "@/components/ui/qr-code";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface Upcoming {
  id: string;
  date: string;
  slotName: string;
  startDate?: string;
  schedule?: {
    startTime: string;
    id: string;
    slot?: {
      id: string;
    };
  };
  wallet?: {
    walletId: string;
  };
}

interface UpcomingListProps {
  upcoming: Upcoming[];
}

export const UpcomingList = ({ upcoming }: UpcomingListProps) => {
  const [selectedTx, setSelectedTx] = useState<Upcoming | null>(null);

  return (
    <section>
      <div className="container px-2 mx-auto gap-2 mb-5 flex flex-col">
        <div className="flex justify-between">
          <p className="text-large font-medium">Up Coming Events</p>
          <Button variant="link" className="p-0 text-secondary">
            View all
          </Button>
        </div>

        <ScrollArea className="h-max w-full">
          <div className="flex flex-col w-full gap-4">
            {upcoming.map((tx) => {
              const startDate = tx.schedule?.startTime
                ? new Date(tx.schedule.startTime)
                : null;

              return (
                <div
                  key={tx.id}
                  className="flex items-center w-full bg-accent/50 h-16 p-2 pl-4 rounded-md cursor-pointer"
                  style={{ gap: "1rem" }}
                  onClick={() => setSelectedTx(tx)}
                >
                  <div className="flex flex-col justify-center space-y-1 flex-grow">
                    <span className="font-mono text-sm font-semibold">
                      {tx.slotName}
                    </span>
                    <span className="text-xs text-muted-foreground">
                      {startDate &&
                        startDate.toLocaleDateString("en-US", {
                          weekday: "long",
                          month: "short",
                          day: "numeric",
                          year: "numeric",
                        })}
                    </span>
                  </div>

                  <div className="flex flex-col items-center justify-center rounded-sm space-y-1 h-full aspect-square bg-secondary">
                    <span className="text-lg font-mono text-center uppercase text-background font-medium">
                      {startDate && startDate.getDate()}
                    </span>
                  </div>
                </div>
              );
            })}
          </div>
        </ScrollArea>
      </div>

      {/* Dialog */}
      {selectedTx && (
        <Dialog open={!!selectedTx} onOpenChange={() => setSelectedTx(null)}>
          <DialogContent>
            <DialogHeader>
              <DialogTitle>Event Details</DialogTitle>
            </DialogHeader>
            <div className="space-y-4">
              <p>
                <strong>Slot Name:</strong> {selectedTx.slotName}
              </p>
              <p>
                <strong>Date:</strong>{" "}
                {selectedTx.schedule?.startTime &&
                  new Date(selectedTx.schedule.startTime).toLocaleString()}
              </p>
              <p>
                <strong>Wallet ID:</strong>{" "}
                {selectedTx.wallet?.walletId || "N/A"}
              </p>

              {/* QR Code */}
              {selectedTx.schedule?.slot?.id &&
                selectedTx.wallet?.walletId && (
                  <div className="flex justify-center mt-4">
                    <QRCode
                      className="w-48 h-48 rounded border bg-white p-4"
                      data={JSON.stringify({
                        booking_id: selectedTx.id,
                        schedule_id: selectedTx.schedule.id,
                        slot_id: selectedTx.schedule.slot.id,
                        wallet_id: selectedTx.wallet.walletId,
                      })}
                    />
                  </div>
                )}

              <Button onClick={() => setSelectedTx(null)}>Close</Button>
            </div>
          </DialogContent>
        </Dialog>
      )}
    </section>
  );
};
