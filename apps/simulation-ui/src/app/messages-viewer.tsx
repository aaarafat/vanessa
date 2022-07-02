import React from 'react';
import styled from 'styled-components';
import { useAppSelector } from './store';

const Container = styled.div<{ open: boolean }>`
  display: flex;
  position: absolute;
  top: 0;
  bottom: 0;
  right: ${(props) => (props.open ? 0 : '-100%')};
  background-color: #ecedf8ec;
  color: #000;
  z-index: 1 !important;
  padding: 1rem 2rem;
  font-weight: bold;
  margin: 1rem;
  width: min(200px, 20%);
  align-items: stretch;
  transition: right 0.3s ease-in-out;
  flex-direction: column;
`;

const Message = styled.div`
  border-bottom: 1px solid #838383;
  padding: 0.5rem;
`;

const MessagesViewer: React.FC<{}> = () => {
  const car = useAppSelector(({ simulation }) =>
    simulation.focusedCar !== null
      ? simulation.cars[simulation.focusedCar]
      : null
  );

  return (
    <Container open={!!car}>
      {car?.receivedMessages.map((message, index) => (
        <Message key={index}>{JSON.stringify(message, null, 2)}</Message>
      ))}
    </Container>
  );
};

export default MessagesViewer;
