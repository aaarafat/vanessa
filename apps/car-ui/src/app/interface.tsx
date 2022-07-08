import React, { useCallback, useContext, useEffect, useState } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled, { keyframes } from 'styled-components';
import { Car, ICar } from '@vanessa/utils';
import * as turf from '@turf/turf';
import mapboxgl from 'mapbox-gl';
import { useEventSource } from '../hooks';
import { useHistory, useParams } from 'react-router-dom';
import MessagesViewer from './messages-viewer';
import { ConnectionErrorAlert } from './connection-error-alert';
import { useAppDispatch, useAppSelector } from '../store';
import { initCar } from '../store/carSlice';

type carState = Omit<ICar, 'map'>;

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

export const Interface: React.FC = () => {
  const { map } = useContext(MapContext);
  const [eventSource, setEventSource] = useEventSource();
  const [connectionError, setConnectionError] = useState(false);
  const [mapLoaded, setMapLoaded] = useState(false);
  const [loading, setLoading] = useState(true);
  // const [obstacles, setObstacles] = useState<turf.Feature[]>([]);
  const { car, obstacles } = useAppSelector((state) => state.car);
  const dispatch = useAppDispatch();
  const { port } = useParams<{
    port: string;
  }>();
  const history = useHistory();

  useEffect(() => {
    if (!port || !Number.isInteger(+port) || +port < 0 || +port > 65535) {
      history.replace('/');
      return;
    }
  }, [port, history, setEventSource]);

  useEffect(() => {
    if (map) {
      map.setMaxZoom(18);
      map.on('load', () => {
        setMapLoaded(true);
      });
    }
    return () => car?.remove();
  }, [map]);

  const connectCar = useCallback(async () => {
    setLoading(true);
    setConnectionError(false);
    try {
      if (!map) {
        setConnectionError(true);
        return;
      }
      const state = await fetch(`http://localhost:${port}/state`);
      const json: carState = await state.json();
      dispatch(initCar(new Car({ ...json, map }, { displayOnly: true })));
      setEventSource(new EventSource(`http://localhost:${port}`));
    } catch (e) {
      setConnectionError(true);
      setEventSource(null);
    } finally {
      setLoading(false);
    }
  }, [map, port, setEventSource]);

  useEffect(() => {
    if (mapLoaded) {
      connectCar();
    }
  }, [mapLoaded, connectCar]);

  useEffect(() => {
    if (mapLoaded && map) {
      const featureCollection: turf.FeatureCollection = {
        type: 'FeatureCollection',
        features: obstacles,
      };
      const obstacle = turf.buffer(featureCollection, 10, { units: 'meters' });
      if (!map.getSource('obstacles')) {
        map.addSource('obstacles', {
          type: 'geojson',
          data: obstacle,
        });

        map.addLayer({
          id: 'obstacles',
          type: 'fill',
          source: 'obstacles',
          layout: {},
          paint: {
            'fill-color': '#f03b20',
            'fill-opacity': 0.5,
            'fill-outline-color': '#f03b20',
          },
        });
      } else {
        (map.getSource('obstacles') as mapboxgl.GeoJSONSource).setData(
          obstacle
        );
      }
    }
  }, [obstacles, map, mapLoaded]);

  return (
    <>
      <ConnectionErrorAlert
        connectCar={connectCar}
        connectionError={connectionError}
      />
      {loading && (
        <LoaderContainer>
          <Loader />
        </LoaderContainer>
      )}
      <div>
        <Map />
        <MessagesViewer />
      </div>
    </>
  );
};
