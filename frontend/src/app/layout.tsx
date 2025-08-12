import type { Metadata } from "next";
import { Inter, JetBrains_Mono } from "next/font/google";
import "./globals.css";


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
        className={`${inter.variable} ${jetbrains.variable} antialiased max-h-dvh flex items-center justify-center`}
      >
        <div className="max-w-[640px] grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full">
        {children}
        </div>
      </body>
    </html>
  );
}
