import mapboxgl, { GeoJSONSourceRaw } from "mapbox-gl";
import { Car } from "../map/src/lib/map.props";


const carDefaultProps = {
  title: 'Car',
  'marker-size': 'large',
  'marker-color': '#f00',
}

export function drawNewCar(map: mapboxgl.Map, sourceId: string, car: Car): void {
  const geojson: GeoJSONSourceRaw = {
    type: "geojson",
    data: getCarLocation(car)
  };
  map.addSource(sourceId, geojson);

  map.addLayer({
    id: sourceId,
    source: sourceId,
    type: 'circle',
    'paint': {
      'circle-radius': 10,
      'circle-color': '#007cbf',
    },
  });
}

export function updateCar(map: mapboxgl.Map, sourceId: string, car: Car) {
  (map.getSource(sourceId) as mapboxgl.GeoJSONSource).setData(getCarLocation(car));
}


export function getCarLocation(car: Car): GeoJSON.FeatureCollection {
  const { lng, lat } = car;
  return {
    type: "FeatureCollection",
    features: [
      {
        type: "Feature",
        geometry: {
          type: "Point",
          coordinates: [lng, lat],
        },
        properties: { ...carDefaultProps, ...car },
      }
    ],
  }
};
