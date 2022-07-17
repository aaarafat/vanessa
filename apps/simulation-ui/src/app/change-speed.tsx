import React, { useContext, useEffect, useState } from 'react';
import styled from 'styled-components';
import { useAppDispatch, useAppSelector } from './store';
import { SocketContext, socketEvents } from '../context';
import { addMessage } from './store/simulationSlice';

const Container = styled.div<{ open: boolean }>`
  display: flex;
  position: absolute;
  bottom: 0;
  right: ${(props) => (props.open ? 0 : '-100%')};
  background-color: #0d0d0df4;
  font-family: 'Bebas Neue', cursive;
  color: #ffffff;
  z-index: 1 !important;
  padding: 1rem 2rem;
  font-weight: bold;
  margin: 1rem;
  width: min(200px, 20%);
  align-items: stretch;
  transition: right 0.3s ease-in-out;
  flex-direction: column;
  overflow: auto;
`;

const Label = styled.label`
  font-size: 1.75rem;
`;

const Form = styled.form`
  display: flex;
  align-items: stretch;
  flex-direction: column;
`;

const PrimaryButton = styled.button`
  margin: 1rem 0;
  padding: 1rem;
  background-color: #ffc000;
  color: #ffffff;
  font-family: 'Bebas Neue', cursive;
  font-weight: bold;
  /* width: 100%; */
  border-radius: 1px;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.3);
  cursor: pointer;
  font-size: 1.5rem;
  &:disabled {
    cursor: not-allowed;
    opacity: 0.5;
  }
`;

const Input = styled.input`
  margin: 1rem 0;
  padding: 1rem;
  background-color: #fff;
  color: #000;
  font-weight: bold;
  /* width: 100%; */
  border-radius: 1px;
  box-shadow: 0 1px 1px rgba(0, 0, 0, 0.3);
  font-family: 'Lato', sans-serif;
  font-size: 1.2rem;
`;

const ChangeSpeedForm: React.FC = () => {
  const [speed, setSpeed] = useState(10);
  const { focusedCar } = useAppSelector((state) => state.simulation);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!focusedCar) return;
    const msg = {
      id: focusedCar.id,
      ...focusedCar.coordinates,
      port: focusedCar.port,
      route: focusedCar.route,
      speed,
      stopped: focusedCar.manualStop,
    };
    socketEvents.updateCar(msg);
  };

  return (
    <Container open={!!focusedCar}>
      <Form onSubmit={handleSubmit}>
        <Label htmlFor="speed">Change Speed (km/h)</Label>
        <Input
          id="speed"
          type="number"
          min="10"
          max="100"
          value={speed}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setSpeed(+e.target.value)
          }
        />
        <PrimaryButton>Apply</PrimaryButton>
      </Form>
    </Container>
  );
};

export default ChangeSpeedForm;
