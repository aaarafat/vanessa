import React, { useContext, useEffect, useState } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled, { keyframes } from 'styled-components';
import { Car, Coordinates, ICar, IRSU, RSU } from '@vanessa/utils';
import * as turf from '@turf/turf';
import mapboxgl from 'mapbox-gl';
import { EventSourceContext } from '../context';
import { useHistory, useParams } from 'react-router-dom';
import MessagesViewer from './messages-viewer';

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

let car: Car;

export const Interface: React.FC = () => {
  const { map } = useContext(MapContext);
  const [eventSource, setEventSource] = useContext(EventSourceContext);
  const [mapLoaded, setMapLoaded] = useState(false);
  const [loading, setLoading] = useState(false);
  const [obstacles, setObstacles] = useState<turf.FeatureCollection>({
    type: 'FeatureCollection',
    features: [],
  });
  const { port } = useParams<{
    port: string;
  }>();

  console.log(eventSource);

  useEffect(() => {
    if (port) {
      setEventSource(new EventSource(`http://localhost:${port}`));
    }
  }, [port, setEventSource]);

  useEffect(() => {
    if (map && eventSource) {
      map.setMaxZoom(18);
      map.on('load', () => {
        setMapLoaded(true);
      });
    }
  }, [map, eventSource]);

  // useEffect(() => {
  //   if(eventSource) {

  //   }
  // }, [eventSource])

  useEffect(() => {
    if (mapLoaded && map) {
      const obstacle = turf.buffer(obstacles, 10, { units: 'meters' });
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
