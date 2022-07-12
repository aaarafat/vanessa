import React, { useEffect, useState } from 'react';
import { useAppDispatch } from '../store';
import {
  addArpEntry,
  addObstacle,
  ArpEntry,
  initRsu,
  ObstacleTableEntry,
  ReceivedPackets,
  removeArpEntry,
  RsuState,
  SentPackets,
  updateArpEntry,
  updateReceivedPackets,
  updateSentPackets,
} from '../store/rsuSlice';

type RefreshData = RsuState;
type AddObstacleData = ObstacleTableEntry;
type AddArpEntryData = ArpEntry;
type RemoveArpEntryData = ArpEntry;
type UpdateArpEntryData = ArpEntry;
type UpdateReceivedPacketsData = ReceivedPackets;
type UpdateSentPackets = SentPackets;

export const useEventSource = () => {
  const [eventSource, setEventSource] = useState<EventSource | null>(null);
  const dispatch = useAppDispatch();

  useEffect(() => {
    if (!eventSource) return;
    eventSource.addEventListener('refresh', ({ data: message }) => {
      const json: RefreshData = JSON.parse(message).data;
      console.log('Refresh!!!!!!!!!!!!!!');
      dispatch(initRsu(json));
    });

    eventSource.addEventListener('add-obstacle', ({ data: message }) => {
      const json: AddObstacleData = JSON.parse(message).data;
      dispatch(addObstacle(json));
    });

    eventSource.addEventListener('add-arp-entry', ({ data: message }) => {
      const json: AddArpEntryData = JSON.parse(message).data;
      dispatch(addArpEntry(json));
    });

    eventSource.addEventListener('remove-arp-entry', ({ data: message }) => {
      const json: RemoveArpEntryData = JSON.parse(message).data;
      dispatch(removeArpEntry(json));
    });

    eventSource.addEventListener('update-arp-entry', ({ data: message }) => {
      const json: UpdateArpEntryData = JSON.parse(message).data;
      dispatch(updateArpEntry(json));
    });

    eventSource.addEventListener(
      'update-received-packets',
      ({ data: message }) => {
        const json: UpdateReceivedPacketsData = JSON.parse(message).data;
        dispatch(updateReceivedPackets(json));
      }
    );

    eventSource.addEventListener('update-sent-packets', ({ data: message }) => {
      const json: UpdateSentPackets = JSON.parse(message).data;
      dispatch(updateSentPackets(json));
    });

    return () => eventSource.close();
  }, [eventSource, dispatch]);

  return [eventSource, setEventSource] as const;
};
