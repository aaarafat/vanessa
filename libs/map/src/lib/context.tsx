import React, { useCallback, useEffect, useRef, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import { MapContextInterface, MapOptions } from './map.props';
// eslint-disable-next-line @typescript-eslint/no-var-requires
import { Directions } from './directions';

export const MapContext = React.createContext<MapContextInterface>({
  map: null,
  mapDirections: null,
  mapRef: null,
  setOptions: () => null,
});

export const MapProvider: React.FC = (props) => {
  const { children } = props;
  const mapRef = useRef<HTMLDivElement>(null);
  const [map, setMap] = useState<mapboxgl.Map | null>(null);
  const [options, setOptions] = useState<MapOptions>({});
  const [initialized, setInitialized] = useState(false);
  const [mapDirections, setMapDirections] = useState<Directions | null>(null);

  useEffect(() => {
    const { onInit, ...restOptions } = options;
    if (!initialized && onInit && mapRef.current) {
      const map = new mapboxgl.Map({
        ...restOptions,
        container: mapRef.current,
      });

      const directions = new Directions({
        accessToken: restOptions.accessToken,
        unit: 'metric',
        profile: 'mapbox/driving',
        alternatives: 'true',
        geometries: 'geojson',
        controls: {
          instructions: false,
          profileSwitcher: false,
        },
      });

      setMap(map);
      setMapDirections(directions);
      setInitialized(true);
      onInit(map);
    }
  }, [map, options, initialized]);

  useEffect(() => {
    return () => map?.remove();
  }, [map]);

  return (
    <MapContext.Provider value={{ map, mapRef, setOptions, mapDirections }}>
      {children}
    </MapContext.Provider>
  );
};
