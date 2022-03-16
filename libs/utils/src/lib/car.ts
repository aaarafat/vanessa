import mapboxgl from 'mapbox-gl';
import { distanceInKm } from './distance';
import { interpolateString } from './string-utils';
import { Coordinates, ICar, CarProps } from './types';

const MS_IN_HOUR = 1000 * 60 * 60;

const carDefaultProps: CarProps = {
  title: 'Car',
  description: `<ul class="popup">
    <li>id: {id}</li>
    <li>speed: {speed} km/h</li>
  </ul>`,
};

/**
 * Car Class
 */
export class Car implements ICar {
  public id: number;
  public lat: number;
  public lng: number;
  public speed: number;
  public route: Coordinates[];
  public originalDirections: GeoJSON.Feature;
  public sourceId: string;
  private routeIndex: number;
  private source: mapboxgl.GeoJSONSource | undefined;
  private layer: mapboxgl.CircleLayer | undefined;
  private directionsSource: mapboxgl.GeoJSONSource | undefined;
  private directionsLayer: mapboxgl.LineLayer | undefined;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private updateIntervalId: any;
  private prevTime: number;
  private map: mapboxgl.Map;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private handlers: Record<string, ((...arg: any) => void)[]>;

  constructor(car: Partial<ICar> & { map: mapboxgl.Map }) {
    this.id = car.id || Date.now();
    this.sourceId = `car-${this.id}`;
    this.lat = car.lat || 0;
    this.lng = car.lng || 0;
    this.speed = car.speed || 10;
    this.route = car.route || [];
    this.map = car.map;
    this.originalDirections = car.originalDirections || {
      type: 'Feature',
      geometry: {
        type: 'Point',
        coordinates: [this.lng, this.lat],
      },
      properties: null,
    };
    this.routeIndex = 0;
    this.handlers = {};

    this.draw();
    this.attachHandlers();

    this.prevTime = Date.now();
    this.update();
  }

  public get coordinates(): Coordinates {
    return { lat: this.lat, lng: this.lng };
  }

  public get arrived(): boolean {
    return this.routeIndex === this.route.length;
  }

  public draw(): void {
    const geojson: mapboxgl.GeoJSONSourceRaw = {
      type: 'geojson',
      data: {
        type: 'Feature',
        geometry: {
          type: 'Point',
          coordinates: [this.lng, this.lat],
        },
        properties: {
          ...this.props,
        },
      },
    };

    this.source = this.map
      .addSource(this.sourceId, geojson)
      .getSource(this.sourceId) as mapboxgl.GeoJSONSource;

    this.layer = this.map
      .addLayer({
        id: this.sourceId,
        source: this.sourceId,
        type: 'circle',
        paint: {
          'circle-radius': 10,
          'circle-color': '#007cbf',
        },
      })
      .getLayer(this.sourceId) as mapboxgl.CircleLayer;

    this.directionsSource = this.map
      .addSource(`car-${this.id}-route`, {
        type: 'geojson',
        data: {
          type: 'FeatureCollection',
          features: [this.originalDirections],
        },
      })
      .getSource(`car-${this.id}-route`) as mapboxgl.GeoJSONSource;

    this.directionsLayer = this.map
      .addLayer(
        {
          id: `car-${this.id}-route`,
          type: 'line',
          source: `car-${this.id}-route`,
          layout: {
            'line-cap': 'round',
            'line-join': 'round',
            visibility: 'none',
          },
          paint: {
            'line-color': '#807515',
            'line-width': 12,
          },
        },
        'first-layer'
      )
      .getLayer(`car-${this.id}-route`) as mapboxgl.LineLayer;
  }

  private get props(): CarProps {
    return {
      ...carDefaultProps,
      id: this.id,
      lat: this.lat,
      lng: this.lng,
      speed: this.speed,
      route: this.route,
    };
  }

  private display = () => {
    console.log(
      `ID: ${this.id}, Lat: ${this.lat}, Lng: ${this.lng}, RouteIndex: ${
        this.routeIndex
      }, Arrived: ${this.arrived ? 'Yes' : 'No'}`
    );
  };

  /**
   * Update Car
   */
  private update = () => {
    this.updateCoordinates();
    this.updateSource();
    if (!this.arrived) requestAnimationFrame(this.update);
    else this.speed = 0;
  };

  private updateCoordinates = () => {
    const now = Date.now();
    let movementAmount =
      this.speed * (((now - this.prevTime) * 1.0) / MS_IN_HOUR);
    this.prevTime = now;
    while (movementAmount && !this.arrived) {
      const dist = distanceInKm(this.coordinates, this.route[this.routeIndex]);

      if (movementAmount >= dist) {
        movementAmount -= dist;
        this.lat = this.route[this.routeIndex].lat;
        this.lng = this.route[this.routeIndex].lng;
        this.routeIndex++;
        if (this.routeIndex === this.route.length) {
          clearInterval(this.updateIntervalId);
        }
      } else {
        const vector: Coordinates = {
          lng: (this.route[this.routeIndex].lng - this.coordinates.lng) / dist,
          lat: (this.route[this.routeIndex].lat - this.coordinates.lat) / dist,
        };
        this.lat += movementAmount * vector.lat;
        this.lng += movementAmount * vector.lng;
        movementAmount = 0;
      }
    }
  };

  private updateSource = () => {
    this.source?.setData({
      type: 'Feature',
      geometry: {
        type: 'Point',
        coordinates: [this.lng, this.lat],
      },
      properties: this.props,
    });
  };

  private attachHandlers = () => {
    this.map.on('click', this.sourceId, this.onClick);
  };

  private onClick = (
    e: mapboxgl.MapMouseEvent & {
      features?: mapboxgl.MapboxGeoJSONFeature[] | undefined;
    }
  ) => {
    const coordinates = (
      e.features?.[0].geometry as GeoJSON.Point
    ).coordinates.slice();

    while (Math.abs(e.lngLat.lng - coordinates[0]) > 180) {
      coordinates[0] += e.lngLat.lng > coordinates[0] ? 360 : -360;
    }

    const popup = new mapboxgl.Popup({
      closeButton: false,
    })
      .setLngLat(coordinates as mapboxgl.LngLatLike)
      .setHTML(this.description)
      .addTo(this.map);

    this.map?.setLayoutProperty(
      `car-${this.id}-route`,
      'visibility',
      'visible'
    );
    const movePopup = () => {
      if (!popup.isOpen()) {
        this.map?.setLayoutProperty(
          `car-${this.id}-route`,
          'visibility',
          'none'
        );
        return;
      }

      popup
        .setLngLat(this.coordinates as mapboxgl.LngLatLike)
        .setHTML(this.description);
      requestAnimationFrame(movePopup);
    };
    movePopup();
    this.emit('click', this);
  };

  private get description() {
    const description =
      '<h1 class="mapboxgl-popup-title">Car</h1>' + this.props.description;
    return interpolateString(description, this);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public on(type: 'click', handler: any) {
    this.subscribe(type, handler);
    return this;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private subscribe(type: 'click' | 'move', handler: (...args: any) => void) {
    if (!this.handlers[type]) this.handlers[type] = [];
    this.handlers[type].push(handler);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private emit(type: 'click' | 'move', ...args: any[]) {
    if (!this.handlers[type]) return;

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    this.handlers[type].forEach((handler: (...args: any[]) => void) =>
      handler.call(this, ...args)
    );
  }
}

export default Car;
