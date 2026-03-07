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
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 9999,
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        background: "rgba(0,0,0,0.3)",
        backdropFilter: "blur(4px)",
        cursor: "pointer",
      }}
      onClick={onDismiss}
    >
      <div
        style={{
          background: "#fff",
          borderRadius: "1rem",
          padding: "2rem 2.5rem",
          textAlign: "center",
          boxShadow: "0 25px 50px -12px rgba(0,0,0,.25)",
          animation: "scale-in 0.3s ease-out",
        }}
      >
        <p style={{ fontSize: "2.25rem", margin: "0 0 0.75rem" }}>&#127881;</p>
        <p style={{ fontSize: "1.25rem", fontWeight: 700, margin: 0 }}>{event.studentName}</p>
        <p style={{ marginTop: "0.25rem", color: "#64748b" }}>Completed all stages!</p>
      </div>
    </div>
  );
}
