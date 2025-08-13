import type { Metadata } from "next";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";

export const metadata: Metadata = {
  title: "SiSu Pass - Hello there",
  description: "Smart ID, Safe Usage ",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <>
      <div className=" grid grid-cols-6 gap-x-4 gap-y-6 p-6 pb-3 w-full sticky top-0 z-1">
        <div className="col-span-6 flex justify-between">
          <Button variant="ghost" size="icon" className="size-12 text-background">
            <ArrowLeft width={24} height={24} />
          </Button>

          <div></div>
        </div>
      </div>

      {children}
    </>

  );
}
