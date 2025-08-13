"use client";
import React, { useEffect, useRef, useState, Children, ReactNode } from "react";
import { motion, AnimatePresence } from "motion/react";
import { useGesture } from "@use-gesture/react";

export const CardStack = ({
  children,
  offset = 10,
  scaleFactor = 0.06,
  maxVisible = 2,
}: {
  children: ReactNode;
  offset?: number;
  scaleFactor?: number;
  maxVisible?: number;
}) => {
  const allCards = Children.toArray(children);
  const [startIndex, setStartIndex] = useState(0);
  const [cardHeight, setCardHeight] = useState<number>();
  const firstCardRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (firstCardRef.current) {
      setCardHeight(firstCardRef.current.offsetHeight);
    }
  }, []);

  const cards = [];
  for (let i = 0; i < maxVisible; i++) {
    cards.push(allCards[(startIndex + i) % allCards.length]);
  }

  const bind = useGesture({
    onDragEnd: ({ swipe: [swipeX] }) => {
      if (swipeX === 1) {
        setStartIndex((prev) => (prev - 1 + allCards.length) % allCards.length);
      } else if (swipeX === -1) {
        setStartIndex((prev) => (prev + 1) % allCards.length);
      }
    },
  });

  const totalHeight = cardHeight
    ? cardHeight + offset * (maxVisible - 1)
    : undefined;

  return (
    <div
      {...bind()}
      className="relative flex items-end w-full touch-none select-none"
      style={{ height: totalHeight, userSelect: "none", marginBottom: offset }}
    >
      <AnimatePresence initial={false}>
        {cards.map((child, index) => {
          const scale = 1 - index * scaleFactor;
          return (
            <motion.div
              key={startIndex + index}
              ref={index === 0 ? firstCardRef : null}
              className="absolute cursor-grab"
              style={{
                transformOrigin: "bottom center",
                height: cardHeight ? `${cardHeight}px` : undefined,
                width: "100%",
                bottom: index * -offset * scale,
                zIndex: maxVisible - index,
              }}
              initial={{ opacity: 0, y: 20, x: 0 }}
              animate={{ opacity: 1, y: 0, scale, x: 0 }}
              exit={{ opacity: 0, y: -20, x: 0 }}
              transition={{ type: "spring", stiffness: 300, damping: 30 }}
              drag="x"
              dragConstraints={{ left: 0, right: 0 }}
              dragElastic={0.2}
              dragMomentum={false}
              onDragEnd={(event, info) => {
                const velocityX = info.velocity.x;
                if (velocityX > 500) {
                  setStartIndex(
                    (prev) => (prev - 1 + allCards.length) % allCards.length
                  );
                } else if (velocityX < -500) {
                  setStartIndex((prev) => (prev + 1) % allCards.length);
                }
              }}
            >
              {child}
            </motion.div>
          );
        })}
      </AnimatePresence>
    </div>
  );
};
