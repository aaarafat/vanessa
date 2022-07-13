import styled from 'styled-components';
import { Route } from 'react-router-dom';
import { Interface } from './interface';
import { PortPrompt } from './port-prompt';

const StyledApp = styled.div`
  margin: 0;
  font-family: 'Bebas Neue', cursive;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
`;

export function App() {
  return (
    <StyledApp>
      <Route path="/" exact render={() => <PortPrompt />} />
      <Route path="/:port" exact render={() => <Interface />} />
    </StyledApp>
  );
}

export default App;
