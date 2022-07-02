import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Car, IRSU, RSU } from '@vanessa/utils';
import mapboxgl from 'mapbox-gl';

export interface SimulationState {
  rsus: RSU[];
  cars: Car[];
  rsusData: Partial<IRSU>[];
  focusedCar: number | null;
  carsReceivedMessages: Record<string, any>;
}

const initialState: SimulationState = {
  rsusData: [
    {
      id: 1,
      lng: 31.213,
      lat: 30.0252,
      radius: 0.25,
    },
    {
      id: 2,
      lng: 31.2029,
      lat: 30.0269,
      radius: 0.5,
    },
    {
      id: 3,
      lng: 31.2129,
      lat: 30.0185,
      radius: 0.5,
    },
  ],
  rsus: [],
  cars: [],
  carsReceivedMessages: {},
  focusedCar: null,
};

export const simulationSlice = createSlice({
  name: 'simulation',
  initialState,
  reducers: {
    initRSUs: (
      state,
      action: PayloadAction<{
        map: mapboxgl.Map;
        socket: any;
      }>
    ) => {
      const { map, socket } = action.payload;
      state.rsusData.forEach((rsu) =>
        state.rsus.push(new RSU({ ...rsu, map, socket }))
      );
    },
    addCar: (state, action: PayloadAction<Car>) => {
      state.cars.push(action.payload);
      state.carsReceivedMessages[action.payload.id] = [];
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
      cars.forEach((car) => car.remove());
      cars.splice(0, cars.length);
      if (removeRSUs) {
        rsus.forEach((rsu) => rsu.remove());
        rsus.splice(0, rsus.length);
      }
    },
    focusCar: (state, action: PayloadAction<number>) => {
      state.cars.forEach((car) => car.hide());
      state.cars[action.payload].show(true);
      state.focusedCar = state.cars[action.payload].id;
    },
    unfocusCar: (state) => {
      state.cars.forEach((car) => car.show());
      state.focusedCar = null;
    },
    addMessage: (
      state,
      action: PayloadAction<{ id: number; message: any }>
    ) => {
      if (!state.carsReceivedMessages[action.payload.id]) return;
      state.carsReceivedMessages[action.payload.id].push(
        action.payload.message
      );
    },
  },
});

// export const { increment, decrement, incrementByAmount } =
//   simulationSlice.actions;

export const {
  initRSUs,
  addCar,
  addRSU,
  clearState,
  focusCar,
  unfocusCar,
  addMessage,
} = simulationSlice.actions;

export default simulationSlice.reducer;
