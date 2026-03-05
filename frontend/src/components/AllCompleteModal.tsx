import { useEffect } from "react";
import confetti from "canvas-confetti";

import { POPUP_CONFIG } from "../config";
import type { AllCompleteEvent } from "../types/student";

interface Props {
  event: AllCompleteEvent;
  onDismiss: () => void;
}

export function AllCompleteModal({ event, onDismiss }: Props) {
  useEffect(() => {
    const timer = setTimeout(onDismiss, POPUP_CONFIG.allCompleteDismissSeconds * 1000);

    const duration = POPUP_CONFIG.allCompleteDismissSeconds * 1000;
    const end = Date.now() + duration;

    function frame() {
      confetti({
        particleCount: 3,
        angle: 60,
        spread: 55,
        origin: { x: 0, y: 0.7 },
      });
      confetti({
        particleCount: 3,
        angle: 120,
        spread: 55,
        origin: { x: 1, y: 0.7 },
      });
      if (Date.now() < end) {
        requestAnimationFrame(frame);
      }
    }
    frame();

    return () => clearTimeout(timer);
  }, [onDismiss]);

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm"
      onClick={onDismiss}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => e.key === "Enter" && onDismiss()}
    >
      <div className="animate-scale-in rounded-2xl bg-white px-10 py-8 text-center shadow-2xl">
        <p className="text-4xl mb-3">&#127881;</p>
        <p className="text-xl font-bold">{event.studentName}</p>
        <p className="mt-1 text-slate-600">Completed all stages!</p>
      </div>
    </div>
  );
}
