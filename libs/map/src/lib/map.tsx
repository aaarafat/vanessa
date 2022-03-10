import styled from 'styled-components';
import React, { useRef, useEffect, useState } from 'react';
import mapboxgl, { GeoJSONSourceRaw } from 'mapbox-gl';
import { Car, MapProps } from './map.props';

import "./mapbox-gl.css";
import "./mapbox-gl-directions.css";

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

const carDefaultProps = {
  title: 'Car',
  'marker-size': 'large',
  'marker-color': '#f00',
}

export const Map: React.FC<MapProps> = ({
  currentZoom = 15.79,
  currentLat = 30.0246,
  currentLng = 31.211,
  cars = [],
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

    const directions = new MapboxDirections({
      accessToken: mapboxgl.accessToken,
      unit: 'metric',
      profile: 'mapbox/driving',
      alternatives: 'true',
      geometries: 'geojson',
    });

    map.addControl(directions, 'top-right');

    map.on('move', () => {
      setLng(Number(map.getCenter().lng.toFixed(4)));
      setLat(Number(map.getCenter().lat.toFixed(4)));
      setZoom(Number(map.getZoom().toFixed(2)));
    });

    map.on('load', () => {
      console.log(cars)
      cars.forEach((car) =>
        carHandler(map, car)
      );

      // todo: update/add event for cars
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

function drawNewCar(map: mapboxgl.Map, sourceId: string, car: Car) {
  const geojson: GeoJSONSourceRaw = {
    type: "geojson",
    data: getCarLocation(car)
  };
  map.addSource(sourceId, geojson);

  map.addLayer({
    id: sourceId,
    // type: 'symbol',
    source: sourceId,
    type: 'circle',
    'paint': {
      'circle-radius': 10,
      'circle-color': '#007cbf',
    },
    // layout: {
    //   "text-field": "Car {id}",
    // }
  });
}

function updateCar(map: mapboxgl.Map, sourceId: string, car: Car) {
  (map.getSource(sourceId) as mapboxgl.GeoJSONSource).setData(getCarLocation(car));
}


function getCarLocation(car: Car): GeoJSON.FeatureCollection {
  const { lng, lat } = car;
  return {
    type: "FeatureCollection",
    features: [
      {
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [lng, lat],
        },
        properties: { ...carDefaultProps, ...car },
      }
    ],
  }
};
