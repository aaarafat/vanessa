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
const SHeading = styled.h1`
  font-size: 5rem;
  font-family: 'Bebas Neue', cursive;
  margin: 1rem;
  color: #0d0d0d;
`;
const SDiscription = styled.h3`
  font-size: 2.5rem;
  font-family: 'Bebas Neue', cursive;
  margin: 0 1rem;
  color: #ffc000;
`;
export const Interface: React.FC = () => {
  return (
    <>
      <SHeading>RSU Overview</SHeading>
      <SDiscription>Total packets</SDiscription>
      <SCards>
        <InfoCard
          title="Sent to RSUs"
          info="9999999"
          color="#ffffff"
          bgc="#0d0d0d"
        />
        <InfoCard
          title="Recieved from RSUs"
          info="9999999"
          color="#ffffff"
          bgc="#ffc000"
        />
        <InfoCard
          title="Sent to cars"
          info="9999999"
          color="#ffffff"
          bgc="#2d2d2d"
        />
        <InfoCard
          title="Recieved from cars"
          info="9999999"
          color="#ffffff"
          bgc="#ffcf33"
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
