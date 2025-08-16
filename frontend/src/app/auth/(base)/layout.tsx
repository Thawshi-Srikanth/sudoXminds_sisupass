import type { Metadata } from "next";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { StickyNavButton } from "./StickyNavButton";

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
      <StickyNavButton />

      {children}
    </>
  );
}
