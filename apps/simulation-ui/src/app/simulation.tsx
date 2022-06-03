import React, { useContext, useEffect, useState } from 'react';
import { Map, MapContext } from '@vanessa/map';
import styled, { keyframes } from 'styled-components';
import { Car, Coordinates, ICar, IRSU, RSU } from '@vanessa/utils';
import * as turf from '@turf/turf';
import mapboxgl from 'mapbox-gl';
import { SocketContext } from '../context';
import { CLICK_SOURCE_ID } from './constants';
import ControlPanel from './control-panel';

const rsusData: Partial<IRSU>[] = [
  {
    id: 1,
    lng: 31.213,
    lat: 30.0252,
    radius: 0.25,
  },
  {
    id: 2,
    lng: 31.2029,
    lat: 30.0269,
    radius: 0.5,
  },
  {
    id: 3,
    lng: 31.2129,
    lat: 30.0185,
    radius: 0.5,
  },
];

const rsus: RSU[] = [];

const cars: Car[] = [];

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

export const Simulation: React.FC = () => {
  const { map, mapDirections } = useContext(MapContext);
  const socket = useContext(SocketContext);
  const [mapLoaded, setMapLoaded] = useState(false);
  const [loading, setLoading] = useState(false);
  const [obstacles, setObstacles] = useState<turf.FeatureCollection>({
    type: 'FeatureCollection',
    features: [],
  });

  useEffect(() => {
    if (map && socket) {
      map.setMaxZoom(18);
      map.on('load', () => {
        rsusData.forEach((rsu) => rsus.push(new RSU({ ...rsu, map, socket })));
        setMapLoaded(true);
      });
    }
  }, [map, socket]);

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

  const handleAddObstacle = (coordinates: Coordinates | null) => {
    if (!coordinates) return;
    setObstacles((prev) => ({
      ...prev,
      features: [
        ...prev.features,
        {
          type: 'Feature',
          geometry: {
            type: 'Point',
            coordinates: [coordinates.lng, coordinates.lat],
          },
          properties: {},
        },
      ],
    }));
  };

  const handleAddCar = (carInputs: Partial<ICar>) => {
    if (!map) return;

    // if no previous click, return
    const source = map.getSource(CLICK_SOURCE_ID);
    if (!source || !carInputs.route) return;

    const car: Car = new Car({
      ...carInputs,
      map,
      socket,
    });

    // add to cars list
    cars.push(car);
    car.on('click', () => {
      mapDirections.reset();
      mapDirections.freeze();
    });
    car.on('popup-closed', () => {
      mapDirections.unfreeze();
    });

    mapDirections.reset();
  };

  const handleExport = () => {
    const info: Array<any> = [...cars, ...rsus].map((item) => item.export());
    info.push({
      coordinates: obstacles.features
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
    (map.getSource('obstacles') as mapboxgl.GeoJSONSource).setData({
      type: 'FeatureCollection',
      features: [],
    });
    cars.forEach((car) => car.remove());
    cars.splice(0, cars.length);
    if (removeRSUs) {
      rsus.forEach((rsu) => rsu.remove());
      rsus.splice(0, rsus.length);
    }
    mapDirections.reset();
  };

  const handleImport = () => {
    const input = document.createElement('input');
    input.type = 'file';
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
          rsus.push(new RSU({ ...item, map, socket }));
        } else if (item.type === 'obstacles') {
          setObstacles({
            type: 'FeatureCollection',
            features: item.coordinates.map((c: turf.Point) => ({
              type: 'Feature',
              geometry: {
                type: 'Point',
                coordinates: c,
              },
              properties: {},
            })),
          });
        }
      });
      setLoading(false);
    });
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
      </div>
    </>
  );
};
