import React from 'react';
import { useAppSelector } from '../store';
import { TableCard } from './table-card';

// const ARPTable = [
//   {
//     ip: '192.168.1.1',
//     mac: '00:00:00:00:00:00',
//   },
//   {
//     ip: '192.168.1.1',
//     mac: '00:00:00:00:00:00',
//   },
//   {
//     ip: '192.168.1.1',
//     mac: '00:00:00:00:00:00',
//   },
//   {
//     ip: '192.168.1.1',
//     mac: '00:00:00:00:00:00',
//   },
// ];

export const ARP = () => {
  const ARPTable = useAppSelector((state) => state.rsu.arp);
  return (
    <TableCard
      table={ARPTable}
      keys={['ip', 'mac']}
      headers={['IP', 'MAC']}
      title="ARP Table"
    />
  );
};
