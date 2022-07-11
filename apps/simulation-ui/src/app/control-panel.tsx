import React, { useContext, useEffect } from 'react';
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
  overflow: auto;
`;

const CancelButtonContainer = styled.div<{ show: boolean }>`
  display: ${(props) => (props.show ? 'flex' : 'none')};
  justify-content: flex-end;
  position: absolute;
  top: 0;
  right: 0;
  padding: 0.5rem;
  margin: 3rem 1rem;
  z-index: 1 !important;
  cursor: pointer;
  background-color: #010942ed;
  border-radius: 100%;
  color: #ffffff;
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

      mapDirections.on('route', () => {
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        const directions = (map.getSource('directions') as any)
          ._data as GeoJSON.FeatureCollection;

        // get origin coordinates
        const feature = directions.features.find(
          (f: GeoJSON.Feature) => f.properties?.route === 'selected'
        );

        if (feature?.geometry.type !== 'LineString') return;
        const coordinates = feature.geometry.coordinates;
        const route: Coordinates[] = coordinates.map(
          ([lng, lat]: number[]) => ({ lng, lat })
        );

        handleCarInputsChange({
          route,
        });
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
      <CancelButtonContainer
        show={!!accidentInput}
        onClick={() => mapDirections?.reset()}
      >
        <svg
          className="svg-icon"
          style={{
            width: '1em',
            height: '1em',
            verticalAlign: 'middle',
            fill: 'currentColor',
            overflow: 'hidden',
          }}
          viewBox="0 0 1024 1024"
          version="1.1"
          xmlns="http://www.w3.org/2000/svg"
        >
          <path d="M896.211 415.947c0 17.51-14.508 32.018-32.018 32.018L640.07 447.965c-13.008 0-24.514-8.004-29.518-20.011-5.001-11.506-2.5-25.514 7.004-34.519l69.039-69.038C639.57 280.873 577.536 255.859 512 255.859c-141.078 0-256.141 115.063-256.141 256.141 0 141.077 115.063 256.141 256.141 256.141 79.544 0 153.084-36.02 202.111-99.555 2.5-3.502 7.004-5.503 11.506-6.003 4.502 0 9.004 1.501 12.506 4.502l68.539 69.038c6.002 5.503 6.002 15.008 1 21.512C734.621 845.684 626.563 896.211 512 896.211c-211.616 0-384.211-172.595-384.211-384.211S300.384 127.789 512 127.789c98.553 0 194.105 39.521 264.645 106.058l65.035-64.535c9.006-9.505 23.014-12.007 35.02-7.004 11.508 5.002 19.512 16.509 19.512 29.516L896.212 415.947z" />
        </svg>
      </CancelButtonContainer>
    </>
  );
};

export default ControlPanel;
