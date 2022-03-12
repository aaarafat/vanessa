import React, { useEffect, useRef, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import { MapContextInterface, MapOptions } from './map.props';

export const MapContext = React.createContext<MapContextInterface>({
  map: null,
  mapRef: null,
  setOptions: () => null,
});

export const MapProvider: React.FC = (props) => {
  const { children } = props;
  const mapRef = useRef<HTMLDivElement>(null);
  const [map, setMap] = useState<mapboxgl.Map | null>(null);
  const [options, setOptions] = useState<MapOptions>({});
  const [initialized, setInitialized] = useState(false);

  useEffect(() => {
    const { onInit, ...restOptions } = options;
    if (!initialized && onInit && mapRef.current) {
      const map = new mapboxgl.Map({ ...restOptions, container: mapRef.current });
      setMap(map);
      setInitialized(true);
      onInit(map);
    }
  }, [map, options, initialized]);

  useEffect(() => {
    return () => map?.remove()
  }, [map]);

  return (
    <MapContext.Provider value={{ map, mapRef, setOptions }}>
      {children}
    </MapContext.Provider>
  );
};

