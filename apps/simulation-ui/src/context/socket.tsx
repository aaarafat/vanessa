import { Car, Coordinates, RSU } from '@vanessa/utils';
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

    socket.on('check-route', (message: any) => {
      console.log('check-route', message);
    });

    socket.on('move', (message: any) => {
      console.log('move', message);
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
      speed: car.speed,
    };
    socket.emit('add-car', message);
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
  addRSU: (rsu: RSU) => {
    const message = {
      id: rsu.id,
      coordinates: rsu.coordinates,
      range: rsu.radius * 1000,
      port: rsu.port,
    };
    socket.emit('add-rsu', message);
  },
  clear: () => {
    socket.emit('clear');
  },
};
