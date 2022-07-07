import styled from 'styled-components';

import { Route } from 'react-router-dom';
import { Simulation } from './simulation';
import { Provider } from 'react-redux';
import { store } from './store';

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
    <Provider store={store}>
      <StyledApp>
        <Route path="/:id?" exact render={() => <Simulation />} />
        {/* END: routes */}
      </StyledApp>
    </Provider>
  );
}

export default App;
