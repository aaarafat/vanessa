import React from 'react';
import { useAppSelector } from '../store';
import { TableCard } from './table-card';

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
