import { Coordinates } from './types';

export function euclideanDistance(c1: Coordinates, c2: Coordinates): number {
  return Math.sqrt(
    Math.pow(c1.lng - c2.lng, 2) + Math.pow(c1.lat - c2.lat, 2) * 1.0
  );
}

/**
 * get distance in KM between 2 coordinates
 *
 * @param c1 coordinates 1
 * @param c2 coordinates 2
 * @returns distance in KM
 * @see https://www.geeksforgeeks.org/program-distance-two-points-earth/
 */
export function distanceInKm(c1: Coordinates, c2: Coordinates): number {
  // The math module contains a function
  // named toRadians which converts from
  // degrees to radians.
  const lon1 = (c1.lng * Math.PI) / 180;
  const lon2 = (c2.lng * Math.PI) / 180;
  const lat1 = (c1.lat * Math.PI) / 180;
  const lat2 = (c2.lat * Math.PI) / 180;

  // Haversine formula
  const dlon = lon2 - lon1;
  const dlat = lat2 - lat1;
  const a =
    Math.pow(Math.sin(dlat / 2), 2) +
    Math.cos(lat1) * Math.cos(lat2) * Math.pow(Math.sin(dlon / 2), 2);

  const c = 2 * Math.asin(Math.sqrt(a));

  // Radius of earth in kilometers. Use 3956
  // for miles
  const r = 6371;

  // calculate the result
  return c * r;
}
