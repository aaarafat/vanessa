import React from 'react';
import styled from 'styled-components';

const STable = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  margin: 0 auto;
  font-family: 'Lato', sans-serif;
`;

const SHeader = styled.div`
  display: flex;
  color: #ffffff;
  justify-content: space-evenly;
  width: 100%;
  margin: 0rem auto;
  margin-bottom: 2px;
`;
const SHeaderCell = styled.div`
  display: flex;
  align-items: center;
  height: 2rem;
  font-size: 1.5rem;
  font-weight: bold;
  background-color: #ffc000;
  font-family: 'Bebas Neue', cursive;
  padding: 0.5rem;
  margin-right: 1px;
  width: 100%;
  &:first-child {
  }
  &:last-child {
    margin-right: 0;
  }
`;
const SRow = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-evenly;
  margin: 0rem auto;

  margin-bottom: 1px;
  &:last-child {
    margin-bottom: 0px;
  }
`;
const SCell = styled.div`
  width: 100%;
  padding: 1rem 0rem;
  padding-left: 0.75rem;
  background-color: #f2f2f2;
  margin-right: 1px;
  &:last-child {
    margin-right: 0;
  }
`;
type Props = {
  table: { [key: string]: string | number }[];
  headers: string[];
  keys: string[];
};

export const Table = (props: Props) => {
  return (
    <STable>
      <SHeader>
        {props.headers.map((header, index) => {
          return <SHeaderCell key={index}>{header}</SHeaderCell>;
        })}
      </SHeader>
      {props.table.map((row, index) => {
        return (
          <SRow key={index}>
            {props.keys.map((key, index) => {
              return <SCell key={index}>{row[key]}</SCell>;
            })}
          </SRow>
        );
      })}
    </STable>
  );
};
