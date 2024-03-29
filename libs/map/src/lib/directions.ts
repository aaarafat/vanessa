import mapboxgl from 'mapbox-gl';

// eslint-disable-next-line @typescript-eslint/no-var-requires
const MapboxDirections = require('vanessa-mapbox-gl-directions/dist/mapbox-gl-directions');
// eslint-disable-next-line @typescript-eslint/no-var-requires

export class Directions extends MapboxDirections {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/no-useless-constructor
  constructor(options?: any) {
    super(options);
    this.freezed = 0;
  }

  freeze() {
    if (this.freezed++) return;
    this._map.off('touchstart', this.move);
    this._map.off('touchstart', this.onDragDown);

    this._map.off('mousedown', this.onDragDown);
    this._map.off('mousemove', this.move);

    this.storeUnsubscribe();
    delete this.storeUnsubscribe;
  }

  unfreeze() {
    if (--this.freezed) return;
    this._map.on('touchstart', this.move);
    this._map.on('touchstart', this.onDragDown);

    this._map.on('mousedown', this.onDragDown);
    this._map.on('mousemove', this.move);

    this.subscribedActions();
  }

  mapState() {
    super.mapState();
    if (this.options.interactive !== false) {
      this.overrideClickEvent();
    }
  }

  private overrideClickEvent() {
    this._map.off('click', this.onClick);
    this._map.on('click', (e: mapboxgl.MapMouseEvent) => {
      if (this.options.interactive === false || this._isPointOnCar(e.point)) {
        return;
      }

      this.onClick(e);
    });
  }

  _isPointOnCar(point: mapboxgl.Point) {
    const features = this._map.queryRenderedFeatures(point, {
      filter: [
        'any',
        ['in', 'Car', ['get', 'title']],
        ['in', 'RSU', ['get', 'title']],
      ],
      validate: false,
    });

    return features?.length;
  }

  reset() {
    this.resetUtil();
    this.actions.eventEmit('reset', {});
    setTimeout(() => this.resetUtil(), 500);
  }
  private resetUtil() {
    this.removeRoutes();
    this._map.getSource('directions').setData({
      type: 'FeatureCollection',
      features: [],
    });
  }
}
