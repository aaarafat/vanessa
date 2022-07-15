export const MS_IN_HOUR = 1000 * 60 * 60;
export const FPS = 30;

export const directionsAPI =
  'https://api.mapbox.com/directions/v5/mapbox/driving/';
export const directionsAPIParams = {
  access_token:
    'pk.eyJ1IjoibWFwYm94IiwiYSI6ImNpejY4M29iazA2Z2gycXA4N2pmbDZmangifQ.-g_vE53SD2WrJ6tFX7QHmA',
  geometries: 'geojson',
  alternatives: 'false',
  overview: 'full',
  language: 'en',
  steps: 'false',
};

export const CAR_PORT_INIT = 10000;
export const RSU_PORT_INIT = 5000;
