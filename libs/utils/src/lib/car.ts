import mapboxgl from 'mapbox-gl';
import { distanceInKm, euclideanDistance } from './distance';
import {
  createFeaturePoint,
  getObstacleFeatures,
  interpolateString,
} from './utils';
import { Coordinates, ICar, CarProps, PartialExcept } from './types';
import {
  MS_IN_HOUR,
  FPS,
  directionsAPI,
  directionsAPIParams,
} from './constants';
import * as turf from '@turf/turf';

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
  public carSpeed: number;
  public route: Coordinates[];
  public originalDirections: GeoJSON.Feature;

  private prevCoordinates: Coordinates;
  private routeIndex: number;

  public sourceId: string;
  public routeSourceId: string;
  public communicationRangeSourceId: string;

  private source: mapboxgl.GeoJSONSource | undefined;
  private routeSource: mapboxgl.GeoJSONSource | undefined;
  private communicationRangeSource: mapboxgl.GeoJSONSource | undefined;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private prevTime: number;
  private map: mapboxgl.Map;
  // private socket: Socket;
  private popup: mapboxgl.Popup | null;
  private obstacleDetected: boolean;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private handlers: Record<string, ((...arg: any) => void)[]>;
  private wasFlyingToCar: boolean;
  private animationFrame: number;
  private removed = false;
  private focused: boolean;
  private initiated = false;
  public stopped = false;

  public port: number;

  constructor(car: PartialExcept<ICar, 'map'>, { displayOnly = false } = {}) {
    this.id = car.id || Date.now();
    this.sourceId = `car-${this.id}`;
    this.routeSourceId = `car-${this.id}-route`;
    this.communicationRangeSourceId = `car-${this.id}-communication-range`;
    this.map = car.map;

    this.lat = car.lat || car.route?.[0].lat || 0;
    this.lng = car.lng || car.route?.[0].lng || 0;
    this.prevCoordinates = this.coordinates;
    this.carSpeed = car.speed ?? 10;
    this.obstacleDetected = false;
    this.route = car.route || [];
    this.routeIndex = car.destinationReached ? this.route.length : 0;
    this.originalDirections = turf.lineString(
      this.route.map((r: Coordinates) => [r.lng, r.lat])
    );

    this.popup = null;
    this.focused = false;

    this.handlers = {};
    this.wasFlyingToCar = false;
    this.animationFrame = 0;

    this.port = car.port || -1;

    this.draw();
    this.attachHandlers();

    this.prevTime = Date.now();

    if (displayOnly) {
      this.focused = true;
      this.initiated = true;
      this.obstacleDetected = car.obstacleDetected || false;
      this.map.panTo([this.lng, this.lat]);
    } else {
      this.map.setPaintProperty(this.sourceId, 'circle-color', '#ababab');
    }
  }

  public get coordinates(): Coordinates {
    return { lat: this.lat, lng: this.lng };
  }

  public get arrived(): boolean {
    return this.routeIndex === this.route.length;
  }

  public get speed(): number {
    if (this.arrived || this.obstacleDetected || this.stopped) return 0;
    return this.carSpeed;
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

    this.map
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

    this.routeSource = this.map
      .addSource(this.routeSourceId, {
        type: 'geojson',
        data: {
          type: 'FeatureCollection',
          features: [this.originalDirections],
        },
      })
      .getSource(this.routeSourceId) as mapboxgl.GeoJSONSource;

    this.map
      .addLayer(
        {
          id: this.routeSourceId,
          source: this.routeSourceId,
          type: 'line',
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
      .getLayer(this.routeSourceId) as mapboxgl.LineLayer;

    this.communicationRangeSource = this.map
      .addSource(this.communicationRangeSourceId, {
        type: 'geojson',
        data: this.getCommunicationRangeFeature(),
      })
      .getSource(this.communicationRangeSourceId) as mapboxgl.GeoJSONSource;

    this.map
      .addLayer({
        id: this.communicationRangeSourceId,
        source: this.communicationRangeSourceId,
        type: 'fill',
        layout: {
          visibility: 'none',
        },
        paint: {
          'fill-color': '#f03b20',
          'fill-opacity': 0.2,
          'fill-outline-color': '#f03b20',
        },
      })
      .getLayer(this.communicationRangeSourceId) as mapboxgl.FillLayer;
  }

  private getCommunicationRangeFeature() {
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

  public startMovement() {
    if (this.initiated) return;
    console.log('Car started movement', this.id);
    this.initiated = true;
    this.map.setPaintProperty(this.sourceId, 'circle-color', '#007cbf');
    this.prevTime = Date.now();
    this.update();
    this.updatePopupProps();
  }

  public updateLocationFromData = (coordinates: Coordinates) => {
    this.updateFromData('update-location', coordinates);
  };

  public updateDestinationFromData = (coordinates: Coordinates) => {
    this.updateFromData('destination-reached', coordinates);
  };

  public updateObstacleDetectedFromData = () => {
    this.updateFromData('obstacle-detected');
  };

  public updateRouteFromData = (route: Coordinates[]) => {
    this.route = route;
    this.originalDirections = turf.lineString(
      this.route.map((r: Coordinates) => [r.lng, r.lat])
    );
    this.routeSource?.setData(
      turf.featureCollection([this.originalDirections])
    );
  };

  private updateFromData(
    type: 'destination-reached' | 'obstacle-detected' | 'update-location',
    data: Partial<ICar> = {}
  ) {
    if (this.removed) return;
    switch (type) {
      case 'destination-reached':
        this.lat = data.lat ?? this.lat;
        this.lng = data.lng ?? this.lng;
        this.routeIndex = this.route.length;
        this.updatePopupProps();
        break;
      case 'obstacle-detected':
        this.obstacleDetected = true;
        this.updatePopupProps();
        break;
      case 'update-location':
        this.lat = data.lat ?? this.lat;
        this.lng = data.lng ?? this.lng;
        break;
      default:
        break;
    }
    this.updateSource();
    this.updateDetails();
  }

  private update = () => {
    if (this.removed || !this.initiated) return;
    this.updateCoordinates();
    this.updateSource();
    this.updateDetails();
    if (!this.willUpdate()) return;
    this.updateNextFrame();
  };

  private updateCoordinates = () => {
    const now = Date.now();
    let movementAmount =
      this.speed * (((now - this.prevTime) * 1.0) / MS_IN_HOUR);
    this.prevTime = now;
    while (this.canMove(movementAmount)) {
      const dist = distanceInKm(this.coordinates, this.route[this.routeIndex]);
      if (movementAmount >= dist) {
        movementAmount -= dist;
        this.lat = this.route[this.routeIndex].lat;
        this.lng = this.route[this.routeIndex].lng;
        this.routeIndex++;
      } else {
        const vector: Coordinates = {
          lng: (this.route[this.routeIndex].lng - this.coordinates.lng) / dist,
          lat: (this.route[this.routeIndex].lat - this.coordinates.lat) / dist,
        };
        this.lat += movementAmount * vector.lat;
        this.lng += movementAmount * vector.lng;
        movementAmount = 0;
      }
      this.emit('move');
    }
  };

  private willUpdate() {
    if (this.arrived) {
      this.updatePopupProps();
      this.emit('destination-reached');
      return false;
    } else if (this.obstacleDetected) {
      return false;
    }
    return true;
  }

  private canMove(movementAmount: number) {
    return movementAmount && !this.arrived && !this.isObstacleDetected();
  }

  private isObstacleDetected(): boolean {
    const obstacles: turf.Polygon = (this.map.getSource('obstacles') as any)
      ._data;
    const obstaclesPoints: turf.FeatureCollection<turf.Point> = (
      this.map.getSource('obstacles-points') as any
    )._data;

    const routeSlice = this.route
      .slice(this.routeIndex, this.routeIndex + 2)
      .map((c) => [c.lng, c.lat]);
    const lineStep = turf.lineString([
      [this.coordinates.lng, this.coordinates.lat],
      ...routeSlice,
    ]);

    const sensorRangeEndPoint = turf.along(lineStep, 100, {
      units: 'meters',
    });
    const sensorRange = turf.lineString([
      [this.coordinates.lng, this.coordinates.lat],
      sensorRangeEndPoint.geometry.coordinates,
    ]);

    const intersections = turf.lineIntersect(sensorRange, obstacles).features;
    if (intersections.length) {
      this.obstacleDetected = true;
      const point = turf.nearestPoint(intersections[0], obstaclesPoints)
        .geometry.coordinates;
      this.emit('obstacle-detected', { lng: point[0], lat: point[1] });
      this.updatePopupProps();
      return true;
    }
    return false;
  }

  public updateRoute = async (obstacles: Coordinates[]) => {
    if (this.obstacleDetected) return false;
    const o = obstacles.map((c) => createFeaturePoint(c));
    if (!this.checkObstaclesOnRoute(o)) return false;

    const result = await this.getRoute(o);
    if (!result) return false;

    this.originalDirections = turf.lineString(result);
    this.route = result.map((c) => ({ lat: c[1], lng: c[0] }));
    this.routeIndex = 0;
    this.routeSource?.setData(
      turf.featureCollection([this.originalDirections])
    );
    return true;
  };

  private getRoute = async (
    obstacles: turf.Feature<turf.Point>[]
  ): Promise<turf.Position[] | null> => {
    const dest = this.route[this.route.length - 1];
    const o = obstacles
      .map(
        (c) =>
          `point(${c.geometry.coordinates[0]} ${c.geometry.coordinates[1]})`
      )
      .join(',');
    const params = {
      ...directionsAPIParams,
      exclude: o,
    };

    const response = await fetch(
      `${directionsAPI}${this.lng},${this.lat};${dest.lng},${
        dest.lat
      }.json?${new URLSearchParams(params).toString()}`
    );
    const data = await response.json();
    if (data.code !== 'Ok') return null;
    return data.routes[0].geometry.coordinates;
  };

  public checkObstaclesOnRoute = (obstacles: turf.Feature<turf.Point>[]) => {
    const obstaclesFeatures = getObstacleFeatures(obstacles);

    const routeSlice = this.route
      .slice(this.routeIndex)
      .map((c) => [c.lng, c.lat]);
    if (routeSlice.length < 1) return false;

    const remainingRoute = turf.lineString([
      [this.lng, this.lat],
      ...routeSlice,
    ]);

    if (!turf.booleanDisjoint(remainingRoute, obstaclesFeatures as any)) {
      return true;
    }
    return false;
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

  public setSpeed = (speed: number) => {
    this.carSpeed = speed;
    this.updatePopupProps();
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
    this.bindElements();

    this.communicationRangeSource?.setData(this.getCommunicationRangeFeature());

    this.popup.once('close', () => {
      this.popup = null;
      this.setDetailsLayersVisibility('none');
      this.emit('popup:close', this);
    });

    this.setDetailsLayersVisibility('visible');
    this.smoothlyFlyToCar(true);
    this.emit('click', this);
  };

  private setDetailsLayersVisibility(visibility: 'visible' | 'none' = 'none') {
    this.map?.setLayoutProperty(this.routeSourceId, 'visibility', visibility);
    this.map?.setLayoutProperty(
      this.communicationRangeSourceId,
      'visibility',
      visibility
    );
  }

  private updatePopupProps = () => {
    if (!this.popup) return;
    this.popup.setHTML(this.description);
    this.bindElements();
  };

  private bindElements() {
    this.bindAnchorElement();
    if (!this.stopped) this.bindSpeedControlElementStop();
    else this.bindSpeedControlElementMove();
  }

  private bindSpeedControlElementStop() {
    if (!this.popup) return;
    const el = this.popup
      .getElement()
      .querySelector(`#control${this.id}-stop`) as HTMLAnchorElement;
    if (el)
      el.onclick = () => {
        this.stopped = true;
        this.emit('change-speed', this);
        this.updatePopupProps();
      };
  }

  private bindSpeedControlElementMove() {
    if (!this.popup) return;
    const el = this.popup
      .getElement()
      .querySelector(`#control${this.id}-move`) as HTMLAnchorElement;
    if (el)
      el.onclick = () => {
        this.stopped = false;
        this.emit('change-speed', this);
        this.updatePopupProps();
      };
  }

  private bindAnchorElement() {
    if (!this.popup) return;
    const el = this.popup
      .getElement()
      .querySelector(`#link${this.id}`) as HTMLAnchorElement;
    if (el)
      el.onclick = () => {
        this.emit('focus');
      };
  }

  private updateDetails() {
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
    this.communicationRangeSource?.setData(this.getCommunicationRangeFeature());
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

    if (this.port > 0) {
      description += `<p>Port: ${this.port}</p>`;
    }
    if (this.obstacleDetected) {
      description += '<p>Obstacle detected</p>';
    }

    if (!this.focused)
      description += '<a id="link{id}">Go to the car interface</a>';

    if (!this.arrived && !this.obstacleDetected && !this.focused) {
      description += !this.stopped
        ? '<a id="control{id}-stop">Stop</a>'
        : '<a id="control{id}-move">Move</a>';
    }

    if (!this.initiated) {
      description += 'Initializing...';
    }

    return interpolateString(description, this);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public on(
    type:
      | 'click'
      | 'focus'
      | 'move'
      | 'popup:close'
      | 'change-speed'
      | 'obstacle-detected'
      | 'destination-reached',
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
      | 'popup:close'
      | 'change-speed'
      | 'obstacle-detected'
      | 'destination-reached',
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
      | 'popup:close'
      | 'change-speed'
      | 'obstacle-detected'
      | 'destination-reached',
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
      speed: this.carSpeed,
      lng: this.route[0].lng,
      lat: this.route[0].lat,
      type: 'car',
    };
  }

  public remove() {
    if (this.removed) return;
    this.removed = true;
    this.popup?.remove();

    this.map.removeLayer(this.sourceId);
    this.map.removeSource(this.sourceId);

    this.map.removeLayer(this.routeSourceId);
    this.map.removeSource(this.routeSourceId);

    this.map.removeLayer(this.communicationRangeSourceId);
    this.map.removeSource(this.communicationRangeSourceId);

    cancelAnimationFrame(this.animationFrame);
    this.map.off('click', this.sourceId, this.onClick);
    this.handlers = {};
  }

  public hide() {
    this.focused = false;
    this.popup?.remove();

    this.map.setLayoutProperty(this.sourceId, 'visibility', 'none');
    this.map.setLayoutProperty(this.routeSourceId, 'visibility', 'none');
    this.map.setLayoutProperty(
      this.communicationRangeSourceId,
      'visibility',
      'none'
    );

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
    this.updatePopupProps();
  }
}

export default Car;
