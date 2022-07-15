import React from 'react';
import styled from 'styled-components';
import { useAppSelector } from '../store';

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
  margin: 3rem 0 1rem 1rem;
  width: min(200px, 20%);
  align-items: stretch;
  transition: right 0.3s ease-in-out;
  flex-direction: column;
  overflow-y: auto;
  overflow-x: visible;
  font-size: 1.5rem;
`;

const Message = styled.div`
  border-bottom: 1px solid #838383;
  padding: 0.5rem;
  width: 100%;
`;

const MessagesViewer: React.FC = () => {
  const messages = useAppSelector((state) => state.car.messages);

  return (
    <Container open={!!messages}>
      {messages?.map((message: string, index: number) => (
        <Message key={index}>{message}</Message>
      ))}
    </Container>
  );
};

export default MessagesViewer;
