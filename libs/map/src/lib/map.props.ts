import { Car } from '@vanessa/utils';
import React from 'react';

export interface MapProps {
  /**Current Zoom value */
  currentZoom?: number;

  /**Current Lng value */
  currentLng?: number;

  /**Current Lat value */
  currentLat?: number;

  cars?: Car[];
}

export interface MapContextInterface {
  map: mapboxgl.Map | null;
  mapRef: React.MutableRefObject<HTMLDivElement | null> | null;
  setOptions: (options: MapOptions) => void;
}

export interface MapOptions extends Omit<mapboxgl.MapboxOptions, 'container'> {
  onInit?: (map: mapboxgl.Map) => void;
}
