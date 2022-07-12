import { createSlice, PayloadAction } from '@reduxjs/toolkit';

export type ArpEntry = {
  ip: string;
  mac: string;
};

export type ObstacleTableEntry = {
  lat: number;
  lng: number;
};

export type ReceivedPackets = {
  receivedFromRsus: number;
  receivedFromCars: number;
};

export type SentPackets = {
  sentToRsus: number;
  sentToCars: number;
};

export interface RsuState extends ReceivedPackets, SentPackets {
  id: number;
  arp: ArpEntry[];
  obstacles: ObstacleTableEntry[];
}

const initialState: RsuState = {
  id: 0,
  arp: [],
  obstacles: [],
  receivedFromRsus: 0,
  receivedFromCars: 0,
  sentToRsus: 0,
  sentToCars: 0,
};

export const rsuSlice = createSlice({
  name: 'rsu',
  initialState,
  reducers: {
    initRsu: (state, action: PayloadAction<RsuState>) => {
      state.obstacles = action.payload.obstacles || [];
      state.arp = action.payload.arp;
      state.id = action.payload.id;
      state.receivedFromCars = action.payload.receivedFromCars;
      state.receivedFromRsus = action.payload.receivedFromRsus;
      state.sentToCars = action.payload.sentToCars;
      state.sentToRsus = action.payload.sentToRsus;
    },
    addObstacle: (state, action: PayloadAction<ObstacleTableEntry>) => {
      console.log(action);
      state.obstacles = [...state.obstacles, action.payload];
    },
    addArpEntry: (state, action: PayloadAction<ArpEntry>) => {
      console.log(action);
      state.arp = [...state.arp, action.payload];
    },
    removeArpEntry: (state, action: PayloadAction<ArpEntry>) => {
      console.log(action);
      state.arp = state.arp.filter((entry) => entry.ip !== action.payload.ip);
    },
    updateArpEntry: (state, action: PayloadAction<ArpEntry>) => {
      console.log(action);
      state.arp = state.arp.map((entry) => {
        if (entry.ip === action.payload.ip) {
          return action.payload;
        }
        return entry;
      });
    },
    updateReceivedPackets: (state, action: PayloadAction<ReceivedPackets>) => {
      state.receivedFromRsus = action.payload.receivedFromRsus;
      state.receivedFromCars = action.payload.receivedFromCars;
    },
    updateSentPackets: (state, action: PayloadAction<SentPackets>) => {
      console.log(action);
      state.sentToRsus = action.payload.sentToRsus;
      state.sentToCars = action.payload.sentToCars;
    },
  },
});

export const {
  addObstacle,
  addArpEntry,
  updateReceivedPackets,
  updateSentPackets,
  initRsu,
  removeArpEntry,
  updateArpEntry,
} = rsuSlice.actions;

export default rsuSlice.reducer;
