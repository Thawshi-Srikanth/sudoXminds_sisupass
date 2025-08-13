import { ExpDatePicker } from "@/components/portal/pass/exp-date-picker";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

export default function Home() {
  return (
    <>
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="flex justify-between flex-col">
          <h1 className="scroll-m-20  text-2xl font-bold text-balance">
            Your Passes
          </h1>
          <p className=" text-base tracking-wide">
            Manage all your passes here.
          </p>
        </div>
      </div>
      <div className="col-span-6 flex justify-between w-full px-2">
        <div className="grid grid-cols-4 w-full gap-3">
          <div className="col-span-4 grid w-full  items-center gap-3">
            <Label htmlFor="email">Name</Label>
            <Input type="name" id="name" placeholder="Name" />
          </div>

          <div className="col-span-2 grid w-full  items-center gap-3">
            <ExpDatePicker />
          </div>
        </div>
      </div>
    </>
  );
}
