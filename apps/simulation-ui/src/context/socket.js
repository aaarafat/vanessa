import { createContext } from 'react';
import io from 'socket.io-client';

export const socket = io('http://localhost:65432');
export const SocketContext = createContext();
export const SocketProvider = ({ children }) => (
  <SocketContext.Provider value={socket}>{children}</SocketContext.Provider>
);
