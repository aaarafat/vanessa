import { euclideanDistance } from './distance';
import { Coordinates, ICar } from './types';

const UPDATE_INTERVAL = 1000; // every 1 second
const SPEED = 0.000003;
const COS_45 = Math.cos(0.25 * Math.PI);

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

  constructor(car: ICar) {
    this.id = car.id;
    this.lat = car.lat;
    this.lng = car.lng;
    this.route = car.route;
    this.routeIndex = 0;

    this.updateIntervalId = setInterval(this.update, UPDATE_INTERVAL);
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
    let movementAmount = SPEED * UPDATE_INTERVAL;
    while (movementAmount && !this.arrived) {
      const dist = euclideanDistance(
        this.coordinates,
        this.route[this.routeIndex]
      );

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
        this.lat += movementAmount * COS_45;
        this.lng += movementAmount * COS_45;
        movementAmount = 0;
      }
    }
  };
}

export default Car;
