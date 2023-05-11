import {List} from 'antd';
import {useQuery} from 'react-query';
import EventService from "@/services/eventService";
import EventItem from "@/components/event/EventItem";

export default function EventList({appId}: { appId: number }) {
  const eventService = new EventService();

  const {data: allEventsResp} = useQuery(['getAllEvents'], () =>
    eventService.getAllEvents(appId)
  );

  return (
    <div>
      <List
        renderItem={(item) => (
          <EventItem eventDetails={item}/>
        )}
        itemLayout="horizontal"
        dataSource={allEventsResp}
      />
    </div>
  );
}
