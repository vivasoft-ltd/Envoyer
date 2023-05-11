import {Badge, Descriptions, Tooltip} from 'antd';
import React, {useEffect, useState} from 'react';
import {ProviderType, SMS_Policy} from '@/utils/Constants';
import {Mobile_code} from '@/utils/CountryList';
import {InfoCircleTwoTone} from "@ant-design/icons";

export default function ProviderInfo({
                                       providerDetails,
                                     }: {
  providerDetails: ProviderResponse;
}) {
  const [tempConfig, setTempConfig] = useState({});

  useEffect(() => {
    setTempConfig(providerDetails.config);
  }, [providerDetails]);

  const getLabel = (value:string) => {
    const item = SMS_Policy.find((item) => item.value === value);
    return item ? item.label : '';
  };

  const getLabelForCode = (value:number) => {
    const item = Mobile_code.find((item) => item.value === value);
    return item ? item.label + " (+"+ item.value+")" : '';
  };

  const getLabelsForValues = (values: string[]) => {
    return values.map((value) => getLabelForCode(Number(value)));
  };

  return (
    <Descriptions title='' layout='horizontal' column={2} bordered>
      <Descriptions.Item label='Type'>
        {providerDetails?.type}
      </Descriptions.Item>
      <Descriptions.Item label='Provider Type'>
        {providerDetails?.provider_type}
      </Descriptions.Item>
      <Descriptions.Item label='Status'>
        <Badge
          status={providerDetails?.active ? 'success' : 'error'}
          text={providerDetails?.active ? 'Active' : 'Inactive'}
        />
      </Descriptions.Item>
      <Descriptions.Item label='Priority'>
        {providerDetails?.priority === 0 ? '' : providerDetails?.priority}
      </Descriptions.Item>
      <Descriptions.Item label='Provider Name' span={2}>
        {providerDetails?.name}
      </Descriptions.Item>
      <Descriptions.Item label='Description' span={2}>
        {providerDetails?.description}
      </Descriptions.Item>
      <Descriptions.Item label='Config' span={2}>
        {providerDetails?.provider_type === ProviderType.SMTP ? (
          <div>
            <Descriptions
              title={ProviderType.SMTP}
              layout='horizontal'
              column={1}
              bordered
            >
              <Descriptions.Item label='SMTP Host'>
                {(tempConfig as SmtpConfig).smtp_host}
              </Descriptions.Item>
              <Descriptions.Item label='SMTP Port'>
                {(tempConfig as SmtpConfig).smtp_port}
              </Descriptions.Item>
              <Descriptions.Item label='SMTP Username'>
                {(tempConfig as SmtpConfig).smtp_user_name}
              </Descriptions.Item>
              <Descriptions.Item label='SMTP Password'>
                {(tempConfig as SmtpConfig).smtp_password}
              </Descriptions.Item>
              <Descriptions.Item label='Sender'>
                {(tempConfig as SmtpConfig).sender}
              </Descriptions.Item>
            </Descriptions>
          </div>
        ) : providerDetails?.provider_type === ProviderType.TWILIO ? (
          <div>
            <Descriptions
              title={ProviderType.TWILIO}
              layout='horizontal'
              column={1}
              bordered
            >
              <Descriptions.Item label='Account SID'>
                {(tempConfig as TwilioSmsConfig).account_sid}
              </Descriptions.Item>
              <Descriptions.Item label='Auth Token'>
                {(tempConfig as TwilioSmsConfig).auth_token}
              </Descriptions.Item>
              <Descriptions.Item label='Sender Id'>
                {(tempConfig as TwilioSmsConfig).sender_id}
              </Descriptions.Item>
            </Descriptions>
          </div>
        ): providerDetails?.provider_type === ProviderType.VONAGE ? (
          <div>
            <Descriptions
              title={ProviderType.VONAGE}
              layout='horizontal'
              column={1}
              bordered
            >
              <Descriptions.Item label='Api Key'>
                {(tempConfig as VonageSmsConfig).api_key}
              </Descriptions.Item>
              <Descriptions.Item label='Api Secret'>
                {(tempConfig as VonageSmsConfig).api_secret}
              </Descriptions.Item>
              <Descriptions.Item label='Sender Id'>
                {(tempConfig as VonageSmsConfig).sender_id}
              </Descriptions.Item>
            </Descriptions>
          </div>
        ) : providerDetails?.provider_type === ProviderType.FIREBASE ? (
          <div>
            <Descriptions
              className='overflow-scroll'
              title={ProviderType.FIREBASE}
              layout='vertical'
              column={1}
              bordered
            >
              <Descriptions.Item
                className='whitespace-pre-wrap'
                label='Service Account JSON'
              >
                {JSON.stringify(providerDetails.config, null, 2)}
              </Descriptions.Item>
            </Descriptions>
          </div>
        ) : providerDetails?.provider_type === ProviderType.WEBHOOK ? (
          <div>
            <Descriptions
              title={ProviderType.WEBHOOK}
              layout='horizontal'
              column={1}
              bordered
            >
              <Descriptions.Item label='Url'>
                {(tempConfig as WebhookConfig).url}
              </Descriptions.Item>
              <Descriptions.Item label='Bearer Token'>
                {(tempConfig as WebhookConfig).token}
              </Descriptions.Item>
            </Descriptions>
          </div>
        ) : (
          <div></div>
        )}
      </Descriptions.Item>
      {providerDetails.policy &&
          <Descriptions.Item label="Policy" span={2}>
              <div>
                  <Descriptions
                      layout='horizontal'
                      column={1}
                      bordered
                      title={<div className="flex items-end gap-1">Policy<Tooltip title="All policy must satisfy and in each policy minimum one condition must satisfy"><InfoCircleTwoTone /></Tooltip></div>}
                  >
                    {Object.entries(providerDetails.policy).map(([key, value]) => (
                      <Descriptions.Item key={key} label={getLabel(key)}>
                        {key === "receiver.prefix" ? JSON.stringify(value,null, 2) :
                          JSON.stringify(getLabelsForValues(value),null,2 )}
                      </Descriptions.Item>
                    ))}
                  </Descriptions>
              </div>
          </Descriptions.Item>
      }
    </Descriptions>
  );
}
