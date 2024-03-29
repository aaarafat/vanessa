import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Car, IRSU, RSU } from '@vanessa/utils';
import mapboxgl from 'mapbox-gl';

export interface SimulationState {
  rsus: RSU[];
  cars: Car[];
  rsusData: Partial<IRSU>[];
  focusedCar: Car | null;
  focusedRSU: RSU | null;
  carsReceivedMessages: Record<string, any>;
}

const initialState: SimulationState = {
  rsusData: [
    {
      id: 1,
      lng: 31.213,
      lat: 30.0252,
      range: 250,
    },
    {
      id: 2,
      lng: 31.2029,
      lat: 30.0269,
      range: 500,
    },
    {
      id: 3,
      lng: 31.2129,
      lat: 30.0185,
      range: 500,
    },
  ],
  rsus: [],
  cars: [],
  carsReceivedMessages: {},
  focusedCar: null,
  focusedRSU: null,
};

export const simulationSlice = createSlice({
  name: 'simulation',
  initialState,
  reducers: {
    initRSUs: (
      state,
      action: PayloadAction<{
        map: mapboxgl.Map;
      }>
    ) => {
      const { map } = action.payload;
      state.rsusData.forEach((rsu) =>
        state.rsus.push(new RSU({ ...rsu, map }))
      );
    },
    addCar: (state, action: PayloadAction<Car>) => {
      const car = action.payload;
      state.cars.push(car);
      state.carsReceivedMessages[car.id] = [];
    },
    addRSU: (state, action: PayloadAction<RSU>) => {
      state.rsus.push(action.payload);
    },
    clearState: (
      state,
      action: PayloadAction<{
        removeRSUs: boolean;
      }>
    ) => {
      const { rsus, cars } = state;
      const { removeRSUs } = action.payload;
      // cars.forEach((car) => car.remove());
      state.focusedCar = null;
      cars.splice(0, cars.length);
      if (removeRSUs) {
        // rsus.forEach((rsu) => rsu.remove());
        rsus.splice(0, rsus.length);
      }
    },
    focusCar: (state, action: PayloadAction<number>) => {
      state.focusedCar =
        state.cars.find((car) => car.id === action.payload) || null;
    },
    unfocusCar: (state) => {
      state.focusedCar = null;
    },
    focusRSU: (state, action: PayloadAction<number>) => {
      state.focusedRSU =
        state.rsus.find((rsu) => rsu.id === action.payload) || null;
    },
    unfocusRSU: (state) => {
      state.focusedRSU = null;
    },
    addMessage: (
      state,
      action: PayloadAction<{ id: number; message: any }>
    ) => {
      if (!state.carsReceivedMessages[action.payload.id]) return;
      state.carsReceivedMessages[action.payload.id].unshift(
        action.payload.message
      );
    },
  },
});

export const {
  initRSUs,
  addCar,
  addRSU,
  clearState,
  focusCar,
  unfocusCar,
  focusRSU,
  unfocusRSU,
  addMessage,
} = simulationSlice.actions;

export default simulationSlice.reducer;
