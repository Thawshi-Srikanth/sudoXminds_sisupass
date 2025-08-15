"use client";

import * as React from "react";
import { cn } from "@/lib/utils";
import { Badge, type VariantProps } from "@/components/ui/badge";
import Image from "next/image";
import { Separator } from "./separator";

interface TicketCardProps extends React.ComponentProps<"div"> {
  cutoutSize?: number;
}

function TicketCard({
  className,
  children,
  cutoutSize = 20,
  ...props
}: TicketCardProps) {
  const clipPathId = React.useId();

  return (
    <>
      <div
        data-slot="ticket-card"
        className={cn(
          "relative bg-card text-card-foreground  rounded-md overflow-hidden min-w-3/4 ",

          className
        )}
        style={{
          clipPath: `url(#${clipPathId})`,
        }}
        {...props}
      >
        {children}
      </div>
    </>
  );
}

interface TicketCardBodyProps extends React.ComponentProps<"div"> {
  imageSrc: string;
  imageAlt?: string;
  badgeText?: string;
  badgeVariant?: VariantProps<typeof Badge>["variant"];
}

function TicketCardBody({
  className,
  imageSrc,
  imageAlt = "Card image",
  badgeText,
  badgeVariant = "default",
  children,
  ...props
}: TicketCardBodyProps) {
  return (
    <div className={cn("relative border-l border-r border-t rounded-lg border-secondary p-2 pb-4", className)} {...props}>
      {/* Image */}

      <span className="absolute bottom-0 left-1/2 -translate-x-1/2 w-3/4 border-1 border-dashed rounded-full" />



      <div className="relative w-full h-48 rounded-md overflow-hidden bg-muted">
        <Image
          src={imageSrc || "/placeholder.svg"}
          alt={imageAlt}
          fill
          className="object-cover"
        />
      </div>

      {/* Badge positioned absolute top-left */}
      {badgeText && (
        <Badge className="absolute top-3 left-3 z-20" variant={badgeVariant}>
          {badgeText}
        </Badge>
      )}

      {children}
    </div>
  );
}

interface TicketCardFooterProps extends React.ComponentProps<"div"> {
  title: string;
  description: string;
  badgeText?: string;
  badgeVariant?: VariantProps<typeof Badge>["variant"];
}

function TicketCardFooter({
  className,
  title,
  description,
  badgeText,
  badgeVariant = "secondary",
  ...props
}: TicketCardFooterProps) {
  return (
    <div
      data-slot="ticket-card-footer"
      className={cn(
        "flex justify-between relative items-center px-6 py-4 bg-muted/30 border-l border-r border-b  rounded-lg border-secondary",
         "before:absolute before:left-0 before:-top-0 before:w-5 before:h-5 before:bg-background before:rounded-full before:-translate-x-1/2 before:-translate-y-1/2 before:z-10 before:border before:border-secondary",
          "after:absolute after:right-0 after:-top-0 after:w-5 after:h-5 after:bg-background after:rounded-full after:translate-x-1/2 after:-translate-y-1/2 after:z-10 after:border after:border-secondary",
        className
      )}
      {...props}
    >
      

      <div className="flex-1">
        <h3 className="font-semibold text-base leading-tight">{title}</h3>
        <p className="text-sm text-muted-foreground mt-1">{description}</p>
      </div>

      {badgeText && (
        <div className="flex flex-col h-full">
          <Badge variant={badgeVariant} className="ml-4">
            {badgeText}
          </Badge>
        </div>
      )}
    </div>
  );
}

// Alternative version with CSS mask for better browser support
interface TicketCardMaskedProps extends React.ComponentProps<"div"> {}

function TicketCardMasked({
  className,
  children,
  ...props
}: TicketCardMaskedProps) {
  return (
    <div
      data-slot="ticket-card-masked"
      className={cn(
        "relative bg-card text-card-foreground shadow-lg border overflow-hidden rounded-lg",
        className
      )}
      style={{
        mask: `
          radial-gradient(circle at 0% 60%, transparent 10px, black 11px),
          radial-gradient(circle at 100% 60%, transparent 10px, black 11px)
        `,
        WebkitMask: `
          radial-gradient(circle at 0% 60%, transparent 10px, black 11px),
          radial-gradient(circle at 100% 60%, transparent 10px, black 11px)
        `,
      }}
      {...props}
    >
      {children}
    </div>
  );
}

export { TicketCard, TicketCardBody, TicketCardFooter, TicketCardMasked };
