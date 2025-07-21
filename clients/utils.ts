import { MAP_BOUNDING_BOX } from "./constants";

export function getRandomCoordinate(): { lat: number; lng: number } {
  const lat =
    Math.random() * (MAP_BOUNDING_BOX.MAX_LAT - MAP_BOUNDING_BOX.MIN_LAT) +
    MAP_BOUNDING_BOX.MIN_LAT;
  const lng =
    Math.random() * (MAP_BOUNDING_BOX.MAX_LON - MAP_BOUNDING_BOX.MIN_LON) +
    MAP_BOUNDING_BOX.MIN_LON;
  return { lat, lng };
}
