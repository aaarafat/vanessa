import { createContext } from 'react';
import io from 'socket.io-client';

export const socket = io('http://127.0.0.1:65432');
export const SocketContext = createContext();
export const SocketProvider = ({ children }) => (
  <SocketContext.Provider value={socket}>{children}</SocketContext.Provider>
);
