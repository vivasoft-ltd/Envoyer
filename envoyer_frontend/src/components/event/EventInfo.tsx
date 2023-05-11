import {Descriptions, Input, List} from "antd";
import React from "react";

export default function EventInfo({eventDetails}: { eventDetails: EventResponse }) {

  return (
    <Descriptions title="" layout="horizontal" column={1} bordered>
      <Descriptions.Item label="Event Id">{eventDetails?.ID}</Descriptions.Item>
      <Descriptions.Item label="Event Name">{eventDetails?.name}</Descriptions.Item>
      <Descriptions.Item label="Description">{eventDetails?.description}</Descriptions.Item>
      <Descriptions.Item label="Usable Variables">
        <List
          bordered
          dataSource={eventDetails?.variables || []}
          renderItem={(item) => <List.Item>{item}</List.Item>}
        />
      </Descriptions.Item>
    </Descriptions>
  );
}