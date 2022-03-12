import { Coordinates } from './types';

export function euclideanDistance(c1: Coordinates, c2: Coordinates): number {
  return Math.sqrt(
    Math.pow(c1.lng - c2.lng, 2) + Math.pow(c1.lat - c2.lat, 2) * 1.0
  );
}
