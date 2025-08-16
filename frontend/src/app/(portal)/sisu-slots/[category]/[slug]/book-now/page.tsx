"use client";

import { useParams, useRouter } from "next/navigation";
import { useState } from "react";
import { ArrowLeft, CheckCircle } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { ScheduleCalendar } from "@/components/schedule-calendar";

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
    return (
      <div className="container mx-auto px-4 py-8">
        <Card className="max-w-md mx-auto text-center">
          <CardContent className="pt-6">
            <CheckCircle className="h-16 w-16 text-green-500 mx-auto mb-4" />
            <h2 className="text-2xl font-bold mb-2">Booking Confirmed!</h2>
            <p className="text-muted-foreground mb-4">
              Your booking for "{bookingSuccess.schedule.slot.title}" has been
              confirmed.
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
