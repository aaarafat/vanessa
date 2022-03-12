export interface Coordinates {
  lng: number;
  lat: number;
}

export interface ICar {
  id: number;
  lat: number;
  lng: number;
  route: Coordinates[];
}
