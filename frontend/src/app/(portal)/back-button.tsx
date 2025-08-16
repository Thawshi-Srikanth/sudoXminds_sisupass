"use client";

import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { useRouter, usePathname } from "next/navigation";

export default function BackButton() {
  const router = useRouter();
  const pathname = usePathname();

  const getBackPath = () => {
    if (pathname === "/pass") return "/";
    if (pathname === "/transactions") return "/";
    if (pathname === "/pass/add-new") return "/pass";

    if (pathname === "/sisu-slots") return "/";
    if (/^\/sisu-slots\/[^/]+$/.test(pathname)) return "/sisu-slots";

    if (/^\/sisu-slots\/[^/]+\/[^/]+$/.test(pathname)) {
      const segments = pathname.split("/");
      return `/${segments.slice(1, 3).join("/")}`; // /sisu-slots/[category]
    }

    if (/^\/sisu-slots\/[^/]+\/[^/]+\/register$/.test(pathname)) {
      const segments = pathname.split("/");
      return `/${segments.slice(1, 4).join("/")}`; // /sisu-slots/[category]/[slug]
    }

    return "/"; // fallback
  };

  const showBackButton = pathname !== "/";

  if (!showBackButton) return <div></div>;

  return (
    <Button
      variant="ghost"
      size="icon"
      className="size-12"
      onClick={() => router.push(getBackPath())}
    >
      <ArrowLeft width={24} height={24} />
    </Button>
  );
}
