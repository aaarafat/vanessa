import { Car, Coordinates, RSU } from '@vanessa/utils';
import React, { createContext } from 'react';
import io, { Socket } from 'socket.io-client';

export const socket = io('http://127.0.0.1:65432');
export const SocketContext = createContext<Socket | null>(null);
export const SocketProvider: React.FC<React.ReactNode> = ({ children }) => (
  <SocketContext.Provider value={socket}>{children}</SocketContext.Provider>
);

export const socketEvents = {
  addCar: (car: Car) => {
    const message = {
      id: car.id,
      coordinates: car.coordinates,
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
  obstacleDetected: (car: Car) => {
    const message = {
      id: car.id,
      coordinates: car.coordinates,
      obstacle_coordinates: car.coordinates,
    };
    socket.emit('obstacle-detected', message);
  },
  sendCarLocation: (car: Car) => {
    const message = {
      id: car.id,
      coordinates: car.coordinates,
    };
    socket.volatile.emit('update-location', message);
  },
  addRSU: (rsu: RSU) => {
    const message = {
      id: rsu.id,
      coordinates: rsu.coordinates,
      range: rsu.radius,
    };
    socket.emit('add-rsu', message);
  },
};
