import React from 'react';
import styled from 'styled-components';

type Props = { title: string; info: string; color: string; bgc: string };

interface ICard {
  color: string;
  bgc: string;
}

const SCard = styled.div<ICard>`
  background-color: ${(p) => p.bgc};
  color: ${(p) => p.color};
  padding: 1rem;
  width: 200px;
  margin: 1rem;
`;
const STitle = styled.div`
  font-family: 'Lato', sans-serif;
  margin-bottom: 1rem;
`;
const SInfo = styled.div`
  font-family: 'Bebas Neue', cursive;
  font-size: 2.5rem;
`;

export const InfoCard = (props: Props) => {
  return (
    <SCard color={props.color} bgc={props.bgc}>
      <STitle>{props.title}</STitle>
      <SInfo>{props.info}</SInfo>
    </SCard>
  );
};
