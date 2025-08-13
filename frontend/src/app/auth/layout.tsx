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
    <div className="min-h-dvh w-full relative items-center flex flex-col h-full">

      {children}
    </div>

  );
}
