import styled from 'styled-components';
import React, { useRef, useEffect, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import { MapProps } from './map.props';

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

export const Map: React.FC<MapProps> = ({
  currentZoom = 1.5,
  currentLat = 34,
  currentLng = 5,
}) => {
  const mapContainerRef = useRef<any>();

  const [lng, setLng] = useState(currentLng);
  const [lat, setLat] = useState(currentLat);
  const [zoom, setZoom] = useState(currentZoom);

  // Initialize map when component mounts
  useEffect(() => {
    const map = new mapboxgl.Map({
      container: mapContainerRef.current,
      style: 'mapbox://styles/mapbox/streets-v11',
      center: [lng, lat],
      zoom: zoom,
    });

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
