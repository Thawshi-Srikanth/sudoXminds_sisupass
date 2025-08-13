"use client";

import { motion } from "motion/react";
import Image from "next/image";
import { Button } from "@/components/ui/button";
import { useForm } from "react-hook-form";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Form, FormField } from "@/components/ui/form";
import { PasswordInput } from "@/components/ui/password-input";

export default function SigninPage() {
  const form = useForm<{
    username: string;
    password: string;
    confirmPassword: string;
  }>();

  const onSubmit = (data: { username: string; password: string }) => {
    console.log("Login data:", data);
  };

  return (
    <div className="flex flex-col w-full flex-1">
      {/* Top image section */}
      <div className="w-full flex min-h-[350px] bg-secondary relative overflow-hidden">
        <motion.div
          initial={{ opacity: 0, x: 50 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.6 }}
          className="absolute right-2 aspect-[1.5/2] bottom-0 h-[250px] z-10"
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
              Almost There!
            </h1>
            <p className=" text-base tracking-wide">
              Enter your details to create your wallet.
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
                  name="username"
                  render={({ field }) => (
                    <div className="flex flex-col gap-3">
                      <Label>Username</Label>
                      <Input
                        placeholder="Enter username"
                        className="h-12"
                        {...field}
                      />
                    </div>
                  )}
                />

                <FormField
                  control={form.control}
                  name="password"
                  render={({ field }) => (
                    <div className="flex flex-col gap-3">
                      <Label>Password</Label>
                      <PasswordInput
                        className="h-12"
                        placeholder="Enter password"
                        {...field}
                      />
                    </div>
                  )}
                />

                <FormField
                  control={form.control}
                  name="confirmPassword"
                  render={({ field }) => (
                    <div className="flex flex-col gap-3">
                      <Label>Confirm Password</Label>
                      <PasswordInput
                        className="h-12"
                        placeholder="Enter password"
                        {...field}
                      />
                    </div>
                  )}
                />

                <Button type="submit" size="xl" className="w-full mt-5">
                  Sign Up
                </Button>
              </form>
            </Form>
          </div>
        </motion.div>
      </div>
    </div>
  );
}
