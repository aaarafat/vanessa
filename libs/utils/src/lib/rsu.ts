import mapboxgl from 'mapbox-gl';
import { distanceInKm } from './distance';
import { interpolateString } from './string-utils';
import { Coordinates, IRSU, RSUProps } from './types';
import * as turf from '@turf/turf';

const carDefaultProps: RSUProps = {
  title: 'RSU',
  description: `<ul class="popup">
    <li>id: {id}</li>
    <li>radius: {radius} km</li>
  </ul>`,
};

/**
 * Car Class
 */
export class RSU implements IRSU {
  public id: number;
  public lat: number;
  public lng: number;
  public radius: number;
  public sourceId: string;
  private source: mapboxgl.GeoJSONSource | undefined;
  private layer: mapboxgl.LineLayer | undefined;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private map: mapboxgl.Map;
  private popup: mapboxgl.Popup | null;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private handlers: Record<string, ((...arg: any) => void)[]>;
  private wasFlyingToRSU: boolean;

  constructor(rsu: Partial<IRSU> & { map: mapboxgl.Map }) {
    this.id = rsu.id || Date.now();
    this.sourceId = `car-${this.id}`;
    this.lat = rsu.lat || 0;
    this.lng = rsu.lng || 0;
    this.radius = rsu.radius || 10;
    this.map = rsu.map;
    this.popup = null;

    this.handlers = {};
    this.wasFlyingToRSU = false;

    this.draw();
    this.attachHandlers();
  }

  public get coordinates(): Coordinates {
    return { lat: this.lat, lng: this.lng };
  }

  public draw(): void {
    const center = turf.point([this.lng, this.lat]);
    const options = {
      steps: 80,
      units: 'kilometers' as turf.Units,
      properties: { ...this.props },
    };

    const data = turf.circle(center, this.radius, options);
    const geojson: mapboxgl.GeoJSONSourceRaw = {
      type: 'geojson',
      data,
    };

    this.source = this.map
      .addSource(this.sourceId, geojson)
      .getSource(this.sourceId) as mapboxgl.GeoJSONSource;

    this.layer = this.map
      .addLayer({
        id: this.sourceId,
        source: this.sourceId,
        type: 'line',
        paint: {
          'line-color': '#ff0000',
          'line-opacity': 0.5,
          'line-width': 2,
        },
      })
      .getLayer(this.sourceId) as mapboxgl.LineLayer;
  }

  private get props(): RSUProps {
    return {
      ...carDefaultProps,
      id: this.id,
      lat: this.lat,
      lng: this.lng,
      radius: this.radius,
    };
  }

  private attachHandlers = () => {
    this.map.on('click', this.sourceId, this.onClick);
  };

  private onClick = () => {
    this.popup = new mapboxgl.Popup()
      .setLngLat(this.coordinates as mapboxgl.LngLatLike)
      .setHTML(this.description)
      .addTo(this.map);

    this.popup.on('close', () => {
      this.popup = null;
      this.map?.setLayoutProperty(`car-${this.id}-route`, 'visibility', 'none');
      this.emit('popup-closed', this);
    });

    this.on('props-updated', () => {
      this.popup?.setHTML(this.description);
    });

    this.map?.setLayoutProperty(
      `car-${this.id}-route`,
      'visibility',
      'visible'
    );

    this.smoothlyFlyToRSU(true);
    this.emit('click', this);
  };

  private smoothlyFlyToRSU(now = false) {
    if (this.wasFlyingToRSU) return;
    this.wasFlyingToRSU = true;
    if (now) this.smoothFlyUtil();
    else this.map.once('moveend', () => this.smoothFlyUtil());
  }
  private smoothFlyUtil() {
    this.map.flyTo({
      center: this.coordinates as mapboxgl.LngLatLike,
      maxDuration: 200,
    });
    this.map.once('moveend', () => {
      this.wasFlyingToRSU = false;
    });
  }

  private get description() {
    const description =
      '<h1 class="mapboxgl-popup-title">RSU</h1>' + this.props.description;
    return interpolateString(description, this);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public on(
    type: 'click' | 'move' | 'popup-closed' | 'props-updated',
    handler: any
  ) {
    this.subscribe(type, handler);
    return this;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private subscribe(
    type: 'click' | 'move' | 'popup-closed' | 'props-updated',
    handler: (...args: any) => void
  ) {
    if (!this.handlers[type]) this.handlers[type] = [];
    this.handlers[type].push(handler);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private emit(
    type: 'click' | 'move' | 'popup-closed' | 'props-updated',
    ...args: any[]
  ) {
    if (!this.handlers[type]) return;

    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    this.handlers[type].forEach((handler: (...args: any[]) => void) =>
      handler.call(this, ...args)
    );
  }
}

export default RSU;
