import styled from 'styled-components';
import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from 'mapbox-gl';

/* eslint-disable-next-line */
export interface MapProps {}

const StyledMap = styled.div``;

const StyledContainer = styled.div`
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  right: 0;
`;

const StyledSidebar = styled.div`
  display: inline-block;
  position: absolute;
  top: 0;
  left: 0;
  margin: 12px;
  background-color: #404040;
  color: #ffffff;
  z-index: 1 !important;
  padding: 6px;
  font-weight: bold;
`;

mapboxgl.accessToken =
  'pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4M29iazA2Z2gycXA4N2pmbDZmangifQ.-g_vE53SD2WrJ6tFX7QHmA';

export const Map: React.FC<MapProps> = () => {
  const mapContainerRef = useRef<any>();

  const [lng, setLng] = useState(5);
  const [lat, setLat] = useState(34);
  const [zoom, setZoom] = useState(1.5);

  // Initialize map when component mounts
  useEffect(() => {
    const map = new mapboxgl.Map({
      container: mapContainerRef.current,
      style: 'mapbox://styles/mapbox/streets-v11',
      center: [lng, lat],
      zoom: zoom,
    });

    // Add navigation control (the +/- zoom buttons)
    map.addControl(new mapboxgl.NavigationControl(), 'top-right');

    map.on('move', () => {
      setLng(Number(map.getCenter().lng.toFixed(4)));
      setLat(Number(map.getCenter().lat.toFixed(4)));
      setZoom(Number(map.getZoom().toFixed(2)));
    });

    // Clean up on unmount
    return () => map.remove();
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  return (
    <StyledMap>
      <StyledSidebar>
        <div>
          Longitude: {lng} | Latitude: {lat} | Zoom: {zoom}
        </div>
      </StyledSidebar>
      <StyledContainer ref={mapContainerRef} />
    </StyledMap>
  );
};

export default Map;
