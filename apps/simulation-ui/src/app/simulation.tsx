import React, { useContext, useEffect, useState } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled, { keyframes } from 'styled-components';
import {
  Car,
  Coordinates,
  createFeaturePoint,
  getObstacleFeatures,
  ICar,
  IRSU,
  RSU,
} from '@vanessa/utils';
import * as turf from '@turf/turf';
import mapboxgl from 'mapbox-gl';
import { socketEvents, SocketContext } from '../context';
import ControlPanel from './control-panel';
import { useAppDispatch, useAppSelector } from './store';
import {
  addCar,
  addRSU,
  clearState,
  focusCar,
  unfocusCar,
} from './store/simulationSlice';
import { useHistory, useParams } from 'react-router-dom';
import MessagesViewer from './messages-viewer';
import { CAR_PORT_INIT, RSU_PORT_INIT } from 'libs/utils/src/lib/constants';

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

let carPortsCounter = CAR_PORT_INIT;
let rsuPortsCounter = RSU_PORT_INIT;

export const Simulation: React.FC = () => {
  const { map, mapDirections } = useContext(MapContext);
  const socket = useContext(SocketContext);
  const [mapLoaded, setMapLoaded] = useState(false);
  const [loading, setLoading] = useState(true);
  const [obstacles, setObstacles] = useState<turf.Feature<turf.Point>[]>([]);
  const dispatch = useAppDispatch();
  const { cars, rsus, rsusData } = useAppSelector((state) => state.simulation);

  useEffect(() => {
    if (map && socket && !mapLoaded) {
      map.setMaxZoom(18);
      map.on('load', () => {
        socketEvents.clear();
        setMapLoaded(true);
        socket.once('cleared', () => {
          rsusData.forEach((r) => {
            const rsu = new RSU({ ...r, map, port: rsuPortsCounter++ });
            dispatch(addRSU(rsu));
            socketEvents.addRSU(rsu);
          });
          setLoading(false);
        });
      });
    }
  }, [map, socket, dispatch, rsusData, mapLoaded]);

  useEffect(() => {
    if (mapLoaded && map) {
      const obstacle = getObstacleFeatures(obstacles);
      const obstaclePoints = turf.featureCollection(
        obstacles.map((o) => createFeaturePoint(o.geometry.coordinates))
      );

      if (!map.getSource('obstacles')) {
        map.addSource('obstacles', {
          type: 'geojson',
          data: obstacle,
        });

        map.addSource('obstacles-points', {
          type: 'geojson',
          data: obstaclePoints,
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
        (map.getSource('obstacles-points') as mapboxgl.GeoJSONSource).setData(
          obstaclePoints
        );
      }
    }
  }, [obstacles, map, mapLoaded]);

  const handleAddObstacle = (coordinates: Coordinates | null) => {
    if (!coordinates) return;
    setObstacles((prev) => [...prev, createFeaturePoint(coordinates)]);
  };

  const handleAddCar = (carInputs: Partial<ICar>) => {
    if (!map) return;

    if (!carInputs.route) return;

    const car: Car = new Car({
      ...carInputs,
      map,
      port: carPortsCounter++,
    });

    car.on('click', () => {
      mapDirections.reset();
      mapDirections.freeze();
    });
    car.on('popup:close', () => {
      mapDirections.unfreeze();
    });

    car.on('move', () => socketEvents.sendCarLocation(car));
    car.on('obstacle-detected', (obstacle: Coordinates) =>
      socketEvents.obstacleDetected(car, obstacle)
    );
    car.on('destination-reached', () => socketEvents.destinationReached(car));

    dispatch(addCar(car));
    socketEvents.addCar(car);

    mapDirections.reset();
  };

  const handleExport = () => {
    const info: Array<any> = [
      {
        coordinates: obstacles
          .map((f) =>
            f.geometry.type === 'Point' ? f.geometry.coordinates : null
          )
          .filter(Boolean),
        type: 'obstacles',
      },
    ];
    info.push([...rsus, ...cars].map((item) => item.export()));
    const fileData = JSON.stringify(info, null, 2);
    const blob = new Blob([fileData], { type: 'text/plain' });
    const url = URL.createObjectURL(blob);

    const link = document.createElement('a');
    link.download = 'exported.json';
    link.href = url;
    link.click();
    link.remove();
  };

  const clearMap = (removeRSUs = false) => {
    if (!map) return;
    setLoading(true);
    socketEvents.clear();

    cars.forEach((car) => car.remove());
    if (removeRSUs) {
      rsus.forEach((rsu) => rsu.remove());
    } else {
      // we need to send it again as it will be cleared from the emulation
      rsus.forEach((r) => socketEvents.addRSU(r));
    }
    // clearing the arrays in the store
    dispatch(
      clearState({
        removeRSUs,
      })
    );

    setObstacles([]);
    mapDirections.reset();
    carPortsCounter = CAR_PORT_INIT;
    rsuPortsCounter = RSU_PORT_INIT;
    setLoading(false);
  };

  const handleImport = () => {
    const reader = new FileReader();
    reader.addEventListener('load', () => {
      setLoading(true);
      clearMap(true);
      socket?.once('cleared', () => {
        const data = JSON.parse(reader.result as string);
        doImport(data);
        setLoading(false);
      });
    });

    const input = document.createElement('input');
    input.type = 'file';
    input.onchange = (e) => {
      if (input.files?.[0]) {
        setLoading(true);
        reader.readAsText(input.files[0]);
      }
    };
    input.click();
    input.remove();
  };

  const doImport = (data: any) => {
    return data.forEach((item: any) => {
      if (item.type === 'car') {
        handleAddCar({
          ...item,
        });
      } else if (item.type === 'rsu') {
        const rsu = new RSU({ ...item, map });
        dispatch(addRSU(rsu));
        socketEvents.addRSU(rsu);
      } else if (item.type === 'obstacles') {
        setObstacles(
          item.coordinates.map((c: turf.Position) => createFeaturePoint(c))
        );
      }
    });
  };

  return (
    <>
      {loading && (
        <LoaderContainer>
          <Loader />
        </LoaderContainer>
      )}
      <div>
        <Map />
        <ControlPanel
          onAddCar={handleAddCar}
          onAddObstacle={handleAddObstacle}
          onExport={handleExport}
          onImport={handleImport}
          onClearMap={() => clearMap(false)}
        />
        <MessagesViewer />
      </div>
    </>
  );
};
