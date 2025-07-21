import "dotenv/config";

if (!process.env.HEARTBEAT_SERVER_URL) {
  throw new Error("HEARTBEAT_SERVER_URL environment variable not configured");
}

if (!process.env.RIDE_REQUEST_SERVER_URL) {
  throw new Error(
    "RIDE_REQUEST_SERVER_URL environment variable not configured"
  );
}

export const serverConfig = {
  services: {
    heartbeat: {
      url: process.env.HEARTBEAT_SERVER_URL,
    },
    rideRequest: {
      url: process.env.RIDE_REQUEST_SERVER_URL,
    },
  },
};
