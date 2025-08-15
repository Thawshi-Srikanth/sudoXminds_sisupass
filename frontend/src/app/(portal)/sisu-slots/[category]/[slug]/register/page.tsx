"use client";

import { DynamicForm } from "@/components/ui/dynamic-form";
import { Separator } from "@/components/ui/separator";
import { TicketCard, TicketCardFooter } from "@/components/ui/slot-card";

const formConfig = {
  title: "Passport Application Form",
  description:
    "Form for applying for a new passport or renewing an existing one.",
  type: "object",
  sections: [
    {
      title: "Personal Information",
      fields: [
        {
          name: "first_name",
          label: "First Name",
          type: "text",
          required: true,
        },
        {
          name: "last_name",
          label: "Last Name",
          type: "text",
          required: true,
        },
        {
          name: "middle_name",
          label: "Middle Name",
          type: "text",
          required: false,
        },
        {
          name: "date_of_birth",
          label: "Date of Birth",
          type: "date",
          required: true,
        },
        {
          name: "place_of_birth",
          label: "Place of Birth",
          type: "text",
          required: true,
        },
        {
          name: "gender",
          label: "Gender",
          type: "select",
          options: ["Male", "Female", "Other"],
          required: true,
        },
        {
          name: "nationality",
          label: "Nationality",
          type: "text",
          required: true,
        },
      ],
    },
    {
      title: "Contact Information",
      fields: [
        {
          name: "email",
          label: "Email Address",
          type: "email",
          required: true,
        },
        {
          name: "phone_number",
          label: "Phone Number",
          type: "tel",
          required: true,
        },
        {
          name: "address",
          label: "Residential Address",
          type: "textarea",
          required: true,
        },
      ],
    },
    {
      title: "Passport Details",
      fields: [
        {
          name: "passport_type",
          label: "Passport Type",
          type: "select",
          options: ["New", "Renewal", "Replacement"],
          required: true,
        },
        {
          name: "previous_passport_number",
          label: "Previous Passport Number",
          type: "text",
          required: false,
          show_if: {
            passport_type: ["Renewal", "Replacement"],
          },
        },
        {
          name: "issue_date",
          label: "Previous Passport Issue Date",
          type: "date",
          required: false,
          show_if: {
            passport_type: ["Renewal", "Replacement"],
          },
        },
      ],
    },
    {
      title: "Document Upload",
      fields: [
        {
          name: "photo",
          label: "Passport-size Photo",
          type: "file",
          accept: ["image/jpeg", "image/png"],
          required: true,
        },
        {
          name: "identity_proof",
          label: "Proof of Identity",
          type: "file",
          accept: ["application/pdf", "image/jpeg", "image/png"],
          required: true,
        },
      ],
    },
  ],
};

export default function Home() {
  const handleFormSubmit = async (data: Record<string, any>) => {
    console.log("Form submitted:", data);

    // Simulate API call
    try {
      await new Promise((resolve) => setTimeout(resolve, 1000));
      alert("Form submitted successfully! Check console for data.");
    } catch (error) {
      alert("Error submitting form. Please try again.");
    }
  };

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full mb-30 h-full">
      <div className="col-span-6 flex flex-col justify-between sticky top-23">
        <TicketCard>
          <TicketCardFooter
            title="Event Title"
            description="Event details"
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
          className="flex items-center w-full border  h-16 p-2 pl-4 rounded-md "
          style={{ gap: "1rem" }} // optional spacing between flex items
        >
          {/* Column 2: Transaction ID (top), Date (bottom) */}
          <div className="flex flex-col justify-center space-y-1 flex-grow">
            <span className="font-mono text-sm font-semibold">
              Saturday, Nov 04 2025
            </span>
            <span className="text-xs text-muted-foreground">
              10:00 am - 05:00 pm
            </span>
          </div>

          {/* Column 3: Amount (top right), Type (bottom right) */}
          <div className="flex flex-col items-center justify-center rounded-sm space-y-1 h-full aspect-square bg-secondary">
            <span className="text-lg  font-mono text-center uppercase text-background  font-medium">
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
