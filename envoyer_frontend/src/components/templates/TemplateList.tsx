import {
  Badge,
  Button,
  Card,
  Col,
  Form,
  FormInstance,
  Input,
  Modal,
  Row,
  Switch,
  Tabs,
} from 'antd';
import {useMutation, useQuery, useQueryClient} from 'react-query';
import TemplateService from '@/services/templateService';
import React, {useEffect, useRef, useState} from 'react';
import parse from 'html-react-parser';
import {NotificationType} from '@/utils/Constants';
import EventService from '@/services/eventService';
import {useRouter} from 'next/router';
import {copyToClipboard} from '@/utils/copyToClipboard';
import {
  CopyOutlined,
  EyeOutlined,
  DeleteOutlined,
  EditOutlined,
} from '@ant-design/icons';
import {toast} from 'react-toastify';
import EmailEditor from 'react-email-editor';
import TemplateItems from './TemplateItems';

const templateService = new TemplateService();

export default function TemplateList({eventId}: { eventId: number }) {
  const {data: allTemplatesResp} = useQuery(['getAllTemplates'], () =>
    templateService.getAllTemplates(eventId)
  );

  const tabs = [
    {
      label: 'Email',
      key: 'email',
      children: (
        <TemplateItems
          templates={allTemplatesResp}
          type={NotificationType.EMAIL}
        />
      ),
    },
    {
      label: 'SMS',
      key: 'sms',
      children: (
        <TemplateItems
          templates={allTemplatesResp}
          type={NotificationType.SMS}
        />
      ),
    },
    {
      label: 'Push',
      key: 'push',
      children: (
        <TemplateItems
          templates={allTemplatesResp}
          type={NotificationType.PUSH}
        />
      ),
    },
  ];

  return (
    <div className='pt-3'>
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
