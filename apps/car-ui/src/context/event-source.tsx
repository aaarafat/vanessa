import React, { createContext, useEffect, useState } from 'react';

export const EventSourceContext = createContext<
  [EventSource | null, React.Dispatch<React.SetStateAction<EventSource | null>>]
>(null as any);
export const EventSourceProvider: React.FC<React.ReactNode> = ({
  children,
}) => {
  const [eventSource, setEventSource] = useState<EventSource | null>(null);
  return (
    <EventSourceContext.Provider value={[eventSource, setEventSource]}>
      {children}
    </EventSourceContext.Provider>
  );
};
