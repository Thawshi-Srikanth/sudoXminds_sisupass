import { Button } from "@/components/ui/button";
import { TransactionList } from "./transaction-list";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
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

const sampleTransactions = [
  {
    id: "TXN12345678",
    date: "2025-08-12",
    amount: 120.0,
    type: "topup",
  },
  {
    id: "TXN12345679",
    date: "2025-08-11",
    amount: 45.0,
    type: "spending",
  },
  {
    id: "TXN12345680",
    date: "2025-08-10",
    amount: 60.0,
    type: "sending",
  },

  {
    id: "TXN12345678s",
    date: "2025-08-12",
    amount: 120.0,
    type: "topup",
  },
  {
    id: "TXN12345679s",
    date: "2025-08-11",
    amount: 45.0,
    type: "spending",
  },
  {
    id: "TXN12345680s",
    date: "2025-08-10",
    amount: 60.0,
    type: "sending",
  },
];


export default async function Page({
  params,
}: {
  params: Promise<{ category: string }>;
}) {
  const { category: slug } = await params;

  // Find category by slug
  const category = categories.find((c) => c.slug === slug);

  if (!category) {
    return <p>Category not found.</p>;
  }

  const Icon = category.icon;

  return (
    <>
      <div className="col-span-6 flex flex-col justify-between">
        <h1 className="scroll-m-20 text-2xl font-bold">{category.title}</h1>
      </div>

      <div className="col-span-6 flex flex-col justify-between">
        <div className="flex justify-between">
          <p className="text-large font-medium ">🔥Trending</p>
          <Button variant="link" className=" p-0 text-secondary">
            View all
          </Button>
        </div>
        <ScrollArea className="w-full whitespace-nowrap">
          <div className="flex w-max space-x-4">
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
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </div>

      <div className="col-span-6 flex flex-1 justify-between flex-col">
        <TransactionList transactions={sampleTransactions} />
      </div>
    </>
  );
}
