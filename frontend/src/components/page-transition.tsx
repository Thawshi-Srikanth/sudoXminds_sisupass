"use client";

import { motion, AnimatePresence } from "motion/react";
import { usePathname } from "next/navigation";
import type { ReactNode } from "react";

interface PageTransitionProps {
  children: ReactNode;
}

const pageVariants = {
  initial: {
    x: "100%",
    opacity: 0,
  },
  in: {
    x: 0,
    opacity: 1,
  },
  out: {
    x: "-100%",
    opacity: 0,
  },
};

const pageTransition = {
  type: "tween",
  ease: [0.25, 0.46, 0.45, 0.94], // iOS-like cubic bezier easing
  duration: 0.3, // Faster iOS-like duration
};

export default function PageTransition({ children }: PageTransitionProps) {
  const pathname = usePathname();

  return (
    <AnimatePresence mode="wait" initial={false}>
      <motion.div
        key={pathname}
        initial="initial"
        animate="in"
        exit="out"
        variants={pageVariants}
        transition={pageTransition}
        className="w-full"
      >
        {children}
      </motion.div>
    </AnimatePresence>
  );
}
