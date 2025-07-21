import { DRIVER_CLIENT } from "../constants";
import { getRandomCoordinate } from "../utils";
import { serverConfig } from "../config";

interface Heartbeat {
  driver_id: string;
  latitude?: number;
  longitude?: number;
  timestamp: number; // Unix timestamp in seconds
  available_status?: 0 | 1;
}

function convertStatusToNumeric(status: "available" | "unavailable"): 0 | 1 {
  return status === "available" ? 1 : 0;
}

async function sendHeartBeat(heartbeat: Heartbeat) {
  const baseUrl = serverConfig.services.heartbeat.url;
  await fetch(`${baseUrl}/heartbeat`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(heartbeat),
  });
}

// Simulate
setInterval(() => {
  Promise.all(
    Array.from({ length: DRIVER_CLIENT.SIMUL_CLIENTS_COUNT }, (_, index) => {
      const driver_id = `driver${index + 1}`;
      const { lat: latitude, lng: longitude } = getRandomCoordinate();
      const timestamp = Math.floor(Date.now() / 1000);
      const available_status = convertStatusToNumeric(
        Math.random() > 0.5 ? "available" : "unavailable"
      );

      const heartbeat: Heartbeat = {
        driver_id,
        latitude,
        longitude,
        timestamp,
        available_status,
      };

      const randomDelayBeforeSend =
        Math.random() * DRIVER_CLIENT.SEND_HEARTBEAT_INTERVAL_MS;
      setTimeout(() => {
        return sendHeartBeat(heartbeat);
      }, randomDelayBeforeSend);
    })
  );
}, DRIVER_CLIENT.SEND_HEARTBEAT_INTERVAL_MS);
