import { useCallback, useEffect, useRef, useState } from "react";

import { getToken } from "../auth";
import type { AllCompleteEvent, StageCompleteEvent } from "../types/student";

export type CelebrationEvent =
  | { type: "stage_complete"; data: StageCompleteEvent }
  | { type: "all_complete"; data: AllCompleteEvent };

export function useSSE() {
  const [connected, setConnected] = useState(false);
  const [registrationVersion, setRegistrationVersion] = useState(0);
  const [stageVersion, setStageVersion] = useState(0);
  const [currentEvent, setCurrentEvent] = useState<CelebrationEvent | null>(null);
  const queueRef = useRef<CelebrationEvent[]>([]);
  const showingRef = useRef(false);

  const showNext = useCallback(() => {
    if (queueRef.current.length > 0) {
      showingRef.current = true;
      const next = queueRef.current.shift()!;
      console.log("[SSE] showing event:", next.type, next.data);
      setCurrentEvent(next);
    } else {
      showingRef.current = false;
      setCurrentEvent(null);
    }
  }, []);

  const enqueue = useCallback(
    (evt: CelebrationEvent) => {
      console.log("[SSE] enqueue:", evt.type, evt.data);
      queueRef.current.push(evt);
      if (!showingRef.current) {
        showNext();
      }
    },
    [showNext]
  );

  const dismissCurrent = useCallback(() => {
    showNext();
  }, [showNext]);

  useEffect(() => {
    let es: EventSource | null = null;
    let reconnectTimer: ReturnType<typeof setTimeout>;

    function connect() {
      const token = getToken();
      const url = token ? `/api/events?token=${encodeURIComponent(token)}` : "/api/events";
      es = new EventSource(url);

      es.onopen = () => {
        console.log("[SSE] connected");
        setConnected(true);
      };
      es.onerror = () => {
        console.log("[SSE] error, reconnecting in 3s...");
        setConnected(false);
        es?.close();
        reconnectTimer = setTimeout(connect, 3000);
      };

      es.addEventListener("stage_complete", (e) => {
        console.log("[SSE] raw stage_complete:", e.data);
        try {
          const data: StageCompleteEvent = JSON.parse(e.data);
          enqueue({ type: "stage_complete", data });
          setStageVersion((v) => v + 1);
        } catch (err) {
          console.error("[SSE] failed to parse stage_complete:", err, e.data);
        }
      });

      es.addEventListener("all_complete", (e) => {
        console.log("[SSE] raw all_complete:", e.data);
        try {
          const data: AllCompleteEvent = JSON.parse(e.data);
          enqueue({ type: "all_complete", data });
          setStageVersion((v) => v + 1);
        } catch (err) {
          console.error("[SSE] failed to parse all_complete:", err, e.data);
        }
      });

      es.addEventListener("new_registration", (e) => {
        console.log("[SSE] new_registration:", e.data);
        setRegistrationVersion((v) => v + 1);
      });
    }

    connect();

    return () => {
      clearTimeout(reconnectTimer);
      es?.close();
    };
  }, [enqueue]);

  return { connected, currentEvent, dismissCurrent, registrationVersion, stageVersion };
}
