import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Toaster } from "@/components/ui/sonner";

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
        className={`${inter.variable} ${jetbrains.variable} antialiased max-h-lvh`}
      >
        <ScrollArea className="relative max-h-lvh flex flex-col items-center justify-start">
          {children}
        </ScrollArea>

        <Toaster richColors   />
      </body>
    </html>
  );
}
