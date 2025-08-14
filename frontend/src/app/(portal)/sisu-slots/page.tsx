"use client";

import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Calendar,
  Building,
  FileText,
  Users,
  GraduationCap,
  HeartPulse,
  Bus,
  Home,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

// Categories with `main: true` for the top section
const categories = [
  {
    title: "Events",
    description: "Workshops, seminars, and student programs",
    icon: Calendar,
    main: true,
  },
  {
    title: "Facilities",
    description: "Book labs, and sports areas",
    icon: Building,
    main: true,
  },
  {
    title: "Documents",
    description: "Request IDs, transcripts, and certificates",
    icon: FileText,
    main: true,
  },
  {
    title: "Appointments",
    description: "Meet with counselors and staff",
    icon: Users,
    main: true,
  },
  {
    title: "Exams",
    description: "Register for tests and practicals",
    icon: GraduationCap,
    main: false,
  },
  {
    title: "Health",
    description: "Checkups, vaccinations, and counseling",
    icon: HeartPulse,
    main: false,
  },
  {
    title: "Transport",
    description: "Book buses, campus shuttles, and permits",
    icon: Bus,
    main: false,
  },
  {
    title: "Accommodation",
    description: "Reserve hostel rooms and campus housing",
    icon: Home,
    main: false,
  },
];

export default function Bookings() {
  const mainCategories = categories.filter((cat) => cat.main);
  const otherCategories = categories.filter((cat) => !cat.main);

  return (
    <>
      {/* Header */}
      <div className="col-span-6 flex flex-col justify-between">
        <h1 className="scroll-m-20 text-2xl font-bold">SiSu Slots</h1>
        <p className="text-base tracking-wide">Fun for them, peace for you!</p>
      </div>

      {/* Main Category Cards - 2 by grid */}
      <div className="col-span-6 grid grid-cols-2 gap-4">
        {mainCategories.map((cat, index) => {
          const Icon = cat.icon;
          return (
            <Card
              key={index}
              className="h-full shadow-none transition-shadow cursor-pointer p-3 rounded-md"
            >
              <CardHeader className="flex flex-col items-start gap-4 p-0">
                <div className="p-2 bg-accent/40 rounded-lg ">
                  <Icon className="w-6 h-6" color="var(--color-secondary)" />
                </div>
                <div>
                  <CardTitle className="text-base font-semibold">
                    {cat.title}
                  </CardTitle>
                  <CardDescription className="text-xs">
                    {cat.description}
                  </CardDescription>
                </div>
              </CardHeader>
            </Card>
          );
        })}
      </div>

      <Separator className="col-span-6" />

      <div className="col-span-6 flex flex-col justify-between">
        <h2 className="text-lg font-semibold mb-2">Other Bookable Slots</h2>
      </div>

      {/* Other Bookable Slots */}
      <div className="col-span-6">
        <div className="flex flex-wrap gap-3">
          {otherCategories.map((cat, index) => {
            const Icon = cat.icon;
            return (
              <Button
                key={index}
                size="lg"
                variant="outline"
                className="shadow-none"
              >
                <Icon className="w-4 h-4" color="var(--color-secondary)" />
                <span className="text-sm font-medium">{cat.title}</span>
              </Button>
            );
          })}
        </div>
      </div>
    </>
  );
}
