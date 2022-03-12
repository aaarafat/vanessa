import React, { useContext, useEffect } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled from 'styled-components';
import { Car, Coordinates, drawNewCar, updateCar } from '@vanessa/utils';

/*
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
*/

const cars: Car[] = [];

const Container = styled.div<{ open: boolean }>`
  display: flex;
  position: absolute;
  top: 0;
  bottom: 0;
  left: ${(props) => (props.open ? 0 : '-100%')};
  background-color: #010942ed;
  color: #ffffff;
  z-index: 1 !important;
  padding: 1rem;
  font-weight: bold;
  margin: 1rem;
  width: min(380px, 20%);
  align-items: stretch;
  transition: left 0.3s ease-in-out;
  flex-direction: column;
`;

const PrimaryButton = styled.button`
  margin: 1rem;
  padding: 1rem;
  background-color: #ccc;
  color: #000;
  font-weight: bold;
  /* width: 100%; */
  border-radius: 1.5rem;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.3);
  cursor: pointer;
  font-size: 1.2rem;
  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
`;

const SmallButton = styled(PrimaryButton)`
  align-self: flex-start;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
`;

const OpenButton = styled(SmallButton) <{ open: boolean }>`
  position: absolute;
  top: 0;
  margin: 3rem;
  ${(props) => (props.open ? 'display: none;' : '')}
`;

export const Simulation: React.FC = () => {
  const { map } = useContext(MapContext);

  useEffect(() => {
    if (map) {
      map.on('load', () => {
        function drawCars() {
          if (map) cars.forEach((car) => updateCar(map, `car-${car.id}`, car))
          requestAnimationFrame(drawCars);
        }
        drawCars();
      });
    }
  }, [map]);

  return (
    <div>
      <Map cars={cars} />
      <ControlPanel />
    </div>
  );
};

const CLICK_SOURCE_ID = 'click';

const ControlPanel: React.FC = () => {
  const { map, mapRef, mapDirections } = useContext(MapContext);
  const [coords, setCoords] = React.useState<Coordinates>();
  const [route, setRoute] = React.useState<Coordinates[]>();
  const [isOpen, setIsOpen] = React.useState(true);

  useEffect(() => {
    if (map) {
      // add directions controller
      map.addControl(mapDirections, 'top-right');

      // add click source on load
      map.on('load', () => {
        map.addSource(CLICK_SOURCE_ID, {
          type: 'geojson',
          data: {
            type: 'FeatureCollection',
            features: [],
          },
        });
      });

      mapDirections.on('route', () => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const directions = (map.getSource('directions') as any)._data;

        // get origin coordinates
        const [lng, lat] = directions.features[2].geometry.coordinates[0];
        setCoords({ lng, lat });

        // set route coordinates
        setRoute(
          directions.features[2].geometry.coordinates.map(
            (el: number[]): Coordinates => {
              return { lng: el[0], lat: el[1] };
            }
          )
        );

        // we can create car here
        (map.getSource(CLICK_SOURCE_ID) as mapboxgl.GeoJSONSource).setData(
          coordsToFeature({
            lng,
            lat,
          })
        );
      });
    }
  }, [map, mapRef, mapDirections]);

  const handleAddCar = () => {
    if (!map) return;

    // if no previous click, return
    const source = map.getSource(CLICK_SOURCE_ID);
    if (!source || !route) return;

    // id is current time
    const car: Car = new Car({
      id: Date.now(),
      lat: coords?.lat ?? 0,
      lng: coords?.lng ?? 0,
      route,
    });

    // add to cars list
    // TODO: replace this with store
    cars.push(car);

    drawNewCar(map, `car-${car.id}`, car);
    setCoords(undefined);
    setRoute(undefined);
    mapDirections.removeRoutes();
  };

  return (
    <>
      <OpenButton
        onClick={() => setIsOpen(true)}
        open={isOpen}
        disabled={isOpen}
      >
        {'>'}
      </OpenButton>
      <Container open={isOpen}>
        <SmallButton onClick={() => setIsOpen(false)}>{'<'}</SmallButton>
        <PrimaryButton onClick={handleAddCar} disabled={!coords}>
          Add Car
        </PrimaryButton>
      </Container>
    </>
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
