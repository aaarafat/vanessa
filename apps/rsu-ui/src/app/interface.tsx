import React from 'react';
import styled from 'styled-components';
import { ARP } from './ARP';
import { ObstaclesTable } from './obstacles-table';

const STables = styled.div`
  display: flex;
`;
const OTable = styled.div`
  width: 30%;
`;
const ARPTable = styled.div`
  width: 70%;
`;
export const Interface: React.FC = () => {
  return (
    <STables>
      <OTable>
        <ObstaclesTable />
      </OTable>
      <ARPTable>
        <ARP />
      </ARPTable>
    </STables>
  );
};
