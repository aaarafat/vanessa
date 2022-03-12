import { StrictMode } from 'react';
import * as ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { MapProvider } from '@vanessa/map';

import App from './app/app';

ReactDOM.render(
  <StrictMode>
    <BrowserRouter>
      <MapProvider>
        <App />
      </MapProvider>
    </BrowserRouter>
  </StrictMode>,
  document.getElementById('root')
);
