import React, { useContext, useEffect, useState } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled, { keyframes } from 'styled-components';
import {
  Car,
  Coordinates,
  createFeaturePoint,
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

let carPortsCounter = 10000;
let rsuPortsCounter = 5000;

export const Simulation: React.FC = () => {
  const { map, mapDirections } = useContext(MapContext);
  const socket = useContext(SocketContext);
  const [mapLoaded, setMapLoaded] = useState(false);
  const [loading, setLoading] = useState(true);
  const [obstacles, setObstacles] = useState<turf.Feature[]>([]);
  const dispatch = useAppDispatch();
  const { cars, rsus, rsusData } = useAppSelector((state) => state.simulation);
  const history = useHistory();
  const { id } = useParams<{
    id: string;
  }>();

  useEffect(() => {
    if (map && socket && !mapLoaded) {
      map.setMaxZoom(18);
      map.on('load', () => {
        setMapLoaded(true);

        rsusData.forEach((r) => {
          const rsu = new RSU({ ...r, map, port: rsuPortsCounter++ });
          dispatch(addRSU(rsu));
          socketEvents.addRSU(rsu);
        });
        setLoading(false);
      });
    }
  }, [map, socket, dispatch, rsusData, mapLoaded]);

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

  useEffect(() => {
    if (id && isNaN(Number(id))) {
      history.replace('/');
      return;
    }
    if (mapLoaded) {
      if (!id) {
        dispatch(unfocusCar());
        map?.setLayoutProperty('obstacles', 'visibility', 'visible');
        mapDirections.options.interactive = true;
        return;
      }
      const index = cars.findIndex((c) => c.id === +id);
      if (index === -1) {
        history.replace('/');
        return;
      }
      dispatch(focusCar(index));
      map?.setLayoutProperty('obstacles', 'visibility', 'none');
      mapDirections.options.interactive = false;
    }
  }, [id, mapLoaded]);

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
    car.on('focus', () => {
      history.push(`/${car.id}`);
    });
    car.on('popup:close', () => {
      mapDirections.unfreeze();
    });

    car.on('move', () => socketEvents.sendCarLocation(car));
    car.on('obstacle-detected', () => socketEvents.obstacleDetected(car));
    car.on('destination-reached', () => socketEvents.destinationReached(car));

    dispatch(addCar(car));
    socketEvents.addCar(car);

    mapDirections.reset();
  };

  const handleExport = () => {
    const info: Array<any> = [...cars, ...rsus].map((item) => item.export());
    info.push({
      coordinates: obstacles
        .map((f) =>
          f.geometry.type === 'Point' ? f.geometry.coordinates : null
        )
        .filter(Boolean),
      type: 'obstacles',
    });
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

    dispatch(
      clearState({
        removeRSUs,
      })
    );
    setObstacles([]);
    mapDirections.reset();
  };

  const handleImport = () => {
    const reader = new FileReader();
    reader.addEventListener('load', () => {
      clearMap(true);
      const data = JSON.parse(reader.result as string);
      data.forEach((item: any) => {
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
      setLoading(false);
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
