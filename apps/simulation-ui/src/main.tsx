import { StrictMode } from 'react';
import * as ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { MapProvider } from '@vanessa/map';
import { SocketProvider } from './context';

import App from './app/app';
import { Provider } from 'react-redux';
import { store } from './app/store';

ReactDOM.render(
  <StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <SocketProvider>
          <MapProvider>
            <App />
          </MapProvider>
        </SocketProvider>
      </BrowserRouter>
    </Provider>
  </StrictMode>,
  document.getElementById('root')
);
