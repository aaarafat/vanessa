import React, { useContext, useEffect } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled from 'styled-components';
import { Car } from '@vanessa/utils';

const cars: Car[] = [
  {
    id: 1,
    lat: 30.02543,
    lng: 31.21146,
  },
  {
    id: 2,
    lat: 30.02763,
    lng: 31.21082,
  },
  {
    id: 3,
    lat: 30.02425,
    lng: 31.20995,
  },
  {
    id: 4,
    lat: 30.02616,
    lng: 31.21075,
  },
];

const Container = styled.div`
  display: flex;
  position: absolute;
  top: 0;
  bottom: 0;
  left: 0;
  background-color: #010942ed;
  color: #ffffff;
  z-index: 1 !important;
  padding: 1rem;
  font-weight: bold;
  margin: 1rem;
  width: 25%;
  align-items: flex-start;
`;

const PrimaryButton = styled.button`
  margin: 1rem;
  padding: 1rem;
  background-color: #ccc;
  color: #000;
  font-weight: bold;
  /* width: 100%; */
  flex: 1;
  border-radius: 1.5rem;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.3);
  cursor: pointer;
  font-size: 1.2rem;
`;

export const Simulation: React.FC = () => {
  return (
    <div>
      <Map cars={cars} />
      <ControlPanel />
    </div>
  );
};

const ControlPanel: React.FC = (props) => {
  const { map, mapRef } = useContext(MapContext);
  useEffect(() => {
    console.log(map, mapRef);
    if (map) {
      map.addSource('click', {
        type: 'geojson',
        data: {
          type: 'FeatureCollection',
          features: [],
        },
      });
      map.on('click', (e) => {
        (map.getSource('click') as mapboxgl.GeoJSONSource).setData(
          coordsToFeature({
            lng: e.lngLat.lng,
            lat: e.lngLat.lat,
          })
        );
      });
    }
  }, [map, mapRef]);

  const handleAddCar = () => {
    console.log('add car');
  };

  // function drawCar() {
  //   // if no map, return
  //   if (!map) return;
  //   const id = "click";

  //   // if no previous click, return
  //   const source = map.getSource(id);
  //   if (!source) return;

  //   // get previous click coordinates, remove it, create new source for the new car, draw a car

  //   // const prevClick = source.getData().features[0].geometry.coordinates;
  //   const location = source;
  //   map.addSource("click", geojson);

  //   map.addLayer({
  //     id: id,
  //     source: id,
  //     type: 'circle',
  //     'paint': {
  //       'circle-radius': 10,
  //       'circle-color': '#a79c06',
  //     },
  //   });
  // }

  return (
    <Container>
      <PrimaryButton onClick={handleAddCar}>Add Car</PrimaryButton>
    </Container>
  );
};

function coordsToFeature({
  lng,
  lat,
}: {
  lng: number;
  lat: number;
}): GeoJSON.FeatureCollection {
  return {
    type: 'FeatureCollection',
    features: [
      {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [+lng, +lat],
        },
        properties: {},
      },
    ],
  };
}
