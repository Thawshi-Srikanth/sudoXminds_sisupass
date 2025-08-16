"use client";

import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { ArrowLeft, CheckCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ScheduleCalendar } from "@/components/schedule-calendar";
import { QRCode } from "@/components/ui/qr-code";

export default function RegisterPage() {
  const params = useParams();
  const router = useRouter();
  const { category, slug } = params;
  const [bookingSuccess, setBookingSuccess] = useState<any>(null);

  const handleBookingComplete = (booking: any) => {
    setBookingSuccess(booking);
  };

  const handleBackToSlot = () => {
    router.push(`/sisu-slots/${category}/${slug}`);
  };

  if (bookingSuccess) {
    const qrData = {
      booking_id: bookingSuccess.id,
      schedule_id: bookingSuccess.schedule.id,
      slot_id: bookingSuccess.schedule.slot.id,
      wallet_id: bookingSuccess.wallet.walletId,
    };
    return (
      <div className="container mx-auto px-4 py-8">
        <Card className="max-w-md mx-auto text-center border-secondary shadow-none">
          <CardContent className="pt-6">
            <div className="w-full flex item-center justify-center mb-5">
              <QRCode
                className="w-48 h-48 rounded border bg-white p-4"
                data={JSON.stringify(qrData)}
              />
            </div>

            <h2 className="text-2xl font-bold mb-2">Booking Confirmed!</h2>
            <p className="text-muted-foreground mb-4">
              Your booking for &quot;{bookingSuccess.schedule.slot.title}&quot;
              has been confirmed.
            </p>
            <div className="space-y-2">
              <Button onClick={handleBackToSlot} className="w-full">
                Back to Slot
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 mb-30">
      <div className="col-span-6 flex flex-col justify-between gap-3">
        {/* Header */}
        <div className="flex items-center gap-4 mb-6">
          <h1 className="text-2xl font-bold">Slot Registration</h1>
        </div>

        {/* Content */}
        <div className="max-w-2xl mx-auto">
          <ScheduleCalendar
            slotId={slug as string}
            onBookingComplete={handleBookingComplete}
          />
        </div>
      </div>
    </div>
  );
}
