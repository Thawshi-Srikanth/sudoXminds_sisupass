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
        <div className="grid grid-cols-4 w-full">
          <div className="col-span-4 grid w-full  items-center gap-3">
            <Label htmlFor="email">Name</Label>
            <Input type="name" id="name" placeholder="Name" />
          </div>

          <div className="col-span-4 grid w-full  items-center gap-3">
            <Label htmlFor="email">Exp</Label>
            <Popover>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          className={cn(
            "w-[240px] justify-start text-left font-normal",
            !date && "text-muted-foreground"
          )}
        >
          <CalendarIcon className="mr-2 h-4 w-4" />
          {date ? format(date, "PPP") : <span>Pick a date</span>}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0" align="start">
        <Calendar mode="single" selected={date} onSelect={setDate} autoFocus />
      </PopoverContent>
    </Popover>
          </div>
        </div>
      </div>
    </>
  );
}
