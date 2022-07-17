import {
  Car,
  Coordinates,
  createFeaturePoint,
  ICar,
  RSU,
} from '@vanessa/utils';
import React, { createContext, useEffect } from 'react';
import io, { Socket } from 'socket.io-client';
import { useAppSelector } from '../app/store';

type RerouteEvent = {
  type: 'reroute';
  id: number;
  data: {
    obstacle_coordinates: Coordinates[];
  };
};

type ChangeSpeedEvent = {
  type: 'change-speed';
  id: number;
  data: {
    speed: number;
  };
};

type CheckRouteEvent = {
  type: 'check-route';
  id: number;
  data: {
    coordinates: Coordinates;
  };
};

export const socket = io('http://127.0.0.1:65432');
export const SocketContext = createContext<Socket | null>(null);
export const SocketProvider: React.FC<React.ReactNode> = ({ children }) => {
  const cars = useAppSelector((state) => state.simulation.cars);

  useEffect(() => {
    socket.on('reroute', (message: RerouteEvent) => {
      socket.receiveBuffer = socket.receiveBuffer.filter(
        ({ data: [type, data] }: any) =>
          !(data.id === message.id && type === 'reroute')
      );
      const car = cars.find((c) => c.id === message.id);
      car
        ?.updateRoute(message.data.obstacle_coordinates)
        .then((res: boolean) => res && socketEvents.addCar(car));
    });

    socket.on('change-speed', (message: ChangeSpeedEvent) => {
      console.log('change-speed', message);

      const car = cars.find((c) => c.id === message.id);
      car?.setSpeed(message.data.speed);
      // socket.emit('change-speed', message);
    });

    socket.on('check-route', (message: CheckRouteEvent) => {
      console.log('check-route', message);
      const car = cars.find((c) => c.id === message.id);
      if (!car) return;
      const coordinates = message.data.coordinates;
      const isInRoute = car?.checkObstaclesOnRoute(
        [createFeaturePoint(coordinates)],
        true
      );
      console.log(isInRoute);
      socketEvents.sendCheckRouteResponse(car, coordinates, isInRoute);
    });

    socket.on('move', (message: any) => {
      cars.find((c) => c.id === message.id)?.startMovement();
    });

    return () => {
      socket.off('reroute');
      socket.off('change-speed');
      socket.off('check-route');
      socket.off('move');
    };
  }, [cars]);

  return (
    <SocketContext.Provider value={socket}>{children}</SocketContext.Provider>
  );
};

export const socketEvents = {
  addCar: (car: Car) => {
    const message = {
      id: car.id,
      coordinates: car.coordinates,
      port: car.port,
      route: car.route,
      speed: car.carSpeed,
      stopped: car.manualStop,
    };
    // console.log('add', message);
    socket.emit('add-car', message);
  },
  updateCar: (car: Omit<ICar, 'map'>) => {
    const message = {
      id: car.id,
      coordinates: {
        lat: car.lat,
        lng: car.lng,
      },
      port: car.port,
      route: car.route,
      speed: car.speed,
      stopped: car.stopped,
    };
    // console.log('update', message);
    socket.emit('add-car', message);
  },
  changeStop: (car: Car) => {
    const message = {
      id: car.id,
      stop: car.manualStop,
    };
    socket.emit('change-stop', message);
  },
  destinationReached: (car: Car) => {
    const message = {
      id: car.id,
      coordinates: car.coordinates,
    };
    socket.emit('destination-reached', message);
  },
  obstacleDetected: (car: Car, obstacle: Coordinates) => {
    console.log(obstacle);
    const message = {
      id: car.id,
      coordinates: car.coordinates,
      obstacle_coordinates: obstacle,
    };
    socket.emit('obstacle-detected', message);
  },
  sendCarLocation: (car: Car) => {
    const message = {
      id: car.id,
      coordinates: car.coordinates,
    };
    socket.sendBuffer = socket.sendBuffer.filter(
      ({ data: [type, data] }: any) =>
        !(data.id === message.id && type === 'update-location')
    );
    socket.emit('update-location', message);
  },
  sendCheckRouteResponse: (
    car: Car,
    coordinates: Coordinates,
    isInRoute: boolean
  ) => {
    const message = {
      id: car.id,
      coordinates,
      in_route: isInRoute,
    };
    socket.emit('check-route-response', message);
  },
  addRSU: (rsu: RSU) => {
    const message = {
      id: rsu.id,
      coordinates: rsu.coordinates,
      range: rsu.range,
      port: rsu.port,
    };
    socket.emit('add-rsu', message);
  },
  clear: () => {
    socket.emit('clear');
  },
};
