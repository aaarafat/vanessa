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
  font-size: 1.5rem;
  margin-bottom: 1rem;
`;

const Button = styled.button`
  display: inline-block;
  margin-top: 1rem;
  align-self: center;
  color: #fff;
  border-radius: 5px;
  padding: 0.5rem;
  background-color: #2a51ff;
  width: max(30%, 150px);
  &:hover {
    cursor: pointer;
    background-color: #2543ca;
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
        <p>It can be that this port was not given to a car yet.</p>
        <ButtonsContainer>
          <Button onClick={() => history.replace('/')}>Change Port</Button>
          <Button onClick={connectCar}>Try Again</Button>
        </ButtonsContainer>
      </AlertWrapper>
    </Alert>
  );
};
