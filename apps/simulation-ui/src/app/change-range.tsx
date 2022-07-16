import React, { useEffect, useState } from 'react';
import styled from 'styled-components';
import { useAppSelector } from './store';
import { socketEvents } from '../context';

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

const ChangeRange: React.FC = () => {
  const [range, setRange] = useState(500);
  const { focusedRSU } = useAppSelector((state) => state.simulation);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!focusedRSU) return;
    focusedRSU.setRadius(range);
    socketEvents.addRSU(focusedRSU);
  };

  return (
    <Container open={!!focusedRSU}>
      <Form onSubmit={handleSubmit}>
        <Label htmlFor="speed">Change Range (meters)</Label>
        <Input
          id="range"
          type="number"
          min="100"
          max="500"
          value={range}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) =>
            setRange(+e.target.value)
          }
        />
        <PrimaryButton>Apply</PrimaryButton>
      </Form>
    </Container>
  );
};

export default ChangeRange;
