import { Button } from "@/components/ui/button";
import Link from "next/link";

export default function Home() {
  return (
    <div className="col-span-6 flex flex-col flex-1 h-full justify-around mt-8">
      <div className=" flex flex-col ">
        <h1 className=" text-2xl font-bold text-balance text-center">
          Join the SisuPass Family
        </h1>
        <p className=" text-base/5 tracking-wide text-center">
         Sign up and make school payments simple, fun, and secure.
        </p>
      </div>

      <div className=" flex flex-col gap-3">
        <Button size="xl">Sign Up Email</Button>
        <Button size="xl" variant="outline">
          Continue With Google
        </Button>
      </div>
      <div className=" flex flex-col">
        <p className="py-5 text-center text-base ">
          Have an account??
          <Link href="#" className="ml-2 font-semibold text-secondary">
            Sign In
          </Link>
        </p>
      </div>
    </div>
  );
}
