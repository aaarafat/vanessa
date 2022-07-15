import { StrictMode } from 'react';
import * as ReactDOM from 'react-dom';
import { BrowserRouter } from 'react-router-dom';
import { MapProvider } from '@vanessa/map';

import App from './app/app';
import { Provider } from 'react-redux';
import { store } from './store';
import './main.css';

ReactDOM.render(
  <StrictMode>
    <Provider store={store}>
      <BrowserRouter>
        <MapProvider>
          <App />
        </MapProvider>
      </BrowserRouter>
    </Provider>
  </StrictMode>,
  document.getElementById('root')
);
