import {
  ArrowUpRight,
  DollarSign,
  Wallet,
  CalendarClock,
} from "lucide-react";
import { Button } from "@/components/ui/button";

import { TransactionList } from "@/components/portal/transaction-list";
import Link from "next/link";
import WalletList from "./wallets";

export default function Home() {
  return (
    <div className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 w-full h-full">
      <div className="col-span-6 flex justify-between flex-col px-2">
        <h1 className="scroll-m-20  text-2xl font-bold text-balance">
          Hello Thawshi
        </h1>
        <p className=" text-base tracking-wide">Fun for them, peace for you!</p>
      </div>

      <div className="col-span-6 flex justify-between flex-col">
        <WalletList />
      </div>
      <div className="col-span-6 flex px-4 justify-between ">
        <div className="flex flex-col gap-2 items-center justify-center">
          <Button variant="muted" size="icon" className=" size-12">
            <ArrowUpRight width={24} height={24} />
          </Button>
          <p className=" text-sm tracking-wide text-center">Send</p>
        </div>
        <div className="flex flex-col gap-2 items-center justify-center">
          <Button variant="muted" size="icon" className=" size-12">
            <DollarSign width={24} height={24} />
          </Button>
          <p className=" text-sm tracking-wide text-center">Top Up</p>
        </div>
        <div className="flex flex-col gap-2 items-center justify-center">
          <Link href="pass">
            <Button variant="muted" size="icon" className=" size-12">
              <Wallet width={24} height={24} />
            </Button>
          </Link>
          <p className=" text-sm tracking-wide text-center">Passes</p>
        </div>
        <div className="flex flex-col gap-2 items-center justify-center">
          <Link href="sisu-slots">
          <Button variant="secondary" size="icon" className=" w-24 h-12">
            <CalendarClock width={24} height={24} />
          </Button>
          </Link>
          <p className=" text-sm tracking-wide text-center">SiSu Slots</p>
        </div>
      </div>
      <div className="col-span-6 flex flex-1 justify-between flex-col">
        <TransactionList limit={5}/>
      </div>
    </div>
  );
}
