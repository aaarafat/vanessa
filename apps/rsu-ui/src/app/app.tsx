import styled from 'styled-components';
import { Route, Link } from 'react-router-dom';
import { Interface } from './interface';

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
      <Route path="/:port" exact render={() => <Interface />} />
    </StyledApp>
  );
}

export default App;
