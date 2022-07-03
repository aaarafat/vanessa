import React, { useContext, useEffect } from 'react';
import styled from 'styled-components';
import { useAppDispatch, useAppSelector } from './store';
import { SocketContext } from '../context';
import { addMessage } from './store/simulationSlice';

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
  overflow: auto;
`;

const Message = styled.div`
  border-bottom: 1px solid #838383;
  padding: 0.5rem;
`;

const MessagesViewer: React.FC = () => {
  const socket = useContext(SocketContext);
  const messages = useAppSelector(({ simulation }) =>
    simulation.focusedCar !== null
      ? simulation.carsReceivedMessages[simulation.focusedCar]
      : null
  );
  const dispatch = useAppDispatch();

  useEffect(() => {
    socket?.on('change', (data: ArrayBuffer) => {
      const message = JSON.parse(
        String.fromCharCode.apply(null, new Uint8Array(data) as any)
      );
      dispatch(
        addMessage({
          id: message.id as number,
          message: message.data,
        })
      );
    });
  }, [socket]);

  return (
    <Container open={!!messages}>
      {messages?.map((message: any, index: number) => (
        <Message key={index}>{message}</Message>
      ))}
    </Container>
  );
};

export default MessagesViewer;
