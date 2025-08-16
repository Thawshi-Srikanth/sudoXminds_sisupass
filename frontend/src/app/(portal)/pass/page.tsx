import { PlusIcon } from "lucide-react";
import { Button } from "@/components/ui/button";

import { ScrollArea } from "@/components/ui/scroll-area";
import WalletList from "./wallets";
import Link from "next/link";

export default function Home() {
  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6  w-full h-full">
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="flex justify-between flex-col">
          <h1 className="scroll-m-20  text-2xl font-bold text-balance">
            Your Passes
          </h1>
          <p className=" text-base tracking-wide">
            Manage all your passes here.
          </p>
        </div>
        <div className="flex ">
          <Link href={"/pass/add-new"}>
            <Button variant="muted" size="icon" className="size-12">
              <PlusIcon width={24} height={24} />
            </Button>
          </Link>
        </div>
      </div>

      <div className="col-span-6 flex justify-between flex-col">
        <ScrollArea className="h-[calc(100vh-200px)] w-full">
          <div className="flex flex-col h-full divide-y w-full  gap-4">
            <WalletList />
          </div>
        </ScrollArea>
      </div>
    </div>
  );
}
