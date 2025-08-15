import {
  Menu,
  BellDot,
  UserRound,
  ArrowUpRight,
  DollarSign,
  Wallet,
  CalendarClock,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  CreditCard,
  CreditCardBack,
  CreditCardBalance,
  CreditCardCvv,
  CreditCardExpiry,
  CreditCardName,
} from "@/components/ui/credit-card";
import { CardStack } from "@/components/ui/card-stack";
import { TransactionList } from "@/components/portal/transaction-list";
import Link from "next/link";



export default function Home() {
  return (
    <>
      <div className="col-span-6 flex justify-between flex-col px-2">
        <h1 className="scroll-m-20  text-2xl font-bold text-balance">
          Hello Thawshi
        </h1>
        <p className=" text-base tracking-wide">Fun for them, peace for you!</p>
      </div>

      <div className="col-span-6 flex justify-between flex-col">
        <CardStack>
          <CreditCard>
            <CreditCardBack className=" bg-primary text-background ">
              <div className="flex flex-col justify-between h-full p-4">
                <CreditCardName className="flex-1">Buss Pass</CreditCardName>
                <CreditCardBalance className="flex flex-col flex-1">
                  <span className=" text-base tracking-wide">Balance</span>
                  Rs. 10,000.00
                </CreditCardBalance>

                <div className="flex justify-between gap-4">
                  <CreditCardExpiry>01/24</CreditCardExpiry>
                  <CreditCardCvv>123</CreditCardCvv>
                </div>
              </div>
            </CreditCardBack>
          </CreditCard>

          <CreditCard>
            <CreditCardBack className=" bg-secondary text-background ">
              <div className="flex flex-col justify-between h-full p-4">
                <CreditCardName className="flex-1">Buss Pass</CreditCardName>
                <CreditCardBalance className="flex flex-col flex-1">
                  <span className=" text-base tracking-wide">Balance</span>
                  Rs. 10,000.00
                </CreditCardBalance>

                <div className="flex justify-between gap-4">
                  <CreditCardExpiry>01/24</CreditCardExpiry>
                  <CreditCardCvv>123</CreditCardCvv>
                </div>
              </div>
            </CreditCardBack>
          </CreditCard>

          <CreditCard>
            <CreditCardBack className=" bg-blue-700 text-background ">
              <div className="flex flex-col justify-between h-full p-4">
                <CreditCardName className="flex-1">Buss Pass</CreditCardName>
                <CreditCardBalance className="flex flex-col flex-1">
                  <span className=" text-base tracking-wide">Balance</span>
                  Rs. 10,000.00
                </CreditCardBalance>

                <div className="flex justify-between gap-4">
                  <CreditCardExpiry>01/24</CreditCardExpiry>
                  <CreditCardCvv>123</CreditCardCvv>
                </div>
              </div>
            </CreditCardBack>
          </CreditCard>
        </CardStack>
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
          <Button variant="secondary" size="icon" className=" w-24 h-12">
            <CalendarClock width={24} height={24} />
          </Button>
          <p className=" text-sm tracking-wide text-center">SiSu Slots</p>
        </div>
      </div>
      <div className="col-span-6 flex flex-1 justify-between flex-col">
        <TransactionList/>
      </div>
    </>
  );
}
