import React from 'react';
import styled from 'styled-components';
import { useHistory } from 'react-router-dom';
import { Alert } from './alert';

const AlertWrapper = styled.div`
  display: flex;
  flex: 1;
  flex-direction: column;
  justify-content: center;
  padding: 2rem;
  margin: 1rem;
`;

const Text = styled.h1`
  font-size: 2rem;
  margin-bottom: 1rem;
  color: #fff;
  text-align: center;
`;

const Paragraph = styled.p`
  font-size: 1rem;
  margin-bottom: 1rem;
  color: #fff;
  text-align: center;
  font-family: 'Lato', sans-serif;
`;

const Button = styled.button`
  display: inline-block;
  margin-top: 1rem;
  align-self: center;
  color: #fff;
  border-radius: 5px;
  padding: 0.5rem;
  background-color: #ffc000;
  width: max(30%, 150px);
  font-family: 'Bebas Neue', cursive;
  font-size: 1.25rem;
  border: none;
  &:hover {
    cursor: pointer;
    background-color: #ffcf33;
  }
`;

const ButtonsContainer = styled.div`
  display: flex;
  flex: 1;
  flex-direction: row;
  justify-content: center;
  gap: 0.5rem;
`;

export const ConnectionErrorAlert: React.FC<{
  connectionError: boolean;
  connectCar: () => void;
}> = ({ connectionError, connectCar }) => {
  const history = useHistory();

  if (!connectionError) return null;
  return (
    <Alert>
      <AlertWrapper>
        <Text>
          Connection Failed.
          <br />
        </Text>
        <Paragraph>
          It can be that this port was not given to a car yet.
        </Paragraph>
        <ButtonsContainer>
          <Button onClick={() => history.replace('/')}>Change Port</Button>
          <Button onClick={connectCar}>Try Again</Button>
        </ButtonsContainer>
      </AlertWrapper>
    </Alert>
  );
};
