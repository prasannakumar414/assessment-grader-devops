import { useEffect } from "react";

import { POPUP_CONFIG } from "../config";
import type { StageCompleteEvent } from "../types/student";

const stageLabels: Record<string, string> = {
  github: "GitHub",
  docker: "Docker",
  k8s: "Kubernetes",
};

interface Props {
  event: StageCompleteEvent;
  onDismiss: () => void;
}

export function StageCompleteModal({ event, onDismiss }: Props) {
  useEffect(() => {
    const timer = setTimeout(onDismiss, POPUP_CONFIG.stageCompleteDismissSeconds * 1000);
    return () => clearTimeout(timer);
  }, [onDismiss]);

  return (
    <div
      style={{
        position: "fixed",
        inset: 0,
        zIndex: 9999,
        display: "flex",
        alignItems: "flex-end",
        justifyContent: "center",
        padding: "1rem",
      }}
      onClick={onDismiss}
    >
      <div
        style={{
          background: "#fff",
          borderRadius: "0.5rem",
          border: "1px solid #e2e8f0",
          padding: "1rem 1.5rem",
          boxShadow: "0 20px 25px -5px rgba(0,0,0,.1), 0 8px 10px -6px rgba(0,0,0,.1)",
          animation: "slide-up 0.3s ease-out",
          cursor: "pointer",
        }}
        onClick={onDismiss}
      >
        <p style={{ fontSize: "1rem", margin: 0 }}>
          <span style={{ fontWeight: 700 }}>{event.studentName}</span> completed the{" "}
          <span style={{ fontWeight: 700, color: "#2563eb" }}>
            {stageLabels[event.stageName] ?? event.stageName}
          </span>{" "}
          stage!
        </p>
      </div>
    </div>
  );
}
