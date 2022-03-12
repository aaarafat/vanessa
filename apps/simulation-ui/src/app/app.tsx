import styled from 'styled-components';

import { Route } from 'react-router-dom';
import { Car, MapContext } from '@vanessa/map';
import { useContext, useEffect } from 'react';
import { Simulation } from './Map';

const StyledApp = styled.div`
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
`;


export function App() {
  return (
    <StyledApp>
      <Route path="/" exact render={() => <Simulation />} />
      {/* END: routes */}
    </StyledApp>
  );
}

export default App;
