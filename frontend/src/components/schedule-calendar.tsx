/* eslint-disable @typescript-eslint/no-explicit-any */
"use client";

import { useState } from "react";
import { Calendar } from "@/components/ui/calendar";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Clock, DollarSign } from "lucide-react";
import { gql, useQuery, useMutation } from "@apollo/client";
import client from "@/lib/apolloClient";
import {
  MiniCalendar,
  MiniCalendarDay,
  MiniCalendarDays,
  MiniCalendarNavigation,
} from "./ui/mini-calendar";
import { set } from "date-fns";
import { Dialog, DialogContent, DialogHeader, DialogTitle } from "./ui/dialog";
import { DynamicForm } from "./ui/dynamic-form";

const GET_SLOT_DETAILS = gql`
  query GetSlotDetails($id: UUID!) {
    slotById(id: $id) {
      id
      title
      description
      fields
    }
  }
`;

const GET_AVAILABLE_SCHEDULES = gql`
  query GetAvailableSchedules($slotId: ID!) {
    slotSchedules(slotId: $slotId) {
      id
      startTime
      durationMinutes
      flexible
      price
      available
    }
  }
`;

const CREATE_BOOKING = gql`
  mutation CreateBooking(
    $scheduleId: ID!
    $desiredTime: DateTime
    $userWalletId: ID
    $details: JSONString
  ) {
    createBooking(
      scheduleId: $scheduleId
      desiredTime: $desiredTime
      userWalletId: $userWalletId
      details: $details
    ) {
      booking {
        id
        status
        schedule {
          id
          startTime
          price
          slot {
            title
            owner {
              email
            }
          }
        }
        wallet {
          walletId
          balance
          walletType
        }
        paymentTransaction {
          transactionId
          amount
          status
        }
        details
      }
    }
  }
`;

interface ScheduleCalendarProps {
  slotId: string;
  onBookingComplete?: (booking: any) => void;
}

export function ScheduleCalendar({
  slotId,
  onBookingComplete,
}: ScheduleCalendarProps) {
  const [selectedDate, setSelectedDate] = useState<Date | undefined>(
    new Date()
  );
  const [selectedSchedule, setSelectedSchedule] = useState<any>(null);
  const [selectedTime, setSelectedTime] = useState<string>("");
  const [showFormDialog, setShowFormDialog] = useState(false);
  const [formConfig, setFormConfig] = useState<any>(null);

  const {
    data: schedulesData,
    loading,
    refetch,
  } = useQuery(GET_AVAILABLE_SCHEDULES, {
    variables: { slotId },
    fetchPolicy: "cache-and-network",
    client,
  });

  const { data: slotData } = useQuery(GET_SLOT_DETAILS, {
    variables: { id: slotId },
    fetchPolicy: "cache-and-network",
    client,
  });

  const [createBooking, { loading: bookingLoading }] = useMutation(
    CREATE_BOOKING,
    {
      onCompleted: (data) => {
        if (data.createBooking.booking) {
          onBookingComplete?.(data.createBooking.booking);
          refetch(); // Refresh schedules after booking
          setShowFormDialog(false); // Close form dialog
        }
      },
      client,
    }
  );
  const availableSchedules =
    schedulesData?.slotSchedules?.filter(
      (schedule: any) => schedule.available
    ) || [];

  const handleBookingClick = () => {
    if (slotData?.slotById?.fields) {
      try {
        const fields = JSON.parse(slotData.slotById.fields);
        if (fields && fields.length > 0) {
          setFormConfig({
            title: slotData.slotById.title,
            description: "Please fill out the required information",
            sections: fields,
          });
          setShowFormDialog(true);
          return;
        }
      } catch (error) {
        console.error("Error parsing slot fields:", error);
      }
    }

    handleBooking();
  };

  const handleFormSubmit = (formData: Record<string, any>) => {
    handleBooking(formData);
  };

  const handleBooking = async (formData?: Record<string, any>) => {
    if (!selectedSchedule) return;

    const variables: any = {
      scheduleId: selectedSchedule.id,
    };

    if (selectedSchedule.flexible && selectedDate && selectedTime) {
      const [hours, minutes] = selectedTime.split(":");
      const desiredDateTime = new Date(selectedDate);
      desiredDateTime.setHours(
        Number.parseInt(hours),
        Number.parseInt(minutes)
      );
      variables.desiredTime = desiredDateTime.toISOString();
    }

    if (formData) {
      variables.details = JSON.stringify(formData);
    }

    try {
      await createBooking({ variables });
    } catch (error) {
      console.error("Booking failed:", error);
    }
  };

  const formatTime = (dateString: string) => {
    return new Date(dateString).toLocaleTimeString("en-US", {
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const formatDuration = (minutes: number) => {
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return hours > 0 ? `${hours}h ${mins}m` : `${mins}m`;
  };

  return (
    <div className="space-y-6">
      {/* Calendar */}

      <MiniCalendar value={selectedDate} onValueChange={setSelectedDate}>
        <MiniCalendarNavigation direction="prev" />
        <MiniCalendarDays>
          {(date) => <MiniCalendarDay date={date} key={date.toISOString()} />}
        </MiniCalendarDays>
        <MiniCalendarNavigation direction="next" />
      </MiniCalendar>
      {/* <Card>
        <CardHeader>
          <CardTitle>Select Date</CardTitle>
        </CardHeader>
        <CardContent>
          

          <Calendar
            mode="single"
            selected={selectedDate}
            onSelect={setSelectedDate}
            disabled={(date) => date < new Date()}
            className="rounded-md border"
          />
        </CardContent>
      </Card> */}

      {/* Available Schedules */}
      <Card className="shadow-none rounded-mds">
        <CardHeader>
          <CardTitle>Available Time Slots</CardTitle>
        </CardHeader>
        <CardContent>
          {loading ? (
            <p className="text-muted-foreground">Loading schedules...</p>
          ) : availableSchedules.length === 0 ? (
            <p className="text-muted-foreground">
              No available schedules for this slot.
            </p>
          ) : (
            <div className="space-y-3">
              {availableSchedules.map((schedule: any) => (
                <div
                  key={schedule.id}
                  className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                    selectedSchedule?.id === schedule.id
                      ? "border-primary bg-primary/5"
                      : "hover:border-primary/50"
                  }`}
                  onClick={() => setSelectedSchedule(schedule)}
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <Clock className="h-4 w-4 text-muted-foreground" />
                      <div>
                        <p className="font-medium">
                          {schedule.flexible
                            ? "Flexible Time"
                            : formatTime(schedule.startTime)}
                        </p>
                        <p className="text-sm text-muted-foreground">
                          Duration: {formatDuration(schedule.durationMinutes)}
                        </p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge
                        variant="secondary"
                        className="flex items-center gap-1"
                      >
                        <DollarSign className="h-3 w-3" />
                        {schedule.price}
                      </Badge>
                      {schedule.flexible && (
                        <Badge variant="outline">Flexible</Badge>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </CardContent>
      </Card>

      {/* Time Selection for Flexible Schedules */}
      {selectedSchedule?.flexible && (
        <Card>
          <CardHeader>
            <CardTitle>Select Preferred Time</CardTitle>
          </CardHeader>
          <CardContent>
            <input
              type="time"
              value={selectedTime}
              onChange={(e) => setSelectedTime(e.target.value)}
              className="w-full p-2 border rounded-md"
            />
          </CardContent>
        </Card>
      )}

      {/* Booking Button */}
      {selectedSchedule && (
        <Button
          onClick={handleBookingClick}
          disabled={
            bookingLoading || (selectedSchedule.flexible && !selectedTime)
          }
          className="w-full"
          size="lg"
        >
          {bookingLoading
            ? "Processing..."
            : `Book for $${selectedSchedule.price}`}
        </Button>
      )}

      <Dialog open={showFormDialog} onOpenChange={setShowFormDialog}>
        <DialogContent className="max-w-2xl max-h-[80vh] overflow-y-auto">
          {formConfig && (
            <DynamicForm
              config={formConfig}
              onSubmit={handleFormSubmit}
              submitButtonText={`Complete Booking for Rs.${selectedSchedule?.price}`}
            />
          )}
        </DialogContent>
      </Dialog>
    </div>
  );
}
