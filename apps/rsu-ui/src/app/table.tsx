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
  border-bottom: 3px solid #ffffff;
  margin: 0rem auto;
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
  border-right: 1px solid #ffffff;
  width: 100%;
  &:first-child {
  }
  &:last-child {
  }
`;
const SRow = styled.div`
  display: flex;
  width: 100%;
  justify-content: space-evenly;
  margin: 0rem auto;
  background-color: #f2f2f2;
  border-bottom: 1px solid #ffffff;
  &:last-child {
    border: 0px;
  }
`;
const SCell = styled.div`
  width: 100%;
  padding: 1rem 0rem;
  border-right: 1px solid #ffffff;
  padding-left: 0.75rem;
  &:first-child {
  }
`;
type Props = {
  table: { [key: string]: string }[];
  headers: string[];
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
            {Object.keys(row).map((key, index) => {
              return <SCell key={index}>{row[key]}</SCell>;
            })}
          </SRow>
        );
      })}
    </STable>
  );
};
