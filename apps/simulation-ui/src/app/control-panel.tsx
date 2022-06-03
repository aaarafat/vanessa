import React, { useContext, useEffect } from 'react';
import { CLICK_SOURCE_ID } from './constants';
import { Coordinates, ICar } from '@vanessa/utils';
import { MapContext } from '@vanessa/map';
import styled from 'styled-components';

const Container = styled.div<{ open: boolean }>`
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
  width: min(200px, 20%);
  align-items: stretch;
  transition: left 0.3s ease-in-out;
  flex-direction: column;
`;

const Form = styled.form`
  display: flex;
  align-items: stretch;
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

const initialState = {
  speed: 10,
};

const ControlPanel: React.FC<{
  onAddObstacle: (coordinates: Coordinates | null) => void;
  onAddCar: (carInputs: Partial<ICar>) => void;
  onExport: () => void;
  onImport: () => void;
  onClearMap: () => void;
}> = ({ onAddObstacle, onAddCar, onExport, onImport, onClearMap }) => {
  const { map, mapRef, mapDirections } = useContext(MapContext);
  const [carInputs, setCarInputs] = React.useState<Partial<ICar>>(initialState);
  const [accidentInput, setAccidentInput] = React.useState<Coordinates>();
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

      mapDirections.on('reset', () => {
        setCarInputs(initialState);
        setAccidentInput(undefined);
      });
    }
  }, [map, mapRef, mapDirections]);

  useEffect(() => {
    if (map) {
      mapDirections.on('origin', (e: any) => {
        const [lng, lat] = e.feature.geometry.coordinates || [];
        const coordinates = { lng, lat };
        setAccidentInput(coordinates);
      });
    }
  }, [map, mapDirections]);

  return (
    <>
      <OpenButton
        onClick={() => setIsOpen(true)}
        open={isOpen}
        disabled={isOpen}
      >
        {'>'}
      </OpenButton>
      <Container
        open={isOpen}
        onSubmit={(e) => {
          e.preventDefault();
          onAddCar(carInputs);
        }}
      >
        <Form>
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
        </Form>
        <PrimaryButton
          disabled={!accidentInput}
          onClick={() => {
            onAddObstacle(accidentInput || null);
            mapDirections.reset();
          }}
        >
          Add Accident
        </PrimaryButton>
        <PrimaryButton onClick={onExport}>Export</PrimaryButton>
        <PrimaryButton onClick={onImport}>Import</PrimaryButton>
        <PrimaryButton onClick={onClearMap}>Clear Map</PrimaryButton>
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

export default ControlPanel;
