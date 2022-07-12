import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Car } from '@vanessa/utils';
import * as turf from '@turf/turf';

export interface CarState {
  car: Car | undefined;
  obstacles: turf.Feature<turf.Point>[];
  messages: string[];
}

const initialState: CarState = {
  car: undefined,
  obstacles: [],
  messages: [],
};

export const carSlice = createSlice({
  name: 'car',
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
  carSlice.actions;

export default carSlice.reducer;
