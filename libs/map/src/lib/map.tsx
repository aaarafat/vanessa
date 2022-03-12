import styled from 'styled-components';
import React, { useEffect, useState, useContext } from 'react';
import mapboxgl from 'mapbox-gl';
import { MapProps } from './map.props';
import { MapContext } from './context';
import { Car, drawNewCar, updateCar } from '@vanessa/utils';

import './mapbox-gl.css';
import './mapbox-gl-directions.css';

// eslint-disable-next-line @typescript-eslint/no-var-requires
const MapboxDirections = require('@mapbox/mapbox-gl-directions/dist/mapbox-gl-directions');

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
  right: 0;
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
  currentZoom = 15.79,
  currentLat = 30.0246,
  currentLng = 31.211,
  cars = [],
}) => {
  const { map, setOptions, mapRef } = useContext(MapContext);

  const [lng, setLng] = useState(currentLng);
  const [lat, setLat] = useState(currentLat);
  const [zoom, setZoom] = useState(currentZoom);

  function onInit(map: mapboxgl.Map) {
    // const directions = new MapboxDirections({
    //   accessToken: mapboxgl.accessToken,
    //   unit: 'metric',
    //   profile: 'mapbox/driving',
    //   alternatives: 'true',
    //   geometries: 'geojson',
    // });

    // console.log(directions);

    // map.addControl(directions, 'top-right');

    map.on('move', () => {
      setLng(Number(map.getCenter().lng.toFixed(4)));
      setLat(Number(map.getCenter().lat.toFixed(4)));
      setZoom(Number(map.getZoom().toFixed(2)));
    });

    map.on('load', () => {
      console.log(cars);
      cars.forEach((car) => carHandler(map, car));

      // todo: update/add event for cars
    });
  }

  // Initialize map when component mounts
  useEffect(() => {
    setOptions({
      style: 'mapbox://styles/mapbox/streets-v11',
      center: [lng, lat],
      zoom: zoom,
      onInit,
    });
  }, []); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    console.log('map', map);
  }, [map]);

  return (
    <StyledMap>
      <StyledSidebar>
        <div>
          Longitude: {lng} | Latitude: {lat} | Zoom: {zoom}
        </div>
      </StyledSidebar>
      <StyledContainer ref={mapRef} />
    </StyledMap>
  );
};

export default Map;
function carHandler(map: mapboxgl.Map, car: Car) {
  const sourceId = `car-${car.id}`;
  const currSource = map.getSource(sourceId);
  const isNewCar = !currSource;

  if (isNewCar) {
    drawNewCar(map, sourceId, car);
  } else {
    updateCar(map, sourceId, car);
    // todo: remove car
  }
}
