import {Tabs} from 'antd';
import {useQuery} from 'react-query';
import ProviderService from "@/services/providerService";
import React from "react";
import {NotificationType} from "@/utils/Constants";
import ProviderListByType from "@/components/provider/ProviderListByType";

export default function ProviderList({appId}: { appId: number }) {
  const providerService = new ProviderService();

  const {data: emailProvidersResp} = useQuery(['getAllEmailProvidersOfApp'], () =>
    providerService.getProvidersByAppIdAndType(appId, NotificationType.EMAIL)
  );

  const {data: smsProvidersResp} = useQuery(['getAllSmsProvidersOfApp'], () =>
    providerService.getProvidersByAppIdAndType(appId, NotificationType.SMS)
  );

  const {data: pushProvidersResp} = useQuery(['getAllPushProvidersOfApp'], () =>
    providerService.getProvidersByAppIdAndType(appId, NotificationType.PUSH)
  );

  const {data: webhookProvidersResp} = useQuery(['getAllWebhookProvidersOfApp'], () =>
    providerService.getProvidersByAppIdAndType(appId, NotificationType.WEBHOOK)
  );

  const tabs = [
    {
      label: 'Email',
      key: 'email',
      children: <ProviderListByType allProvidersResp={emailProvidersResp} appId={appId}
                                    notificationType={NotificationType.EMAIL}/>,
    },
    {
      label: 'SMS',
      key: 'sms',
      children: <ProviderListByType allProvidersResp={smsProvidersResp} appId={appId}
                                    notificationType={NotificationType.SMS}/>
    },
    {
      label: 'Push',
      key: 'push',
      children: <ProviderListByType allProvidersResp={pushProvidersResp} appId={appId}
                                    notificationType={NotificationType.PUSH}/>
    },
    {
      label: 'Webhook',
      key: 'webhook',
      children: <ProviderListByType allProvidersResp={webhookProvidersResp} appId={appId}
                                    notificationType={NotificationType.WEBHOOK}/>
    },
  ];

  return (
    <div>
      <Tabs
        defaultActiveKey='1'
        type='card'
        size='middle'
        items={tabs.map((tab, i) => {
          return {
            label: tab.label,
            key: tab.key,
            children: tab.children,
          };
        })}
      />
    </div>
  );
}

