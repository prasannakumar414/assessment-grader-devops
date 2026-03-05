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
      className="fixed inset-0 z-50 flex items-end justify-center p-4 pointer-events-none"
      onClick={onDismiss}
    >
      <div
        className="pointer-events-auto animate-slide-up rounded-lg border border-slate-200 bg-white px-6 py-4 shadow-xl"
        onClick={onDismiss}
        role="button"
        tabIndex={0}
        onKeyDown={(e) => e.key === "Enter" && onDismiss()}
      >
        <p className="text-base">
          <span className="font-bold">{event.studentName}</span> completed the{" "}
          <span className="font-bold text-blue-600">{stageLabels[event.stageName] ?? event.stageName}</span>{" "}
          stage!
        </p>
      </div>
    </div>
  );
}
