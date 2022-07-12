import React from 'react';
import { TableCard } from './table-card';

const table = [
  {
    Lat: '52.5',
    Lon: '13.4',
  },
  {
    Lat: '52.5',
    Lon: '13.4',
  },
  {
    Lat: '52.5',
    Lon: '13.4',
  },
  {
    Lat: '52.5',
    Lon: '13.4',
  },
  {
    Lat: '52.5',
    Lon: '13.4',
  },
  {
    Lat: '52.5',
    Lon: '13.4',
  },
];

export const ObstaclesTable = () => {
  return (
    <TableCard table={table} headers={['Lat', 'Lon']} title="Obstacles Table" />
  );
};
