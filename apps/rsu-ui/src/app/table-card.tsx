import React from 'react';
import styled from 'styled-components';
import { Table } from './table';

const SCard = styled.div`
  width: 80%;
  margin: 1rem auto;
  background-color: rgba(255, 255, 255, 1);
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

type Props = {
  table: { [key: string]: string }[];
  headers: string[];
  title: string;
};
export const TableCard = (props: Props) => {
  return (
    <SCard>
      <STitle>{props.title}</STitle>
      <Table table={props.table} headers={props.headers} />
    </SCard>
  );
};
