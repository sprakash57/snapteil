import { useEffect, useLayoutEffect, useRef } from "react";
import type { Image } from "../../lib/types";

const WS_RECONNECT_DELAY = 3000;

export function useWebSocket(onNewImage: (image: Image) => void) {
  const onNewImageRef = useRef(onNewImage);
  useLayoutEffect(() => {
    onNewImageRef.current = onNewImage;
  });

  useEffect(() => {
    let ws: WebSocket | null = null;
    let timeoutId: ReturnType<typeof setTimeout> | null = null;
    let destroyed = false;

    function connect() {
      if (destroyed) return;

      const protocol = window.location.protocol === "https:" ? "wss" : "ws";
      ws = new WebSocket(`${protocol}://${window.location.host}/ws`);

      ws.onmessage = (event) => {
        try {
          const image: Image = JSON.parse(event.data as string);
          onNewImageRef.current(image);
        } catch {
          // ignore malformed messages
        }
      };

      ws.onclose = () => {
        if (!destroyed) {
          timeoutId = setTimeout(connect, WS_RECONNECT_DELAY);
        }
      };

      ws.onerror = () => {
        ws?.close();
      };
    }

    connect();

    return () => {
      destroyed = true;
      if (timeoutId !== null) clearTimeout(timeoutId);
      ws?.close();
    };
  }, []);
}
