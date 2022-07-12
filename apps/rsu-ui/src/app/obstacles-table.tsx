import React from 'react';
import { useAppSelector } from '../store';
import { TableCard } from './table-card';

export const ObstaclesTable = () => {
  const table = useAppSelector((state) => state.rsu.obstacles);
  return (
    <TableCard
      table={table}
      headers={['Longitude', 'Latitude']}
      keys={['lng', 'lat']}
      title="Obstacles Table"
    />
  );
};
