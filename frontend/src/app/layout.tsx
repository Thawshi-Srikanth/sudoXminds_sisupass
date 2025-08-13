import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { Button } from "@/components/ui/button";
import { BellDot, Menu, UserRound } from "lucide-react";
import { ScrollArea } from "@/components/ui/scroll-area";

const inter = Inter({
  variable: "--font-sans",
  subsets: ["latin"],
});

const jetbrains = JetBrains_Mono({
  variable: "--font-mono",
  subsets: ["latin"],
});

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
    <html lang="en">
      <body className={`${inter.variable} ${jetbrains.variable} antialiased max-h-lvh overflow-hidden`}>
        <ScrollArea className="relative max-h-lvh flex flex-col items-center justify-start">
          {children}
        </ScrollArea>
      </body>
    </html>
  );
}
