import { StrictMode } from 'react';
import * as ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { MapProvider } from '@vanessa/map';
import { EventSourceProvider } from './context';

import App from './app/app';

ReactDOM.render(
  <StrictMode>
    <BrowserRouter>
      <EventSourceProvider>
        <MapProvider>
          <App />
        </MapProvider>
      </EventSourceProvider>
    </BrowserRouter>
  </StrictMode>,
  document.getElementById('root')
);
