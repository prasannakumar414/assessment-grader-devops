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
        alignItems: "center",
        justifyContent: "center",
        background: "rgba(0,0,0,0.15)",
        cursor: "pointer",
      }}
      onClick={onDismiss}
    >
      <div
        style={{
          background: "#fff",
          borderRadius: "1rem",
          border: "1px solid #e2e8f0",
          padding: "2.5rem 3rem",
          boxShadow: "0 25px 50px -12px rgba(0,0,0,.25)",
          animation: "scale-in 0.3s ease-out",
          textAlign: "center",
          minWidth: "360px",
        }}
        onClick={(e) => e.stopPropagation()}
      >
        <p style={{ fontSize: "1.5rem", margin: 0, lineHeight: 1.4 }}>
          <span style={{ fontWeight: 700 }}>{event.studentName}</span>
          <br />
          completed the{" "}
          <span style={{ fontWeight: 700, color: "#2563eb" }}>
            {stageLabels[event.stageName] ?? event.stageName}
          </span>{" "}
          stage!
        </p>
      </div>
    </div>
  );
}
