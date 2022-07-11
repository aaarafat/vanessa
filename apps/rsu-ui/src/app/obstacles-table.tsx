import React from 'react';
import styled from 'styled-components';
import { Table } from './table';

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

const SCard = styled.div`
  width: 80%;
  margin: 1rem auto;
  background-color: #fff;
  border-radius: 1px;
  box-shadow: 0 0.5rem 1rem rgba(0, 0, 0, 0.01);
  padding: 1rem;
`;
const STitle = styled.div`
  font-size: 2.5rem;
  font-weight: bold;
  margin-bottom: 1rem;
  color: #0d0d0d;
  font-family: 'Bebas Neue', cursive;
`;
export const ObstaclesTable = () => {
  return (
    <SCard>
      <STitle>Obstacles Table</STitle>
      <Table table={table} header={['Lat', 'Lon']} />
    </SCard>
  );
};