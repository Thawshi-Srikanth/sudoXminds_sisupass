import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="col-span-6 flex flex-col flex-1 h-full justify-around mt-8">
      <div className=" flex flex-col ">
        <h1 className=" text-2xl font-bold text-balance text-center">
          Tap in to SisuPass
        </h1>
        <p className=" text-base/5 tracking-wide text-center">
          Access your smart wallet. Keep the good times rolling.
        </p>
      </div>

      <div className=" flex flex-col gap-3">
        <Button size="xl">Sign In Email</Button>
        <Button size="xl" variant="outline">
          Sign In Google
        </Button>
      </div>
      <div className=" flex flex-col">
        <p className="py-5 text-center text-base ">
          Don’t have an account?
          <Link href="#" className="ml-2 font-semibold text-secondary">
            Sign Up
          </Link>
        </p>
      </div>
    </div>
  );
}
