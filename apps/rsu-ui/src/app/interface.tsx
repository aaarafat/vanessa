import React from 'react';
import styled from 'styled-components';
import { ARP } from './ARP';
import { InfoCard } from './info-card';
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
const SCards = styled.div`
  display: flex;
  flex-wrap: wrap;
`;
export const Interface: React.FC = () => {
  return (
    <>
      <SCards>
        <InfoCard
          title="Received Packets"
          info="9999999"
          color="#ffffff"
          bgc="#ffc000"
        />
        <InfoCard
          title="Received Packets"
          info="9999999"
          color="#ffffff"
          bgc="#0d0d0d"
        />
      </SCards>
      <STables>
        <OTable>
          <ObstaclesTable />
        </OTable>
        <ARPTable>
          <ARP />
        </ARPTable>
      </STables>
    </>
  );
};
