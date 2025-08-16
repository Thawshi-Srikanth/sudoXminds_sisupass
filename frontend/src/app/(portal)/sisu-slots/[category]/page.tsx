"use client";

import { useParams } from "next/navigation";
import Image from "next/image";
import { gql, useQuery } from "@apollo/client";
import client from "@/lib/apolloClient";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { UpcomingList } from "./upcoming-list";
import { ScrollArea, ScrollBar } from "@/components/ui/scroll-area";
import {
  TicketCard,
  TicketCardBody,
  TicketCardFooter,
} from "@/components/ui/slot-card";
import { Input } from "@/components/ui/input";
import Link from "next/link";

// GraphQL queries matching your backend resolvers
const GET_SLOTS = gql`
  query SlotsPage($typeId: UUID!, $search: String) {
    slotTypes {
      id
      name
    }
    slotsByType(typeId: $typeId, search: $search) {
      id
      title
      coverImage
      description
      action
    }
    trendingSlots(typeId: $typeId, limit: 5) {
      id
      title
      coverImage
      description
      action
    }
    upcomingBookings {
      id
      bookingDate
      status
      slotName
      schedule {
        startTime
      }
    }
  }
`;

export default function Slots() {
  const params = useParams();
  const { category } = params; // This will be the UUID slug for the type

  const [search, setSearch] = useState("");

  // Apollo client-side query
  const { data, loading, error } = useQuery(GET_SLOTS, {
    client,
    variables: { typeId: category, search },
    fetchPolicy: "no-cache",
  });

  if (loading) return <p>Loading...</p>;
  if (error) return <p>Error loading slots.</p>;

  const categoryName =
    data?.slotTypes.find((t: any) => t.id === category)?.name || category;
  const trending = data?.trendingSlots || [];
  const slots = data?.slotsByType || [];
  const upcoming = data?.upcomingBookings || [];

  const fallbackImage = "/static/images/card-noise.png";

  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full h-full">
      {/* Category Name */}
      <div className="col-span-6 flex flex-col justify-between">
        <h1 className="text-2xl font-bold">{categoryName}</h1>
      </div>

      {/* Search Bar */}
      <div className="col-span-6 w-full flex flex-col justify-between">
        <Input
          placeholder="Search event..."
          className="h-12 border-secondary shadow-none"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {/* Trending Slots */}
      <div className="col-span-6 flex flex-col justify-between">
        <div className="flex justify-between">
          <p className="text-large font-medium">🔥Trending</p>
          <Button variant="link" className="p-0 text-secondary">
            View all
          </Button>
        </div>
        <ScrollArea className="w-full whitespace-nowrap">
          <div className="flex w-max space-x-4">
            {trending.map((slot: any) => (
              <Link
                key={slot.id}
                href={`/sisu-slots/${category}/${slot.id}`}
                className="flex w-[300px]"
              >
                <TicketCard className="max-w-[300px]">
                  <TicketCardBody
                    imageSrc={slot.coverImage || fallbackImage}
                    badgeText="VIP"
                    badgeVariant="destructive"
                  ></TicketCardBody>
                  <TicketCardFooter
                    title={slot.title}
                    description={slot.description?.about || ""}
                    badgeText="$99"
                  />
                </TicketCard>
              </Link>
            ))}
          </div>
          <ScrollBar orientation="horizontal" />
        </ScrollArea>
      </div>

      {/* Upcoming Bookings */}
      <div className="col-span-6 flex flex-1 justify-between flex-col mt-6">
        <UpcomingList upcoming={upcoming} />
      </div>
    </div>
  );
}
