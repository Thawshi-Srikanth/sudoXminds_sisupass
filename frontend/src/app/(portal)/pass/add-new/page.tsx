"use client";

import { useState } from "react";
import { ExpDatePicker } from "@/components/portal/pass/exp-date-picker";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { AvatarMultiSelect } from "@/components/portal/pass/avatar-select";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import PassCard from "./pass-card";

const users = [
  { label: "Thawshi Krish", value: "TK" },
  { label: "Thawshi Srikanth", value: "TS" },
  { label: "Jane Doe", value: "JD" },
];

export default function Home() {
  const [selected, setSelected] = useState<string[]>([]);

  return (
    <>
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="flex justify-between flex-col">
          <h1 className="scroll-m-20  text-2xl font-bold text-balance">
            Add New Pass
          </h1>
          <p className=" text-base tracking-wide">
            Fill out the form to create a new pass.
          </p>
        </div>
      </div>
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="grid grid-cols-4 w-full gap-3">
          <div className="col-span-4 grid w-full  items-center gap-3">
            <Label htmlFor="email">Name</Label>
            <Input type="name" id="name" placeholder="Name" />
          </div>

          <div className="col-span-2 grid w-full  items-center gap-3">
            <ExpDatePicker />
          </div>

          <div className="col-span-2 grid w-full  items-center gap-3">
            <Label htmlFor="pass-holder">Pass Holder</Label>
            <AvatarMultiSelect
              options={users}
              value={selected}
              onChange={setSelected}
            />
          </div>
        </div>
      </div>

      <div className="col-span-6 flex flex-1 h-full  w-full ">
        <PassCard />
      </div>
    </>
  );
}
