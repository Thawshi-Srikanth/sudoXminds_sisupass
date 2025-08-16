import type { Metadata } from "next";
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
