"use client";

import { motion, AnimatePresence } from "motion/react";
import Image from "next/image";
import { useState } from "react";
import { Button } from "@/components/ui/button";
import axios from "axios";

import GoogleSignInButton from "./google-button";
import { GoogleOAuthProvider } from "@react-oauth/google";

export default function AuthPage() {
  const [mode, setMode] = useState<"signin" | "signup">("signin");
  const isSignIn = mode === "signin";

  const GOOGLE_AUTH_CLIENT_ID =  "167431025530-ardvei2ptlk4gd4u6r81k7mp263f7p1j.apps.googleusercontent.com";
  const API_URL = "http://localhost:8000";


  const handleGoogleSignIn = async (response) => {
    try {
      const { data } = await axios.post(`${API_URL}/auth/google/`, {
        access_token: response.access_token,
      });
      // Handle JWT token storage and logic here
    } catch (e) {
      // Handle error
    }
  };

  return (
    <div className="flex flex-col flex-1">
      {/* Top image section */}
      <div className="w-full flex min-h-[400px] bg-secondary relative overflow-hidden">
        <AnimatePresence mode="wait">
          {isSignIn ? (
            <motion.div
              key="signin-image"
              initial={{ opacity: 0, x: 50 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -50 }}
              transition={{ duration: 0.5 }}
              className="absolute right-10 aspect-[2.2/3] bottom-0 h-[300px] z-10"
            >
              <Image
                src="/static/images/sign-in-human.png"
                alt="sign-in-human"
                fill
                priority
                sizes="x2"
              />
            </motion.div>
          ) : (
            <motion.div
              key="signup-image"
              initial={{ opacity: 0, x: -50 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 50 }}
              transition={{ duration: 0.5 }}
              className="absolute left-10 aspect-[2/3.3] bottom-0 h-[300px] z-10"
            >
              <Image
                src="/static/images/young-boy-sign-up.png"
                alt="young-boy-sign-up"
                fill
                priority
                sizes="x2"
              />
            </motion.div>
          )}
        </AnimatePresence>

        <AnimatePresence mode="wait">
          <motion.div
            className="absolute aspect-[3.5/2] h-[200px] -rotate-12 top-1/2 left-1/2 -translate-y-1/2 flex items-end justify-end z-2"
            animate={{ x: isSignIn ? "-50%" : "0%" }}
            transition={{ type: "spring", stiffness: 300, damping: 30 }}
          >
            <Image
              src="/static/images/detailed-sisu-pass.png"
              alt="detailed-sisu-pass"
              fill
            />
          </motion.div>
        </AnimatePresence>
        <AnimatePresence mode="wait">
          <motion.div
            className="absolute aspect-[3.5/2] h-[200px] rotate-12 top-1/2 left-1/2 -translate-y-1/2 flex items-end justify-end"
            animate={{ x: isSignIn ? "-50%" : "0%" }}
            transition={{ type: "spring", stiffness: 300, damping: 30 }}
          >
            <Image src="/static/images/card-noise.png" alt="card-noise" fill />
          </motion.div>
        </AnimatePresence>
      </div>

      {/* Content section */}
      <div className="flex flex-1">
        <div className="max-w-[640px] grid grid-cols-6 gap-x-4 gap-y-6 p-6 pt-3 w-full relative min-h-full">
          <AnimatePresence mode="wait">
            {isSignIn ? (
              <motion.div
                key="signin-content"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.4 }}
                className="col-span-6 flex flex-col flex-1 h-full justify-around mt-8"
              >
                <div className="flex flex-col">
                  <h1 className="text-2xl font-bold text-center">
                    Tap in to SisuPass
                  </h1>
                  <p className="text-base tracking-wide text-center">
                    Access your smart wallet. Keep the good times rolling.
                  </p>
                </div>

                <div className="flex flex-col gap-3">
                  <Button size="xl">Sign In Email</Button>
                  <GoogleOAuthProvider clientId={GOOGLE_AUTH_CLIENT_ID}>
                    <GoogleSignInButton
                      handleGoogleSignIn={handleGoogleSignIn}
                    />
                  </GoogleOAuthProvider>
                </div>

                <div className="flex flex-col">
                  <p className="py-5 text-center text-base">
                    Don’t have an account?
                    <button
                      onClick={() => setMode("signup")}
                      className="ml-2 font-semibold text-secondary"
                    >
                      Sign Up
                    </button>
                  </p>
                </div>
              </motion.div>
            ) : (
              <motion.div
                key="signup-content"
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.4 }}
                className="col-span-6 flex flex-col flex-1 h-full justify-around mt-8"
              >
                <div className="flex flex-col">
                  <h1 className="text-2xl font-bold text-center">
                    Join the SisuPass Family
                  </h1>
                  <p className="text-base tracking-wide text-center">
                    Sign up and make school payments simple, fun, and secure.
                  </p>
                </div>

                <div className="flex flex-col gap-3">
                  <Button size="xl">Sign Up Email</Button>
                  <Button size="xl" variant="outline">
                    Continue With Google
                  </Button>
                </div>

                <div className="flex flex-col">
                  <p className="py-5 text-center text-base">
                    Have an account?
                    <button
                      onClick={() => setMode("signin")}
                      className="ml-2 font-semibold text-secondary"
                    >
                      Sign In
                    </button>
                  </p>
                </div>
              </motion.div>
            )}
          </AnimatePresence>
        </div>
      </div>
    </div>
  );
}
