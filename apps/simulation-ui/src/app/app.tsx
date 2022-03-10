import styled from 'styled-components';

import { Route } from 'react-router-dom';
import { Car, Map } from '@vanessa/map';

const StyledApp = styled.div`
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
    'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
    sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
`;

const cars: Car[] = [
  {
    id: 1,
    lat: 30.02543,
    lng: 31.21146,
  },
  {
    id: 2,
    lat: 30.02763,
    lng: 31.21082,
  },
  {
    id: 3,
    lat: 30.02425,
    lng: 31.20995,
  }
]

export function App() {
  return (
    <StyledApp>
      <Route path="/" exact render={() => <Map cars={cars} />} />
      {/* END: routes */}
    </StyledApp>
  );
}

export default App;
