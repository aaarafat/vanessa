export interface MapProps {
  /**Current Zoom value */
  currentZoom?: number;

  /**Current Lng value */
  currentLng?: number;

  /**Current Lat value */
  currentLat?: number;

  cars?: Car[];
}

export interface Car {
  id: number;
  lat: number;
  lng: number;
}