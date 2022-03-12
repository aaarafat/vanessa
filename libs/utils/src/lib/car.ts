import { distanceInKm, euclideanDistance } from './distance';
import { Coordinates, ICar } from './types';

const MS_IN_HOUR = 1000 * 60 * 60;
const SPEED_KM_H = 100; // KM/H

/**
 * Car Class
 */
export class Car implements ICar {
  public id: number;
  public lat: number;
  public lng: number;
  public route: Coordinates[];
  private routeIndex: number;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  private updateIntervalId: any;
  private prevTime: number;

  constructor(car: ICar) {
    this.id = car.id;
    this.lat = car.lat;
    this.lng = car.lng;
    this.route = car.route;
    this.routeIndex = 0;

    this.prevTime = Date.now();
    this.update();
  }

  public get coordinates(): Coordinates {
    return { lat: this.lat, lng: this.lng };
  }

  public get arrived(): boolean {
    return this.routeIndex === this.route.length;
  }

  private display = () => {
    // console.log(
    //   `ID: ${this.id}, Lat: ${this.lat}, Lng: ${this.lng}, RouteIndex: ${
    //     this.routeIndex
    //   }, Arrived: ${this.arrived ? 'Yes' : 'No'}`
    // );
  };

  /**
   * Update Car
   */
  private update = () => {
    this.updateRouteIndex();
    this.updateCoordinates();
    this.display();
    if(!this.arrived)
      requestAnimationFrame(this.update)
  };

  private updateRouteIndex = () => {
    while (this.routeIndex < this.route.length - 1) {
      const coords = this.coordinates;
      const dist = euclideanDistance(coords, this.route[this.routeIndex + 1]);
      const distBetweenIndexes = euclideanDistance(
        this.route[this.routeIndex],
        this.route[this.routeIndex + 1]
      );

      if (distBetweenIndexes > dist) {
        // the car in the middle -> update the routeIndex
        this.routeIndex++;
      } else {
        // nothing to be updated
        break;
      }
    }
  };

  private updateCoordinates = () => {
    const now = Date.now();
    let movementAmount = SPEED_KM_H * (((now - this.prevTime) * 1.0) / MS_IN_HOUR);
    this.prevTime = now;
    while (movementAmount && !this.arrived) {
      const dist = distanceInKm(this.coordinates, this.route[this.routeIndex]);

      if (movementAmount >= dist) {
        movementAmount -= dist;
        this.lat = this.route[this.routeIndex].lat;
        this.lng = this.route[this.routeIndex].lng;
        this.routeIndex++;
        if (this.routeIndex === this.route.length) {
          // you arrived
          this.display();
          // console.log('ARRRRIIIVEEEEEEEEEEEEED!!!!!!!!!!!!!!!!!!');
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
}

export default Car;
