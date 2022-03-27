import React, { useContext, useEffect } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled from 'styled-components';
import { Car, Coordinates, ICar, IRSU, RSU } from '@vanessa/utils';
import mapboxgl from 'mapbox-gl';

const rsus: IRSU[] = [
  {
    id: 1,
    lng: 31.213,
    lat: 30.0252,
    radius: 0.25,
  },
  {
    id: 2,
    lng: 31.2029,
    lat: 30.0269,
    radius: 0.5,
  },
  {
    id: 3,
    lng: 31.2129,
    lat: 30.0185,
    radius: 0.5,
  },
];

const cars: Car[] = [];

const Container = styled.form<{ open: boolean }>`
  display: flex;
  position: absolute;
  top: 0;
  bottom: 0;
  left: ${(props) => (props.open ? 0 : '-100%')};
  background-color: #010942ed;
  color: #ffffff;
  z-index: 1 !important;
  padding: 1rem 2rem;
  font-weight: bold;
  margin: 1rem;
  width: min(380px, 20%);
  align-items: stretch;
  transition: left 0.3s ease-in-out;
  flex-direction: column;
`;

const PrimaryButton = styled.button`
  margin: 1rem 0;
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

const OpenButton = styled(SmallButton)<{ open: boolean }>`
  position: absolute;
  top: 0;
  margin: 3rem;
  ${(props) => (props.open ? 'display: none;' : '')}
`;

const Input = styled.input`
  margin: 1rem 0;
  padding: 1rem;
  background-color: #ccc;
  color: #000;
  font-weight: bold;
  /* width: 100%; */
  border-radius: 1.5rem;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.3);
  font-size: 1.2rem;
`;

export const Simulation: React.FC = () => {
  const { map } = useContext(MapContext);

  useEffect(() => {
    if (map) {
      map.setMaxZoom(18);
      map.on('load', () => {
        cars.map((car) => new Car({ ...car, map }));
        rsus.map((rsu) => new RSU({ ...rsu, map }));
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

const initialState = {
  speed: 10,
};

const ControlPanel: React.FC = () => {
  const { map, mapRef, mapDirections } = useContext(MapContext);
  const [carInputs, setCarInputs] = React.useState<Partial<ICar>>(initialState);
  const [isOpen, setIsOpen] = React.useState(true);

  const handleCarInputsChange = (newValue: Partial<ICar>) => {
    setCarInputs((prev) => ({
      ...prev,
      ...newValue,
    }));
  };

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
        const directions = (map.getSource('directions') as any)
          ._data as GeoJSON.FeatureCollection;

        // get origin coordinates
        const feature = directions.features.find(
          (f: GeoJSON.Feature) => f.properties?.route === 'selected'
        );

        if (feature?.geometry.type === 'LineString') {
          const [lng, lat] = feature.geometry.coordinates[0] || [];
          const route = feature?.geometry.coordinates.map(
            (el: number[]): Coordinates => {
              return { lng: el[0], lat: el[1] };
            }
          );

          handleCarInputsChange({
            lng,
            lat,
            route,
            originalDirections: feature,
          });

          // we can create car here
          (map.getSource(CLICK_SOURCE_ID) as mapboxgl.GeoJSONSource).setData(
            coordsToFeature({
              lng,
              lat,
            })
          );
        }
      });
    }
  }, [map, mapRef, mapDirections]);

  const handleAddCar = (e: React.FormEvent) => {
    e.preventDefault();
    if (!map) return;

    // if no previous click, return
    const source = map.getSource(CLICK_SOURCE_ID);
    if (!source || !carInputs.route) return;

    const car: Car = new Car({
      ...carInputs,
      map,
    });

    // add to cars list
    cars.push(car);
    car.on('click', () => {
      mapDirections.reset();
      mapDirections.freeze();
    });
    car.on('popup-closed', () => {
      mapDirections.unfreeze();
    });

    setCarInputs(initialState);
    mapDirections.reset();
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
      <Container open={isOpen} onSubmit={handleAddCar}>
        <SmallButton onClick={() => setIsOpen(false)} type="button">
          {'<'}
        </SmallButton>
        <label htmlFor="speed">Speed (km/h)</label>
        <Input
          id="speed"
          type="number"
          min="10"
          max="100"
          value={carInputs.speed}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            handleCarInputsChange({ speed: +e.target.value })
          }
        />
        <PrimaryButton disabled={!carInputs.route}>Add Car</PrimaryButton>
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
