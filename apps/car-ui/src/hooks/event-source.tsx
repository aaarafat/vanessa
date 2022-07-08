import { Coordinates, createFeaturePoint } from '@vanessa/utils';
import React, { createContext, useContext, useEffect, useState } from 'react';
import { useAppDispatch, useAppSelector } from '../store';
import { addMessage, addObstacle, addObstacles } from '../store/carSlice';
import * as turf from '@turf/turf';

interface DestinationReachedData {
  coordinates: Coordinates;
}

interface ObstacleDetectedData {
  coordinates: Coordinates;
  obstacle_coordinates: Coordinates;
}

interface UpdateLocationData {
  coordinates: Coordinates;
}

interface ObstacleReceivedData {
  obstacle_coordinates: Coordinates[];
}

export const useEventSource = () => {
  const [eventSource, setEventSource] = useState<EventSource | null>(null);
  const { car } = useAppSelector((state) => state.car);
  const dispatch = useAppDispatch();

  useEffect(() => {
    if (!eventSource) return;
    eventSource.addEventListener('destination-reached', ({ data: message }) => {
      const json: DestinationReachedData = JSON.parse(message).data;
      car?.updateDestinationFromData(json.coordinates);
      dispatch(addMessage(message));
    });

    eventSource.addEventListener('obstacle-detected', ({ data: message }) => {
      const json: ObstacleDetectedData = JSON.parse(message).data;
      dispatch(addObstacle(createFeaturePoint(json.obstacle_coordinates)));
      car?.updateObstacleDetectedFromData();
      dispatch(addMessage(message));
    });

    eventSource.addEventListener('update-location', ({ data: message }) => {
      const json: UpdateLocationData = JSON.parse(message).data;
      car?.updateLocationFromData(json.coordinates);
    });

    eventSource.addEventListener('obstacle-received', ({ data: message }) => {
      const json: ObstacleReceivedData = JSON.parse(message).data;
      const obstacles = json.obstacle_coordinates.map((coordinates) =>
        createFeaturePoint(coordinates)
      );
      dispatch(addObstacles(obstacles));
      dispatch(addMessage(message));
    });

    return () => eventSource.close();
  }, [eventSource, car, dispatch]);

  return [eventSource, setEventSource] as const;
};
