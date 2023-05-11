import { useQuery } from 'react-query';
import { Button, Descriptions, Input } from 'antd';
import React, { useState } from 'react';
import AppServices from '@/services/appServices';
import {
  CopyOutlined,
  EyeInvisibleOutlined,
  EyeOutlined,
} from '@ant-design/icons';
import { toast } from 'react-toastify';
import { copyToClipboard } from '@/utils/copyToClipboard';

export default function ClientInfo({
  clientDetails,
}: {
  clientDetails: ClientResponse;
}) {
  const appServices = new AppServices();
  const [visibleAppKey, setVisibleAppKey] = useState(false);
  const [visibleClientKey, setVisibleClientKey] = useState(false);

  const { data: appResp } = useQuery(['getApp'], () =>
    appServices.getApp(Number(clientDetails.app_id))
  );

  return (
    <Descriptions title='' layout='horizontal' column={1} bordered>
      <Descriptions.Item label='Client Id'>
        {clientDetails?.ID}
      </Descriptions.Item>
      <Descriptions.Item label='Client Name'>
        {clientDetails?.name}
      </Descriptions.Item>
      <Descriptions.Item label='Description'>
        {clientDetails?.description}
      </Descriptions.Item>
      <Descriptions.Item label='App Key'>
        <div className='flex items-center border pr-2'>
          <Input
            type={visibleAppKey ? 'text' : 'password'}
            bordered={false}
            disabled
            value={appResp?.app_key}
          />
          {visibleAppKey && (
            <EyeOutlined onClick={() => setVisibleAppKey(false)} />
          )}
          {!visibleAppKey && (
            <EyeInvisibleOutlined onClick={() => setVisibleAppKey(true)} />
          )}
          <CopyOutlined
            className='pl-2'
            onClick={() => {
              copyToClipboard(appResp?.app_key || '');
            }}
          />
        </div>
      </Descriptions.Item>
      <Descriptions.Item label='Client Key'>
        <div className='flex items-center border pr-2'>
          <Input
            type={visibleClientKey ? 'text' : 'password'}
            bordered={false}
            disabled
            value={clientDetails?.client_key}
          />
          {visibleClientKey && (
            <EyeOutlined onClick={() => setVisibleClientKey(false)} />
          )}
          {!visibleClientKey && (
            <EyeInvisibleOutlined onClick={() => setVisibleClientKey(true)} />
          )}
          <CopyOutlined
            className='pl-2'
            onClick={() => {
              copyToClipboard(clientDetails?.client_key || '');
            }}
          />
        </div>
      </Descriptions.Item>
    </Descriptions>
  );
}
