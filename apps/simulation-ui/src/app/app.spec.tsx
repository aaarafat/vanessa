import { render } from '@testing-library/react';
import { Provider } from 'react-redux';

import { BrowserRouter } from 'react-router-dom';

import App from './app';
import { store } from './store';

jest.mock('mapbox-gl/dist/mapbox-gl', () => ({
  GeolocateControl: jest.fn(),
  Map: jest.fn(() => ({
    addControl: jest.fn(),
    on: jest.fn(),
    remove: jest.fn(),
  })),
  NavigationControl: jest.fn(),
}));

describe('App', () => {
  it('should render successfully', () => {
    const { baseElement } = render(
      <BrowserRouter>
        <Provider store={store}>
          <App />
        </Provider>
      </BrowserRouter>
    );

    expect(baseElement).toBeTruthy();
  });
});
