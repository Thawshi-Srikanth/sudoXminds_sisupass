import type { Metadata } from "next";
import { Button } from "@/components/ui/button";
import { BellDot, Menu, UserRound } from "lucide-react";

export const metadata: Metadata = {
  title: "SiSu Pass",
  description: "Smart ID, Safe Usage",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <div className="max-w-[640px] grid grid-cols-6 gap-x-4 gap-y-6 p-6 pb-3 w-full bg-background fixed top-0 z-[9999]">
        <div className="col-span-6 flex justify-between">
          <Button variant="ghost" size="icon" className="size-12 ">
            <Menu width={24} height={24} />
          </Button>

          <div>
            <Button variant="ghost" size="icon" className="size-12 ">
              <BellDot width={24} height={24} />
            </Button>
            <Button variant="ghost" size="icon" className="size-12 ">
              <UserRound width={24} height={24} />
            </Button>
          </div>
        </div>
      </div>
      <div className="max-w-[640px]  pt-20 ">
        {children}
      </div>
    </>
  );
}
