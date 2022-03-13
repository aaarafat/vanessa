import mapboxgl, { GeoJSONSourceRaw } from 'mapbox-gl';
import { Car } from './car';
import { CarProps } from './types';

const carDefaultProps: CarProps = {
  title: 'Car',
  description: `id: {id}<br>
    speed: {speed} km/h`,
};

export function drawNewCar(
  map: mapboxgl.Map,
  sourceId: string,
  car: Car
): void {
  const geojson: GeoJSONSourceRaw = {
    type: 'geojson',
    data: getCarLocation(car),
  };
  map.addSource(sourceId, geojson);

  map.addLayer({
    id: sourceId,
    source: sourceId,
    type: 'circle',
    paint: {
      'circle-radius': 10,
      'circle-color': '#007cbf',
    },
  });
}

export function updateCar(map: mapboxgl.Map, sourceId: string, car: Car) {
  (map.getSource(sourceId) as mapboxgl.GeoJSONSource).setData(
    getCarLocation(car)
  );
}

export function getCarLocation(car: Car): GeoJSON.FeatureCollection {
  const { lng, lat } = car;
  return {
    type: 'FeatureCollection',
    features: [
      {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [lng, lat],
        },
        properties: { ...carDefaultProps, ...car },
      },
    ],
  };
}
