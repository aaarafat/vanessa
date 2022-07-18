import mapboxgl from 'mapbox-gl';
import { interpolateString } from './utils';
import { Coordinates, IRSU, PartialExcept, RSUProps } from './types';
import * as turf from '@turf/turf';

const rsuDefaultProps: RSUProps = {
  title: 'RSU',
  description: `<ul class="popup">
    <li>ID: {id}</li>
    <li>Range: {range} m</li>
  </ul>`,
};

/**
 * RSU Class
 */
export class RSU {
  public id: number;
  public lat: number;
  public lng: number;
  public range: number;
  public sourceId: string;
  public clickableSourceId: string;
  private map: mapboxgl.Map;
  private popup: mapboxgl.Popup | null;
  private removed = false;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private handlers: Record<string, ((...arg: any) => void)[]>;
  private wasFlyingToRSU: boolean;

  public port: number;

  constructor(rsu: PartialExcept<IRSU, 'map'>) {
    this.id = rsu.id || Date.now();
    this.sourceId = `rsu-${this.id}`;
    this.clickableSourceId = `rsu-clickable-${this.id}`;
    this.lat = rsu.lat || 0;
    this.lng = rsu.lng || 0;
    this.range = rsu.range || 500;
    this.map = rsu.map;

    this.popup = null;

    this.handlers = {};
    this.wasFlyingToRSU = false;

    this.port = rsu.port || -1;

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
      units: 'meters' as turf.Units,
      properties: { ...this.props },
    };

    const data = turf.circle(center, this.range, options);
    const geojson: mapboxgl.GeoJSONSourceRaw = {
      type: 'geojson',
      data,
    };

    this.map.addSource(this.sourceId, geojson);

    this.map.addLayer({
      id: this.sourceId,
      source: this.sourceId,
      type: 'line',
      paint: {
        'line-color': '#ff0000',
        'line-opacity': 0.5,
        'line-width': 2,
      },
    });

    const g: mapboxgl.GeoJSONSourceRaw = {
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

    this.map.addSource(this.clickableSourceId, g);

    this.map.addLayer({
      id: this.clickableSourceId,
      source: this.clickableSourceId,
      type: 'circle',
      paint: {
        'circle-radius': 10,
        'circle-color': '#ff0000',
      },
    });
  }

  public setRange(range: number) {
    this.range = range;
    if (!this.removed) {
      this.map.removeLayer(this.sourceId);
      this.map.removeSource(this.sourceId);

      this.map.removeLayer(this.clickableSourceId);
      this.map.removeSource(this.clickableSourceId);
    }
    this.draw();
    this.updatePopup();
  }

  private get props(): RSUProps {
    return {
      ...rsuDefaultProps,
      id: this.id,
      lat: this.lat,
      lng: this.lng,
      range: this.range,
    };
  }

  private attachHandlers = () => {
    this.map.on('click', this.clickableSourceId, this.onClick);
  };

  private onClick = (e: mapboxgl.MapMouseEvent) => {
    if (e?.originalEvent.defaultPrevented) return;
    e?.originalEvent.preventDefault();
    if (this.popup) {
      this.popup.remove();
      this.popup = null;
    }

    this.popup = new mapboxgl.Popup()
      .setLngLat(this.coordinates as mapboxgl.LngLatLike)
      .setHTML(this.description)
      .addTo(this.map);

    this.popup.on('close', () => {
      this.popup = null;
      this.emit('popup:close', this);
    });

    this.smoothlyFlyToRSU(true);
    this.emit('click', this);
  };

  private updatePopup = () => {
    if (!this.popup) return;
    this.popup.setHTML(this.description);
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
    let description =
      '<h1 class="mapboxgl-popup-title">RSU</h1>' + this.props.description;

    if (this.port > 0) {
      description += `<p>Port: ${this.port}</p>`;
    }

    description +=
      '<a id="link{id}" href="http://localhost:4202/{port}" target="_blank">Go to the rsu interface</a>';

    return interpolateString(description, this);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  public on(
    type: 'click' | 'move' | 'popup:close' | 'props:change',
    handler: any
  ) {
    this.subscribe(type, handler);
    return this;
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private subscribe(
    type: 'click' | 'move' | 'popup:close' | 'props:change',
    handler: (...args: any) => void
  ) {
    if (!this.handlers[type]) this.handlers[type] = [];
    this.handlers[type].push(handler);
  }

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private emit(
    type: 'click' | 'move' | 'popup:close' | 'props:change',
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
      lng: this.lng,
      lat: this.lat,
      range: this.range,
      type: 'rsu',
    };
  }

  public remove() {
    if (this.removed) return;
    this.removed = true;
    this.popup?.remove();
    this.popup = null;

    this.map.removeLayer(this.sourceId);
    this.map.removeSource(this.sourceId);

    this.map.removeLayer(this.clickableSourceId);
    this.map.removeSource(this.clickableSourceId);

    this.map.off('click', this.clickableSourceId, this.onClick);
    this.handlers = {};
  }
}

export default RSU;
