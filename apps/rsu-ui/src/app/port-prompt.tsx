import React, { useState } from 'react';
import styled from 'styled-components';
import { useHistory } from 'react-router-dom';
import { Alert } from './alert';

const Form = styled.form`
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
  text-align: center;
  color: #fff;
  font-size: 2rem;
`;

const Input = styled.input`
  display: inline-block;
  border: 1px solid #000;
  border-radius: 1px;
  padding: 0.5rem;
  font-size: 1rem;
  font-family: 'Lato', sans-serif;
  &:focus {
    outline: none;
  }
`;

const Button = styled.button`
  display: inline-block;
  margin-top: 2rem;
  align-self: center;
  color: #fff;
  border-radius: 1px;
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

const Error = styled.p`
  color: red;
  font-size: 1rem;
  margin-top: 0.5rem;
  text-align: center;
  font-family: 'Lato', sans-serif;
`;

export const PortPrompt: React.FC = () => {
  const [port, setPort] = useState<string>('');
  const [error, setError] = useState<string | null>(null);
  const history = useHistory();

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!isValidPort(port)) {
      setError('Please, prompt for a valid port number to continue');
      return;
    }
    history.push(`/${port}`);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setError(null);
    const p = e.target.value;
    isValidPort(p);
    setPort(p);
  };

  const isValidPort = (p: string) => {
    if (p === '' || !Number.isInteger(+p) || +p <= 0 || +p > 65535) {
      setError('Port number should be an integer between 1 and 65535');
      return false;
    }
    setError(null);
    return true;
  };

  return (
    <Alert>
      <Form onSubmit={handleSubmit}>
        <Text>Enter The RSU Port </Text>
        <Input
          type="text"
          value={port}
          onChange={handleChange}
          placeholder="For example: 3000"
        />
        {error && <Error>{error}</Error>}
        <Button type="submit">Continue</Button>
      </Form>
    </Alert>
  );
};
