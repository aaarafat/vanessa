import { StrictMode } from 'react';
import * as ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { MapProvider } from '@vanessa/map';
import { SocketProvider } from './context';

import App from './app/app';

ReactDOM.render(
  <StrictMode>
    <BrowserRouter>
      <SocketProvider>
        <MapProvider>
          <App />
        </MapProvider>
      </SocketProvider>
    </BrowserRouter>
  </StrictMode>,
  document.getElementById('root')
);
