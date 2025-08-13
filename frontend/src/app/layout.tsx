import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { Button } from "@/components/ui/button";
import { BellDot, Menu, UserRound } from "lucide-react";

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
      <body
        className={`${inter.variable} ${jetbrains.variable} antialiased relative max-h-dvh flex flex-col items-center justify-start`}
      >
        <div className="max-w-[640px] grid grid-cols-6 gap-x-4 gap-y-6 p-6 pb-3 w-full bg-background sticky top-0">
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
        <div className="max-w-[640px] grid grid-cols-6 gap-x-4 gap-y-6 p-6 pt-3 w-full max-h-dvh h-full">
          {children}
        </div>
      </body>
    </html>
  );
}
