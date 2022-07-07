import React, { useContext, useEffect, useState } from 'react';
import { Socket } from 'socket.io-client';
import styled from 'styled-components';
import { EventSourceContext } from '../context';

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
  const eventSource = useContext(EventSourceContext)[0];
  const [messages, setMessages] = useState<any[]>([]);

  useEffect(() => {
    function testHandler({ data }: EventSourceEventMap['message']) {
      console.log(JSON.parse(data));
      setMessages((messages) => [...messages, data]);
    }
    eventSource?.addEventListener('test', testHandler);
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    return () => eventSource?.removeEventListener('test', testHandler);
  }, [eventSource]);
  return (
    <Container open={!!messages}>
      {messages?.map((message: any, index: number) => (
        <Message key={index}>{message}</Message>
      ))}
    </Container>
  );
};

export default MessagesViewer;
