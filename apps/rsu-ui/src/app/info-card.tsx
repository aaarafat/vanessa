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
`;
const STitle = styled.div``;
const SInfo = styled.div``;

export const InfoCard = (props: Props) => {
  return (
    <SCard color={props.color} bgc={props.bgc}>
      <STitle>{props.title}</STitle>
      <SInfo>{props.info}</SInfo>
    </SCard>
  );
};
