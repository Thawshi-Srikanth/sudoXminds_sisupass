"use client";

import Link from "next/link";
import { gql, useQuery } from "@apollo/client";
import client from "@/lib/apolloClient";

import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Separator } from "@/components/ui/separator";

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

// --- GraphQL Query ---
const GET_SLOT_TYPES = gql`
  query {
    slotTypes(trending: true) {
      id
      name
      frequency
    }
  }
`;

// --- Icon Mapping ---
const iconMap = {
  Events: Calendar,
  Facilities: Building,
  Documents: FileText,
  Appointments: Users,
  Exams: GraduationCap,
  Health: HeartPulse,
  Transport: Bus,
  Accommodation: Home,
};

export default function Slots() {
  const { data, loading, error } = useQuery(GET_SLOT_TYPES, { client });

  if (loading) {
    return <p className="p-6">Loading slot types...</p>;
  }

  if (error) {
    return <p className="p-6 text-red-500">Error loading slot types</p>;
  }

  const mainCategories = data.slotTypes.filter((cat) => cat.frequency >= 8);
  const otherCategories = data.slotTypes.filter((cat) => cat.frequency < 8);

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full h-full">
      {/* Header */}
      <div className="col-span-6 flex flex-col justify-between">
        <h1 className="scroll-m-20 text-2xl font-bold">SiSu Slots</h1>
        <p className="text-base tracking-wide">Fun for them, peace for you!</p>
      </div>

      {/* Main Category Cards */}
      <div className="col-span-6 grid grid-cols-2 gap-4">
        {mainCategories.map((cat) => {
          const Icon = iconMap[cat.name] || Calendar;
          return (
            <Link
              key={cat.id}
              href={`/sisu-slots/${cat.id}`} // or `/slots/${slug}` if you have slug
              passHref
              className="block"
            >
              <Card className="h-full shadow-none transition-shadow cursor-pointer p-3 rounded-md hover:shadow-md">
                <CardHeader className="flex flex-col items-start gap-4 p-0">
                  <div className="p-2 bg-accent/40 rounded-lg">
                    <Icon className="w-6 h-6" color="var(--color-secondary)" />
                  </div>
                  <div>
                    <CardTitle className="text-base font-semibold">
                      {cat.name}
                    </CardTitle>
                    <CardDescription className="text-xs">
                      Book now and enjoy benefits
                    </CardDescription>
                  </div>
                </CardHeader>
              </Card>
            </Link>
          );
        })}
      </div>

      <Separator className="col-span-6" />

      {/* Other Slots */}
      <div className="col-span-6 flex flex-col justify-between">
        <h2 className="text-lg font-semibold mb-2">Other Bookable Slots</h2>
      </div>

      <div className="col-span-6">
        <div className="flex flex-wrap gap-3">
          {otherCategories.map((cat) => {
            const Icon = iconMap[cat.name] || Calendar;
            return (
              <Link key={cat.id} href={`/sisu-slots/${cat.id}`} passHref>
                <Button
                  size="lg"
                  variant="outline"
                  className="shadow-none flex items-center gap-2"
                >
                  <Icon className="w-4 h-4" color="var(--color-secondary)" />
                  <span className="text-sm font-medium">{cat.name}</span>
                </Button>
              </Link>
            );
          })}
        </div>
      </div>
    </div>
  );
}
