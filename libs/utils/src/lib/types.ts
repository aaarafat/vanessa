export interface Coordinates {
  lng: number;
  lat: number;
}

export interface ICar {
  id: number;
  lat: number;
  lng: number;
  speed: number;
  route: Coordinates[];
}

export interface CarProps extends Partial<ICar> {
  title?: string;
  name?: string;
  description?: string;
}