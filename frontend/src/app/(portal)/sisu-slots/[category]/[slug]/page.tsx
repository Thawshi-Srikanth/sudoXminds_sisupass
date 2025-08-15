"use client";

import { useParams } from "next/navigation";
import Image from "next/image";
import { gql } from "@apollo/client";
import client from "@/lib/apolloClient";

import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";
import {
  TicketCard,
  TicketCardBody,
  TicketCardFooter,
} from "@/components/ui/slot-card";

import { useEffect, useState } from "react";

const GET_SLOT_BY_ID = gql`
  query SlotById($id: UUID!) {
    slotById(id: $id) {
      id
      title
      coverImage
      description
      action
      slotType {
        name
      }
    }
  }
`;

export default function SlotDetail() {
  const params = useParams();
  const { slug } = params;

  const [slot, setSlot] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<any>(null);

  useEffect(() => {
    const fetchSlot = async () => {
      try {
        const { data } = await client.query({
          query: GET_SLOT_BY_ID,
          variables: { id: slug },
          fetchPolicy: "no-cache",
        });

        if (data.slotById) {
          // Parse JSON strings
          const parsedDesc = JSON.parse(data.slotById.description || "{}");
          const parsedAction = JSON.parse(data.slotById.action || "{}");

          setSlot({
            ...data.slotById,
            description: parsedDesc,
            action: parsedAction,
          });
        }
      } catch (err) {
        setError(err);
      } finally {
        setLoading(false);
      }
    };
    fetchSlot();
  }, [slug]);

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error loading slot</p>;
  if (!slot) return <p>Slot not found</p>;

  const fallbackImage = "/static/images/card-noise.png";
  const isBooked = true;

  const desc = slot.description;
  const startDate = desc?.event?.start_date
    ? new Date(desc.event.start_date)
    : null;
  const endDate = desc?.event?.end_date ? new Date(desc.event.end_date) : null;
  const location = desc?.location || null;

  return (
    <>
      <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 mb-30">
        {/* Header Slot Card */}
        <div className="col-span-6 flex flex-col justify-between sticky top-23">
          <TicketCard>
            <TicketCardBody imageSrc={slot.coverImage || fallbackImage}>
         
            </TicketCardBody>
            <TicketCardFooter
              title={slot.title}
              description={desc?.about || ""}
              badgeText={
                slot.action?.type === "form"
                  ? "Form"
                  : slot.action?.type === "link"
                  ? "Link"
                  : ""
              }
            />
          </TicketCard>
        </div>

        <div className="col-span-6 flex flex-col justify-between gap-3">
          <Separator />
        </div>

        {/* Date & Time */}
        {startDate && endDate && (
          <div className="col-span-6 flex flex-col px-2 gap-3">
            <h1 className="text-lg font-bold">Date & Time</h1>
            <div className="flex items-center w-full border h-16 p-2 pl-4 rounded-md gap-4">
              <div className="flex flex-col justify-center space-y-1 flex-grow">
                <span className="font-mono text-sm font-semibold">
                  {startDate.toLocaleDateString("en-US", {
                    weekday: "long",
                    month: "short",
                    day: "numeric",
                    year: "numeric",
                  })}
                </span>
                <span className="text-xs text-muted-foreground">
                  {desc.event.time || ""}
                </span>
              </div>
              <div className="flex flex-col items-center justify-center rounded-sm space-y-1 h-full aspect-square bg-secondary">
                <span className="text-lg font-mono text-center text-background font-medium">
                  {startDate.getDate()}
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Location */}
        {location && (
          <div className="col-span-6 flex flex-col px-2 gap-3">
            <h1 className="text-lg font-bold">Location</h1>
            <div className="flex items-center w-full border h-16 p-2 pl-4 rounded-md gap-2">
              <div className="flex flex-col justify-center space-y-1 flex-grow">
                <span className="font-mono text-sm font-semibold">
                  {location.address || location}
                </span>
                <span className="text-xs text-muted-foreground">
                  Lat: {location.latitude || ""}, Lng:{" "}
                  {location.longitude || ""}
                </span>
              </div>
            </div>
          </div>
        )}

        {/* Other Description Fields */}
        {desc?.schedules && (
          <div className="col-span-6 flex flex-col gap-2">
            <h2 className="text-lg font-semibold">Schedules</h2>
            <pre className="text-xs bg-gray-50 p-2 rounded-md">
              {JSON.stringify(desc.schedules, null, 2)}
            </pre>
          </div>
        )}

        {desc?.facilities && (
          <div className="col-span-6 flex flex-col gap-2">
            <h2 className="text-lg font-semibold">Facilities</h2>
            <ul className="text-xs list-disc pl-4">
              {desc.facilities.map((f: string, idx: number) => (
                <li key={idx}>{f}</li>
              ))}
            </ul>
          </div>
        )}

        {/* Action Buttons */}
        {slot.action && (
          <div className="col-span-6 mt-4">
            {slot.action.type === "form" && (
              <Button size="lg" variant="secondary">
                Open Form
              </Button>
            )}
            {slot.action.type === "link" && (
              <Button
                size="lg"
                variant="secondary"
                onClick={() => window.open(slot.action.url, "_blank")}
              >
                {slot.action.label || "Go to Link"}
              </Button>
            )}
          </div>
        )}
      </div>

      {/* Bottom Action */}
      <div className="fixed bottom-0 flex p-8 border-t z-50 h-30 items-center justify-center w-full gap-3 bg-background">
        {isBooked ? (
          <>
            <Button className="flex-1" variant="destructive" size="lg">
              Cancel
            </Button>
            <Button className="flex-2" variant="secondary" size="lg">
              View Qr
            </Button>
          </>
        ) : (
          <Button className="w-full" variant="secondary" size="lg">
            Register
          </Button>
        )}
      </div>
    </>
  );
}
