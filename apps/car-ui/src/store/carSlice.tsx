import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Car } from '@vanessa/utils';
import * as turf from '@turf/turf';

export interface SimulationState {
  car: Car | undefined;
  obstacles: turf.Feature<turf.Point>[];
  messages: string[];
}

const initialState: SimulationState = {
  car: undefined,
  obstacles: [],
  messages: [],
};

export const simulationSlice = createSlice({
  name: 'simulation',
  initialState,
  reducers: {
    initCar: (state, action: PayloadAction<Car>) => {
      const car = action.payload;
      state.car = car;
    },
    addObstacle: (state, action: PayloadAction<turf.Feature<turf.Point>>) => {
      state.obstacles = [...state.obstacles, action.payload];
    },
    addObstacles: (
      state,
      action: PayloadAction<turf.Feature<turf.Point>[]>
    ) => {
      state.obstacles = action.payload;
    },
    addMessage: (state, action: PayloadAction<string>) => {
      state.messages.unshift(action.payload);
    },
  },
});

export const { initCar, addObstacle, addObstacles, addMessage } =
  simulationSlice.actions;

export default simulationSlice.reducer;
