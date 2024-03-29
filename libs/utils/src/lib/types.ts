export type PartialExcept<T, K extends keyof T> = Partial<T> & Pick<T, K>;

export interface Coordinates {
  lng: number;
  lat: number;
}

export interface ICommon {
  id: number;
  lat: number;
  lng: number;
  map: mapboxgl.Map;
  port: number;
}

export interface ICar extends ICommon {
  speed: number;
  route: Coordinates[];
  obstacleDetected?: boolean;
  destinationReached?: boolean;
  stopped?: boolean;
}
export interface IRSU extends ICommon {
  range: number;
}

export interface CarProps extends Partial<ICar> {
  title?: string;
  name?: string;
  description?: string;
}

export interface RSUProps extends Partial<IRSU> {
  title?: string;
  name?: string;
  description?: string;
}
