import {
  PlusIcon,
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
import { ScrollArea } from "@/components/ui/scroll-area";


export default function Home() {
  return (
    <>
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="flex justify-between flex-col">
          <h1 className="scroll-m-20  text-2xl font-bold text-balance">
          Your Passes
        </h1>
        <p className=" text-base tracking-wide">Manage all your passes here.</p>
          </div>
           <div className="col-span-6 flex px-4 justify-between ">
        <Button variant="muted" size="icon" className="size-12" >
          <PlusIcon width={24} height={24} />
        </Button>
      </div>
        
      </div>

      <div className="col-span-6 flex justify-between flex-col">
        <ScrollArea className="h-[calc(100vh-200px)] w-full">
          <div className="flex flex-col h-full divide-y w-full  gap-4">
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
          </div>
        </ScrollArea>
      </div>
     
    </>
  );
}
