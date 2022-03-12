import React, { useState, useRef, useEffect } from 'react';
import mapboxgl from 'mapbox-gl';

interface useMapProps extends Omit<mapboxgl.MapboxOptions, "container"> {
  onInit?: (map: mapboxgl.Map) => void
}

interface useMapReturn {
  map: mapboxgl.Map | null;
  mapRef: React.MutableRefObject<HTMLDivElement | null>;
}

export const useMap = (props: useMapProps = {}): useMapReturn => {
    const { onInit, ...rest } = props;
    const mapRef = useRef<HTMLDivElement>(null);
    const [map, setMap] = useState<mapboxgl.Map | null>(null);
    const [initialized, setInitialized] = useState(false);

    useEffect(() => {
      if(!initialized && onInit && mapRef.current) {
        const map = new mapboxgl.Map({...rest, container: mapRef.current});
        setMap(map);
        setInitialized(true);
        onInit(map);
      }
    }, [map, rest, onInit, initialized]);

    useEffect(() => {
      return () => map?.remove()
    }, [map]);

    return { map, mapRef };
}