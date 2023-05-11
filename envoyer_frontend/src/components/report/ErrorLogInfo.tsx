import {Descriptions} from "antd";
import React from "react";
import {format} from "date-fns";
import {useQuery} from "react-query";
import ProviderService from "@/services/providerService";

export default function ErrorLogInfo({errorDetails}: { errorDetails: ErrorLog }) {

  const providerService = new ProviderService();

  function IsEmpty(data: any) {
    return (data?.length === 1 && Object.keys(data[0]).length === 0);
  }

  const { data: providerResp } = useQuery(['getProvider'], () =>
    providerService.getProvider(Number(errorDetails.provider_id))
  );

  return (
    <Descriptions title="" layout="horizontal" column={1} bordered>
      <Descriptions.Item label="Datetime">{format(new Date(errorDetails.date!), 'yyyy-MM-dd hh:mm:ss a z')}</Descriptions.Item>
      { errorDetails.event_name &&
          <Descriptions.Item label="Event Name">{errorDetails?.event_name}</Descriptions.Item>}
      {errorDetails.provider_id &&
          <Descriptions.Item label="Provider Name">{providerResp?.name}</Descriptions.Item>}
      <Descriptions.Item label="Error Message">{errorDetails?.message}</Descriptions.Item>
      <Descriptions.Item label="Retry">{errorDetails?.is_requeue === true ? "True" : "False"}</Descriptions.Item>
      <Descriptions.Item
        className='whitespace-pre-wrap'
        label='Request'
      >
        {JSON.stringify(errorDetails.request, null, 2)}
      </Descriptions.Item>
      {!IsEmpty(errorDetails?.data) && <Descriptions.Item
        className='whitespace-pre-wrap'
        label='Notification Data'
      >
        {JSON.stringify(errorDetails.data, null, 2)}
      </Descriptions.Item>
      }
    </Descriptions>
  );
}