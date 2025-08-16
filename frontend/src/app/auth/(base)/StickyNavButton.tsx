"use client";

import { usePathname } from "next/navigation";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { cn } from "@/lib/utils";
import Link from "next/link";

export function StickyNavButton() {
  const pathname = usePathname();
  const isVerificationPage = pathname === "/auth/sign-up/user-verification";

  return (
    <div
      className={cn(
        " top-0 z-[100] w-full transition-colors duration-300",
        isVerificationPage ? "bg-background sticky" : " bg-transparent fixed"
      )}
    >
      <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 pb-3 w-full">
        <div className="col-span-6 flex justify-between relative">
          <Link href={"/auth"}>
            <Button
              variant="secondary"
              size="icon"
              className="size-12 text-background shadow-none"
            >
              <ArrowLeft width={24} height={24} color="background" />
            </Button>
          </Link>
          <div />
        </div>
      </div>
    </div>
  );
}
