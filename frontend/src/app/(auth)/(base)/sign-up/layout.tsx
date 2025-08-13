import type { Metadata } from "next";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import Image from "next/image";

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
      <div className=" w-full flex min-h-[400px] max-h-1/2 h-full bg-secondary relative">
        <div className=" absolute left-10 aspect-[2/3.3] bottom-0 h-[300px] flex items-end justify-end z-3">
          <Image
            src="/static/images/young-boy-sign-up.png"
            alt="young-boy-sign-up"
            fill
            priority
            sizes="x2"
          />
        </div>
        <div className="absolute aspect-[3.5/2] h-[200px] -rotate-12 top-1/2 left-1/2  -translate-y-1/2 flex items-end justify-end z-2">
          <Image
            src="/static/images/detailed-sisu-pass.png"
            alt="detailed-sisu-pass"
            fill
          />
        </div>
        <div className="absolute aspect-[3.5/2] h-[200px] rotate-12 top-1/2 left-1/2 -translate-y-1/2 flex items-end justify-end">
          <Image
            src="/static/images/card-noise.png"
            alt="detailed-sisu-pass"
            fill
          />
        </div>
      </div>
      <div className="flex flex-1">
        <div className="max-w-[640px] grid grid-cols-6 gap-x-4 gap-y-6 p-6 pt-3 w-full relative min-h-full ">
          {children}
        </div>
      </div>
    </>
  );
}
