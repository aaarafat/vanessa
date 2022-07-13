import React, { useCallback, useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import styled, { keyframes } from 'styled-components';
import { useAppDispatch, useAppSelector } from '../store';
import { RsuState } from '../store/rsuSlice';
import { ARP } from './ARP';
import { InfoCard } from './info-card';
import { ObstaclesTable } from './obstacles-table';
import { initRsu } from '../store/rsuSlice';
import { useEventSource } from '../hooks';
import { ConnectionErrorAlert } from './connection-error-alert';

const spin = keyframes`
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
`;

const Loader = styled.div`
  border: 10px solid #f3f3f3;
  border-top: 10px solid #3498db;
  border-radius: 50%;
  width: 80px;
  height: 80px;
  animation: ${spin} 1s linear infinite;
`;

const LoaderContainer = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  z-index: 999999;
  background-color: rgba(0, 0, 0, 0.61);

  backdrop-filter: blur(5px);
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const STables = styled.div`
  display: flex;
`;
const OTable = styled.div`
  width: 30%;
`;
const ARPTable = styled.div`
  width: 70%;
`;
const SCards = styled.div`
  display: flex;
  flex-wrap: wrap;
`;
const SHeading = styled.h1`
  font-size: 5rem;
  font-family: 'Bebas Neue', cursive;
  margin: 1rem;
  color: #0d0d0d;
`;
const SDiscription = styled.h3`
  font-size: 2.5rem;
  font-family: 'Bebas Neue', cursive;
  margin: 0 1rem;
  color: #ffc000;
`;
export const Interface: React.FC = () => {
  const rsuState = useAppSelector((state) => state.rsu);
  const [eventSource, setEventSource] = useEventSource();
  const [loading, setLoading] = useState(true);
  const [connectionError, setConnectionError] = useState(false);
  const dispatch = useAppDispatch();
  const { port } = useParams<{
    port: string;
  }>();

  const connectRsu = useCallback(async () => {
    setLoading(true);
    setConnectionError(false);
    try {
      const state = await fetch(`http://localhost:${port}/state-rsu`);
      const json: RsuState = await state.json();
      dispatch(initRsu(json));

      setEventSource(new EventSource(`http://localhost:${port}`));

      document.title = `Rsu - ${json.id}`;
    } catch (e) {
      setConnectionError(true);
      setEventSource(null);
    } finally {
      setLoading(false);
    }
  }, [port, setEventSource, dispatch]);

  useEffect(() => {
    connectRsu();
    return () => {
      document.title = `RsuUi`;
    };
  }, [connectRsu]);

  return (
    <>
      <ConnectionErrorAlert
        connectRSU={connectRsu}
        connectionError={connectionError}
      />
      {loading && (
        <LoaderContainer>
          <Loader />
        </LoaderContainer>
      )}
      <SHeading>RSU Overview</SHeading>
      <SDiscription>Total packets</SDiscription>
      <SCards>
        <InfoCard
          title="Sent to RSUs"
          info={rsuState.sentToRsus}
          color="#ffffff"
          bgc="#0d0d0d"
        />
        <InfoCard
          title="Recieved from RSUs"
          info={rsuState.receivedFromRsus}
          color="#ffffff"
          bgc="#ffc000"
        />
        <InfoCard
          title="Sent to cars"
          info={rsuState.sentToCars}
          color="#ffffff"
          bgc="#2d2d2d"
        />
        <InfoCard
          title="Recieved from cars"
          info={rsuState.receivedFromCars}
          color="#ffffff"
          bgc="#ffcf33"
        />
      </SCards>

      <STables>
        <OTable>
          <ObstaclesTable />
        </OTable>
        <ARPTable>
          <ARP />
        </ARPTable>
      </STables>
    </>
  );
};
