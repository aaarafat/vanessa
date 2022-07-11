import { Coordinates } from './types';
import * as turf from '@turf/turf';

export function interpolateString(str: string, obj: any): string {
  return str.replace(/{([^{}]*)}/g, (a: string, b: string) => {
    const r = obj[b];
    return typeof r === 'string' || typeof r === 'number' ? String(r) : b;
  });
}

export function createFeaturePoint(
  c: Coordinates | turf.Position
): turf.Feature<turf.Point> {
  const coordinates = Array.isArray(c) ? c : [c.lng, c.lat];
  return {
    type: 'Feature',
    geometry: {
      type: 'Point',
      coordinates,
    },
    properties: {},
  };
}

export function getObstacleFeatures(obstacles: turf.Feature<turf.Point>[]) {
  const featureCollection: turf.FeatureCollection = {
    type: 'FeatureCollection',
    features: obstacles,
  };
  const obstacle = turf.buffer(featureCollection, 2, { units: 'meters' });
  return obstacle;
}
