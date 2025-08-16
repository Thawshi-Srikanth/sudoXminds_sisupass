"use client";
import * as React from "react";
import {
  Popover,
  PopoverTrigger,
  PopoverContent,
} from "@/components/ui/popover";
import {
  Command,
  CommandList,
  CommandItem,
  CommandInput,
  CommandEmpty,
  CommandGroup,
} from "@/components/ui/command";
import { Button } from "@/components/ui/button";
import { Avatar, AvatarImage, AvatarFallback } from "@/components/ui/avatar";
import { CheckIcon, PlusIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { AvatarStack } from "@/components/ui/avatar-stack";

export interface AvatarOption {
  label: string;
  value: string;
  image?: string;
}

interface AvatarMultiSelectProps {
  options: AvatarOption[];
  value: string[];
  onChange: (value: string[]) => void;
  size?: number;
  searchable?: boolean;
}

export function AvatarMultiSelect({
  options,
  value,
  onChange,
  size = 40,
  searchable = true,
}: AvatarMultiSelectProps) {
  const [open, setOpen] = React.useState(false);
  const toggleOption = (val: string) => {
    onChange(
      value.includes(val) ? value.filter((v) => v !== val) : [...value, val]
    );
  };

  const getOption = (val: string) => options.find((o) => o.value === val);

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="ghost"
          className="flex items-start gap-0 h-auto p-0 justify-start"
          onClick={() => setOpen(true)}
        >
          {/* Avatar stack */}
          <AvatarStack>
            <span className="inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium transition-all disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none  shrink-0 [&_svg]:shrink-0 outline-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive border bg-accent/50 hover:text-accent-foreground dark:bg-input/30 dark:border-input dark:hover:bg-input/50 size-full">
              <PlusIcon width={24} height={24} />
            </span>
            {value.length > 0 &&
              value.map((val) => {
                const opt = getOption(val);
                if (!opt) return null;
                return (
                  <Avatar
                    key={val}
                    className="border-2 border-background"
                    style={{ width: size, height: size }}
                  >
                    <AvatarImage src={opt.image || ""} />
                    <AvatarFallback>
                      {opt.label
                        .split(" ")
                        .map((n) => n[0])
                        .join("")}
                    </AvatarFallback>
                  </Avatar>
                );
              })}
          </AvatarStack>
        </Button>
      </PopoverTrigger>

      <PopoverContent className=" w-52 p-0">
        <Command>
          {searchable && (
            <CommandInput placeholder="Search users..." className="h-9" />
          )}
          <CommandList>
            <CommandEmpty>No users found.</CommandEmpty>
            <CommandGroup>
              {options.map((opt) => {
                const selected = value.includes(opt.value);
                return (
                  <CommandItem
                    key={opt.value}
                    onSelect={() => toggleOption(opt.value)}
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <div
                      className={cn(
                        "flex h-4 w-4 items-center justify-center rounded-sm border",
                        selected
                          ? "bg-primary text-primary-foreground"
                          : "opacity-50 [&_svg]:invisible"
                      )}
                    >
                      <CheckIcon className="h-4 w-4" />
                    </div>
                    <Avatar
                      className="border"
                      style={{ width: 24, height: 24 }}
                    >
                      <AvatarImage src={opt.image || ""} />
                      <AvatarFallback>
                        {opt.label
                          .split(" ")
                          .map((n) => n[0])
                          .join("")}
                      </AvatarFallback>
                    </Avatar>
                    <span>{opt.label}</span>
                  </CommandItem>
                );
              })}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
}
