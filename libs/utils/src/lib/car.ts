import mapboxgl from 'mapbox-gl';
import { distanceInKm, euclideanDistance } from './distance';
import { interpolateString } from './string-utils';
import { Coordinates, ICar, CarProps, PartialExcept } from './types';
import { MS_IN_HOUR, FPS } from './constants';
import * as turf from '@turf/turf';
import { Socket } from 'socket.io-client';

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
export class Car {
  public id: number;
  public lat: number;
  public lng: number;
  private prevCoordinates: Coordinates;
  public speed: number;
  public route: Coordinates[];
  public originalDirections: GeoJSON.Feature;
  public sourceId: string;
  private routeIndex: number;
  private source: mapboxgl.GeoJSONSource | undefined;
  private layer: mapboxgl.CircleLayer | undefined;
  private directionsSource: mapboxgl.GeoJSONSource | undefined;
  private directionsLayer: mapboxgl.LineLayer | undefined;
  private communicationRangeSource: mapboxgl.GeoJSONSource | undefined;
  private communicationRangeLayer: mapboxgl.FillLayer | undefined;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private prevTime: number;
  private map: mapboxgl.Map;
  private socket: Socket;
  private popup: mapboxgl.Popup | null;
  private obstacleDetected: boolean;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private handlers: Record<string, ((...arg: any) => void)[]>;
  private wasFlyingToCar: boolean;
  private animationFrame: number;
  private removed = false;
  private focused: boolean;

  constructor(car: PartialExcept<ICar, 'map' | 'socket'>) {
    this.id = car.id || Date.now();
    this.sourceId = `car-${this.id}`;
    this.lat = car.lat || 0;
    this.lng = car.lng || 0;
    this.prevCoordinates = this.coordinates;
    this.speed = car.speed ?? 10;
    this.route = car.route || [];
    this.map = car.map;
    this.socket = car.socket;
    this.popup = null;
    this.focused = false;
    this.originalDirections = turf.lineString(
      this.route.map((r: Coordinates) => [r.lng, r.lat])
    );
    this.routeIndex = 0;
    this.handlers = {};
    this.wasFlyingToCar = false;
    this.animationFrame = 0;
    this.obstacleDetected = car.obstacleDetected || false;

    this.socket.emit('add-car', {
      id: this.id,
      coordinates: this.coordinates,
    });

    // TODO: REMOVE THIS!!!!!!!!!!!!!!!!!!!!
    this.socket.on('change', (data: any) => {
      console.log(data);
    });

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

    this.communicationRangeSource = this.map
      .addSource(`car-${this.id}-com-range`, {
        type: 'geojson',
        data: this.communicationRangeFeature,
      })
      .getSource(`car-${this.id}-com-range`) as mapboxgl.GeoJSONSource;

    this.communicationRangeLayer = this.map
      .addLayer({
        id: `car-${this.id}-com-range`,
        type: 'fill',
        source: `car-${this.id}-com-range`,
        layout: {
          visibility: 'none',
        },
        paint: {
          'fill-color': '#f03b20',
          'fill-opacity': 0.2,
          'fill-outline-color': '#f03b20',
        },
      })
      .getLayer(`car-${this.id}-com-range`) as mapboxgl.FillLayer;
  }

  private get communicationRangeFeature() {
    return turf.buffer(
      {
        type: 'FeatureCollection',
        features: [
          {
            type: 'Feature',
            geometry: {
              type: 'Point',
              coordinates: [this.coordinates.lng, this.coordinates.lat],
            },
            properties: {},
          },
        ],
      },
      78,
      { units: 'meters' }
    );
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

  private updateNextFrame() {
    setTimeout(
      () => (this.animationFrame = requestAnimationFrame(this.update)),
      1000 / FPS
    );
  }

  /**
   * Update Car
   */
  private update = () => {
    if (this.removed) return;
    this.updateCoordinates();
    this.updateSource();
    this.updatePopup();
    if (!this.arrived) this.updateNextFrame();
  };

  private updateCoordinates = () => {
    const now = Date.now();
    let movementAmount =
      this.speed * (((now - this.prevTime) * 1.0) / MS_IN_HOUR);
    this.prevTime = now;
    while (movementAmount && !this.arrived && !this.obstacleDetected) {
      if (this.checkObstacles()) {
        this.speed = 0;
        this.emit('props-updated');
        break;
      }
      const dist = distanceInKm(this.coordinates, this.route[this.routeIndex]);
      if (movementAmount >= dist) {
        movementAmount -= dist;
        this.lat = this.route[this.routeIndex].lat;
        this.lng = this.route[this.routeIndex].lng;
        this.routeIndex++;
        this.emit('move');
      } else {
        const vector: Coordinates = {
          lng: (this.route[this.routeIndex].lng - this.coordinates.lng) / dist,
          lat: (this.route[this.routeIndex].lat - this.coordinates.lat) / dist,
        };
        this.lat += movementAmount * vector.lat;
        this.lng += movementAmount * vector.lng;
        movementAmount = 0;
        this.emit('move');
      }
    }
    if (this.arrived) {
      this.speed = 0;
      this.emit('props-updated');
      this.socket.emit('destination-reached', {
        id: this.id,
        coordinates: this.coordinates,
      });
    }
  };

  private checkObstacles(): boolean {
    const obstacles: turf.Geometry = (this.map.getSource('obstacles') as any)
      ._data;
    const lineStep = turf.lineString([
      [this.coordinates.lng, this.coordinates.lat],
      ...this.route
        .slice(this.routeIndex, this.routeIndex + 2)
        .map((c) => [c.lng, c.lat]),
    ]);
    const sensorRangeEndPoint = turf.along(lineStep, 100, {
      units: 'meters',
    });
    const sensorRange = turf.lineString([
      [this.coordinates.lng, this.coordinates.lat],
      sensorRangeEndPoint.geometry.coordinates,
    ]);

    if (!turf.booleanDisjoint(sensorRange, obstacles)) {
      this.obstacleDetected = true;
      this.emit('props-updated');
      this.socket.emit('obstacle-detected', {
        id: this.id,
        coordinates: this.coordinates,
        obstacle_coordinates: this.coordinates,
      });
      return true;
    }
    return false;
  }

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
    this.on('move', () => {
      this.socket.emit('update-location', {
        id: this.id,
        coordinates: this.coordinates,
      });
    });
  };

  private onClick = () => {
    if (this.popup) {
      this.popup.remove();
      this.popup = null;
    }
    this.popup = new mapboxgl.Popup()
      .setLngLat(this.coordinates as mapboxgl.LngLatLike)
      .setHTML(this.description)
      .addTo(this.map);

    const el = this.popup
      .getElement()
      .querySelector(`#link${this.id}`) as HTMLAnchorElement;
    if (el)
      el.onclick = () => {
        this.emit('focus');
      };

    this.popup.once('close', () => {
      this.popup = null;
      this.map?.setLayoutProperty(`car-${this.id}-route`, 'visibility', 'none');
      this.map?.setLayoutProperty(
        `car-${this.id}-com-range`,
        'visibility',
        'none'
      );
      this.emit('popup-closed', this);
    });

    this.on('props-updated', () => {
      if (!this.popup) return;
      this.popup.setHTML(this.description);
      const el = this.popup
        .getElement()
        .querySelector(`#link${this.id}`) as HTMLAnchorElement;
      if (el)
        el.onclick = () => {
          this.emit('focus');
        };
    });

    this.map?.setLayoutProperty(
      `car-${this.id}-route`,
      'visibility',
      'visible'
    );
    this.map?.setLayoutProperty(
      `car-${this.id}-com-range`,
      'visibility',
      'visible'
    );

    this.smoothlyFlyToCar(true);
    this.emit('click', this);
  };

  private updatePopup() {
    if (!this.popup) return;

    if (!this.map.isMoving()) {
      if (euclideanDistance(this.coordinates, this.map.getCenter()) > 1e-3) {
        this.prevCoordinates = this.map.getCenter();
        this.map.flyTo({
          center: this.coordinates as mapboxgl.LngLatLike,
        });
      }
    } else {
      this.smoothlyFlyToCar();
    }

    this.popup.setLngLat(this.coordinates as mapboxgl.LngLatLike);
    this.communicationRangeSource?.setData(this.communicationRangeFeature);
  }

  private smoothlyFlyToCar(now = false) {
    if (this.wasFlyingToCar) return;
    this.wasFlyingToCar = true;
    if (now) this.smoothFlyUtil();
    else this.map.once('moveend', () => this.smoothFlyUtil());
  }
  private smoothFlyUtil() {
    this.map.flyTo({
      center: this.coordinates as mapboxgl.LngLatLike,
      maxDuration: 200,
    });
    this.map.once('moveend', () => {
      this.wasFlyingToCar = false;
    });
  }

  private get description() {
    let description =
      '<h1 class="mapboxgl-popup-title">Car</h1>' + this.props.description;
    if (this.obstacleDetected) {
      description += '<p>Obstacle detected</p>';
    }
    if (!this.focused)
      description += '<a id="link{id}">Go to the car interface</a>';
    return interpolateString(description, this);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public on(
    type:
      | 'click'
      | 'focus'
      | 'move'
      | 'popup-closed'
      | 'props-updated'
      | 'obstacle-detected',
    handler: any
  ) {
    this.subscribe(type, handler);
    return this;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private subscribe(
    type:
      | 'click'
      | 'focus'
      | 'move'
      | 'popup-closed'
      | 'props-updated'
      | 'obstacle-detected',
    handler: (...args: any) => void
  ) {
    if (!this.handlers[type]) this.handlers[type] = [];
    this.handlers[type].push(handler);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private emit(
    type:
      | 'click'
      | 'focus'
      | 'move'
      | 'popup-closed'
      | 'props-updated'
      | 'obstacle-detected',
    ...args: any[]
  ) {
    if (!this.handlers[type]) return;

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    this.handlers[type].forEach((handler: (...args: any[]) => void) =>
      handler.call(this, ...args)
    );
  }

  public export() {
    return {
      id: this.id,
      route: this.route,
      speed: this.speed,
      lng: this.route[0].lng,
      lat: this.route[0].lat,
      obstacleDetected: this.obstacleDetected,
      type: 'car',
    };
  }

  public remove() {
    this.popup?.remove();
    this.map.removeLayer(this.sourceId);
    this.map.removeSource(this.sourceId);
    this.map.removeLayer(`car-${this.id}-route`);
    this.map.removeSource(`car-${this.id}-route`);
    this.map.removeLayer(`car-${this.id}-com-range`);
    this.map.removeSource(`car-${this.id}-com-range`);
    cancelAnimationFrame(this.animationFrame);
    this.removed = true;
    this.map.off('click', this.sourceId, this.onClick);
  }

  public hide() {
    this.popup?.remove();
    this.map.setLayoutProperty(this.sourceId, 'visibility', 'none');
    this.map.setLayoutProperty(`car-${this.id}-route`, 'visibility', 'none');
    this.map.setLayoutProperty(
      `car-${this.id}-com-range`,
      'visibility',
      'none'
    );
    this.focused = false;
    this.map.off('click', this.sourceId, this.onClick);
  }
  public show(focus = false) {
    this.map.setLayoutProperty(this.sourceId, 'visibility', 'visible');
    this.map
      .off('click', this.sourceId, this.onClick)
      .on('click', this.sourceId, this.onClick);
    this.focused = false;
    if (focus) {
      this.focused = true;
      this.onClick();
    }
    this.emit('props-updated');
  }
}

export default Car;
