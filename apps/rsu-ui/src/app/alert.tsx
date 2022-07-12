import React from 'react';
import styled from 'styled-components';

const Container = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  z-index: 9999;
  background-color: rgba(0, 0, 0, 0.35);

  backdrop-filter: blur(5px);
  width: 100vw;
  height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
`;

const Wrapper = styled.div`
  display: flex;
  border: 1px solid #000;
  border-radius: 5px;
  background-color: #fff;
  width: max(350px, 30%);
  box-shadow: 0px 0px 10px rgba(0, 0, 0, 0.5);
`;

export const Alert: React.FC = ({ children }) => {
  return (
    <Container>
      <Wrapper>{children}</Wrapper>
    </Container>
  );
};
