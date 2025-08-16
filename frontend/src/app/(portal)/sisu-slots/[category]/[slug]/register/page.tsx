"use client";

import { useQuery, gql } from "@apollo/client";
import client from "@/lib/apolloClient";
import { DynamicForm } from "@/components/ui/dynamic-form";
import { Separator } from "@/components/ui/separator";
import { TicketCard, TicketCardFooter } from "@/components/ui/slot-card";
import { useState, useEffect } from "react";
import { useParams } from "next/navigation";

const GET_SLOT_BY_ID = gql`
  query GetSlotById($id: UUID!) {
    slotById(id: $id) {
      id
      title
      description
      fields
    }
  }
`;

export default function Registration() {
  const params = useParams();
  const { slug } = params;

  const { data, loading, error } = useQuery(GET_SLOT_BY_ID, {
    client,
    variables: { id: slug },
  });

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const [formConfig, setFormConfig] = useState<any>(null);
  const [slotInfo, setSlotInfo] = useState<{
    title: string;
    description: string;
  }>({ title: "", description: "" });

  useEffect(() => {
    if (data?.slotById) {
      // Parse description for date & location if it's stored as JSON
      let descriptionText = "";
      if (data.slotById.description) {
        try {
          const desc =
            typeof data.slotById.description === "string"
              ? JSON.parse(data.slotById.description)
              : data.slotById.description;

          descriptionText = `${
            desc.date ? new Date(desc.date).toDateString() : ""
          } - ${desc.location || ""}`;
        } catch {
          descriptionText = data.slotById.description;
        }
      }

      setSlotInfo({
        title: data.slotById.title,
        description: descriptionText,
      });

      setFormConfig({
        title: data.slotById.title,
        description: descriptionText,
        type: "object",
        sections: JSON.parse(data.slotById.fields) || [],
      });
    }
  }, [data]);

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const handleFormSubmit = async (formData: Record<string, any>) => {
    console.log("Form submitted:", formData);
    alert("Form submitted successfully!");
  };

  if (loading) return <p>Loading slot...</p>;
  if (error) return <p>Error loading slot: {error.message}</p>;
  if (!formConfig) return <p>No form configuration found</p>;

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full mb-30 h-full">
      <div className="col-span-6 flex flex-col justify-between sticky top-23">
        <TicketCard>
          <TicketCardFooter
            title={slotInfo.title}
            description={slotInfo.description}
            badgeText="$99"
          />
        </TicketCard>
      </div>

      <div className="col-span-6 flex flex-col justify-between gap-3">
        <Separator />
      </div>

      <div className="col-span-6 flex flex-col px-2 justify-between gap-3">
        <h1 className="scroll-m-20 text-lg font-bold">Date & Time</h1>
        <div
          className="flex items-center w-full border h-16 p-2 pl-4 rounded-md"
          style={{ gap: "1rem" }}
        >
          <div className="flex flex-col justify-center space-y-1 flex-grow">
            <span className="font-mono text-sm font-semibold">
              {slotInfo.description.split(" - ")[0]}
            </span>
            <span className="text-xs text-muted-foreground">
              {slotInfo.description.split(" - ")[1]}
            </span>
          </div>

          <div className="flex flex-col items-center justify-center rounded-sm space-y-1 h-full aspect-square bg-secondary">
            <span className="text-lg font-mono text-center uppercase text-background font-medium">
              08
            </span>
          </div>
        </div>
      </div>

      <div className="col-span-6 flex flex-col px-2 justify-between gap-3">
        <DynamicForm config={formConfig} onSubmit={handleFormSubmit} />
      </div>
    </div>
  );
}
