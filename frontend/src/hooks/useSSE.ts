import { useCallback, useEffect, useRef, useState } from "react";

import { getToken } from "../auth";
import type { AllCompleteEvent, StageCompleteEvent } from "../types/student";

export type CelebrationEvent =
  | { type: "stage_complete"; data: StageCompleteEvent }
  | { type: "all_complete"; data: AllCompleteEvent };

export function useSSE() {
  const [connected, setConnected] = useState(false);
  const [registrationVersion, setRegistrationVersion] = useState(0);
  const [currentEvent, setCurrentEvent] = useState<CelebrationEvent | null>(null);
  const queueRef = useRef<CelebrationEvent[]>([]);
  const showingRef = useRef(false);

  const showNext = useCallback(() => {
    if (queueRef.current.length > 0) {
      showingRef.current = true;
      setCurrentEvent(queueRef.current.shift()!);
    } else {
      showingRef.current = false;
      setCurrentEvent(null);
    }
  }, []);

  const enqueue = useCallback(
    (evt: CelebrationEvent) => {
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

      es.onopen = () => setConnected(true);
      es.onerror = () => {
        setConnected(false);
        es?.close();
        reconnectTimer = setTimeout(connect, 3000);
      };

      es.addEventListener("stage_complete", (e) => {
        const data: StageCompleteEvent = JSON.parse(e.data);
        enqueue({ type: "stage_complete", data });
      });

      es.addEventListener("all_complete", (e) => {
        const data: AllCompleteEvent = JSON.parse(e.data);
        enqueue({ type: "all_complete", data });
      });

      es.addEventListener("new_registration", () => {
        setRegistrationVersion((v) => v + 1);
      });
    }

    connect();

    return () => {
      clearTimeout(reconnectTimer);
      es?.close();
    };
  }, [enqueue]);

  return { connected, currentEvent, dismissCurrent, registrationVersion };
}
