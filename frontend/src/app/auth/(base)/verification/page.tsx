"use client";

import { motion } from "motion/react";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Form, FormField } from "@/components/ui/form";
import { PasswordInput } from "@/components/ui/password-input";
import {
  InputOTP,
  InputOTPGroup,
  InputOTPSeparator,
  InputOTPSlot,
} from "@/components/ui/input-otp";

export default function SigninPage() {
  const form = useForm<{
    otp: string;
  }>();

  const onSubmit = (data: { otp: string }) => {
    console.log("Otp data:", data);
  };

  return (
    <div className="flex flex-col w-full flex-1">
      {/* Top image section */}
      <div className="w-full flex min-h-[400px] bg-secondary relative overflow-hidden">
        <motion.div
          initial={{ opacity: 0, x: 50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="absolute right-0 aspect-[2.5/2.9] bottom-0 h-[250px] z-10"
        >
          <Image
            src="/static/images/verification-girl.png"
            alt="verification-girl"
            fill
            priority
            sizes="x2"
          />
        </motion.div>

        <motion.div
          className="absolute top-1/2 left-0 -translate-y-1/2 z-4"
          initial={{ x: "-100%" }}
          animate={{ x: 20 }}
          transition={{
            type: "spring",
            stiffness: 100,
            damping: 20,
            duration: 1.5,
          }}
        >
          <h1 className="scroll-m-20 text-5xl font-bold text-balance text-background">
            It’s hiding
            <br />
            in your
            <br />
            inbox…
          </h1>
        </motion.div>

        <motion.div
          initial={{ opacity: 0, x: -50, rotate: -12 }}
          animate={{ opacity: 0.3, x: 0 }}
          transition={{ type: "spring", stiffness: 300, damping: 30 }}
          className="absolute aspect-[3.5/2] h-[200px] top-1/2 left-1/2 -translate-y-1/2 flex items-end justify-end"
        >
          <Image
            src="/static/images/card-noise.png"
            alt="card-noise"
            fill
            sizes="x1"
          />
        </motion.div>

        <motion.div
          initial={{ opacity: 0, x: 50, rotate: 12 }}
          animate={{ opacity: 0.3, x: 0 }}
          transition={{
            type: "spring",
            stiffness: 300,
            damping: 30,
            delay: 0.1,
          }}
          className="absolute aspect-[3.5/2] h-[200px] top-1/2 left-1/2 -translate-y-1/2 flex items-end justify-end"
        >
          <Image
            src="/static/images/card-noise.png"
            alt="card-noise"
            fill
            sizes="x1"
          />
        </motion.div>
      </div>

      {/* Content section */}
      <div className="flex flex-1 justify-center items-center">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.6, delay: 0.2 }}
          className="grid grid-cols-6 gap-x-4 gap-y-6 p-6 pb-3 w-full"
        >
          <div className="col-span-6 flex justify-between flex-col px-2 w-full h-full">
            <h1 className="scroll-m-20  text-2xl font-bold text-balance">
              Verify your email
            </h1>
            <p className=" text-base tracking-wide">
              Code has been set to you email.
            </p>
          </div>
          <div className="col-span-6 flex justify-between flex-col px-2 w-full h-full">
            <Form {...form}>
              <form
                onSubmit={form.handleSubmit(onSubmit)}
                className="space-y-4"
              >
                <FormField
                  control={form.control}
                  name="otp"
                  render={({ field }) => (
                    <div className="flex flex-col gap-3">
                      <InputOTP className="h-12" {...field} maxLength={6}>
                        <InputOTPGroup>
                          <InputOTPSlot index={0} />
                          <InputOTPSlot index={1} />
                          <InputOTPSlot index={2} />
                        </InputOTPGroup>
                        <InputOTPSeparator />
                        <InputOTPGroup>
                          <InputOTPSlot index={3} />
                          <InputOTPSlot index={4} />
                          <InputOTPSlot index={5} />
                        </InputOTPGroup>
                      </InputOTP>
                    </div>
                  )}
                />

                <p className=" text-base tracking-wide">
                  sent to your email abc*****@email.com
                </p>

                <div className="grid grid-cols-2 justify-between flex-col p-x w-full">
                  <Button
                    type="submit"
                    size="lg"
                    variant="ghost"
                    className="w-full mt-5 p-2 text-left justify-start"
                  >
                    Resend email
                  </Button>
                  <Button type="submit" size="lg" className="w-full mt-5">
                    Verify
                  </Button>
                </div>
              </form>
            </Form>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
