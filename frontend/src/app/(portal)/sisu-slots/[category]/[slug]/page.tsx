import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import {
  TicketCard,
  TicketCardBody,
  TicketCardFooter,
} from "@/components/ui/slot-card";
import {
  Calendar,
  Building,
  FileText,
  Users,
  GraduationCap,
  HeartPulse,
  Bus,
  Home,
} from "lucide-react"; // import your icons

import Image from "next/image";

// Categories array
const categories = [
  {
    title: "Events",
    description: "Workshops, seminars, and student programs",
    icon: Calendar,
    main: true,
    slug: "events",
  },
  {
    title: "Facilities",
    description: "Book labs, and sports areas",
    icon: Building,
    main: true,
    slug: "facilities",
  },
  {
    title: "Documents",
    description: "Request IDs, transcripts, and certificates",
    icon: FileText,
    main: true,
    slug: "documents",
  },
  {
    title: "Appointments",
    description: "Meet with counselors and staff",
    icon: Users,
    main: true,
    slug: "appointments",
  },
  {
    title: "Exams",
    description: "Register for tests and practicals",
    icon: GraduationCap,
    main: false,
    slug: "exams",
  },
  {
    title: "Health",
    description: "Checkups, vaccinations, and counseling",
    icon: HeartPulse,
    main: false,
    slug: "health",
  },
  {
    title: "Transport",
    description: "Book buses, campus shuttles, and permits",
    icon: Bus,
    main: false,
    slug: "transport",
  },
  {
    title: "Accommodation",
    description: "Reserve hostel rooms and campus housing",
    icon: Home,
    main: false,
    slug: "accommodation",
  },
];

export default async function Page({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug: slug } = await params;

  // Find category by slug
  const category = categories.find((c) => c.slug === slug);

  if (!category) {
    return <p>Category not found.</p>;
  }

  const Icon = category.icon;

  const isBooked = true;

  return (
    <>
      <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full mb-30 h-full">
        <div className="col-span-6 flex flex-col justify-between sticky top-23">
          <TicketCard>
            <TicketCardBody
              imageSrc="/static/images/event-image.png"
              badgeText="VIP"
              badgeVariant="destructive"
            />
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
          <h1 className="scroll-m-20 text-lg font-bold">About</h1>
          <p className="text-xs tracking-wide">
            Organized by the Sri Lanka Schools’ Chess Association, this
            championship brings together the nation’s most promising young minds
            for a battle of strategy, skill, and focus. It serves as a platform
            to showcase talent, foster sportsmanship, and inspire the next
            generation of chess champions in Sri Lanka.
          </p>
        </div>

        <div className="col-span-6 flex flex-col px-2 justify-between gap-3">
          <h1 className="scroll-m-20 text-lg font-bold">Location</h1>
          <div className="flex items-center w-full border  h-16 p-2 pl-4 rounded-md gap-2">
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
      </div>

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
