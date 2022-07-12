import React from 'react';
import { TableCard } from './table-card';

const ARPTable = [
  {
    ip: '192.168.1.1',
    mac: '00:00:00:00:00:00',
  },
  {
    ip: '192.168.1.1',
    mac: '00:00:00:00:00:00',
  },
  {
    ip: '192.168.1.1',
    mac: '00:00:00:00:00:00',
  },
  {
    ip: '192.168.1.1',
    mac: '00:00:00:00:00:00',
  },
];

export const ARP = () => {
  return (
    <TableCard table={ARPTable} headers={['IP', 'MAC']} title="ARP Table" />
  );
};
