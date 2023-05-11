import {useMutation, useQueryClient} from 'react-query';
import React, {useEffect, useState} from 'react';
import {toast} from 'react-toastify';
import {Button, Form, FormListFieldData, Input, InputNumber, Select, Switch, Tooltip, Typography} from 'antd';
import {useRouter} from 'next/router';
import ProviderService from '@/services/providerService';
import {NotificationType, ProviderType} from '@/utils/Constants';
import {Mobile_code} from '@/utils/CountryList';
import {InfoCircleTwoTone, MinusCircleOutlined, PlusOutlined} from "@ant-design/icons";

export default function ProviderCreate({
                                         hideModal,
                                         IsEdit = false,
                                         EditProviderDetails,
                                       }: {
  hideModal: () => void;
  IsEdit?: boolean;
  EditProviderDetails?: ProviderResponse;
}) {
  const queryClient = useQueryClient();
  const providerService = new ProviderService();
  const [form] = Form.useForm();

  const router = useRouter();
  const app_id = router.query?.['app-id'] as string;

  const [name, setName] = useState('');
  const [priority, setPriority] = useState(0);
  const [active, setActive] = useState(false);
  const [description, setDescription] = useState('');
  const [notificationType, setNotificationType] = useState(
    NotificationType.SMS
  );
  const [providerType, setProviderType] = useState('');

  //twilio
  const [senderId, setSenderId] = useState('');
  const [accountSid, setAccountSid] = useState('')
  const [authToken, setAuthToken] = useState('')

  //vonage
  const [apiKey, setApiKey] = useState('');
  const [apiSecret, setApiSecret] = useState('')

  //smtp
  const [smtpUsername, setSmtpUsername] = useState('');
  const [smtpPassword, setSmtpPassword] = useState('');
  const [smtpPort, setSmtpPort] = useState('');
  const [smtpHost, setSmtpHost] = useState('');
  const [smtpSender, setSmtpSender] = useState('');

  //fcm
  const [fcmConfig, setFcmConfig] = useState<string>();

  //webhook
  const [webhookUrl, setWebhookUrl] = useState('');
  const [token, setToken] = useState('');

  useEffect(() => {
    if (IsEdit && EditProviderDetails != undefined) {
      setName(EditProviderDetails?.name);
      setDescription(EditProviderDetails?.description);
      setProviderType(EditProviderDetails.provider_type);
      setNotificationType(EditProviderDetails.type);
      setActive(EditProviderDetails.active);
      setPriority(EditProviderDetails.priority);
      if (EditProviderDetails.provider_type === ProviderType.TWILIO) {
        let tempConfig = EditProviderDetails.config as TwilioSmsConfig;
        setAccountSid(tempConfig.account_sid);
        setAuthToken(tempConfig.auth_token);
        setSenderId(tempConfig.sender_id);
      } else if (EditProviderDetails.provider_type === ProviderType.VONAGE) {
        let tempConfig = EditProviderDetails.config as VonageSmsConfig;
        setApiKey(tempConfig.api_key);
        setApiSecret(tempConfig.api_secret);
        setSenderId(tempConfig.sender_id);
      } else if (EditProviderDetails.provider_type === ProviderType.SMTP) {
        let tempConfig = EditProviderDetails.config as SmtpConfig;
        setSmtpSender(tempConfig.sender);
        setSmtpHost(tempConfig.smtp_host);
        setSmtpPort(tempConfig.smtp_port?.toString());
        setSmtpUsername(tempConfig.smtp_user_name);
        setSmtpPassword(tempConfig.smtp_password);
      } else if (EditProviderDetails.provider_type === ProviderType.FIREBASE) {
        setFcmConfig(JSON.stringify(EditProviderDetails.config));
      } else if (EditProviderDetails.provider_type === ProviderType.WEBHOOK) {
        let tempConfig = EditProviderDetails.config as WebhookConfig;
        setWebhookUrl(tempConfig.url);
        setToken(tempConfig.token || '');
      }

      let values = EditProviderDetails.policy || {}
      const initialValues = Object.entries(values).flatMap(([policyType, rules]) => {
        let acc: any[] = []
        rules.forEach((rule: string, index: number) => {
          acc.push({
            rule: rule,
            policyType: policyType,
          });
        });
        return acc;
      });
      if (initialValues.length != 0){
        form.setFieldsValue({ policy: initialValues });
      }
    }
  }, [IsEdit, EditProviderDetails]);

  const clearForm = () => {
    setName('');
    setDescription('');
    setProviderType('');
    setNotificationType(NotificationType.SMS);
    setActive(false);
    setPriority(0);
    setSenderId('');
    setApiKey('');
    setSmtpSender('');
    setSmtpHost('');
    setSmtpPort('');
    setSmtpUsername('');
    setSmtpPassword('');
    setFcmConfig('');
    setWebhookUrl('');
    setToken('');
    setApiSecret('');
    setAccountSid('');
    setAuthToken('');
    form.resetFields();
    form.setFieldsValue({ policy: null });
  };

  const handleSubmit = (values: any) => {
    let policies = undefined;
    if(notificationType === NotificationType.SMS) {
      policies = values?.policy?.reduce(
        (acc: Record<string, any>, row: any) => {
          const policyType = row.policyType;
          const rule = row.rule;
          acc[policyType] = acc[policyType] || [];
          acc[policyType].push(String(rule));
          return acc;
        },
        {}
      );
    }

    if (IsEdit) {
      handleEditProvider(policies);
    } else {
      handleCreateProvider(policies);
    }
  };

  function handleCreateProvider(policies?: object) {
    let ProviderConfig: object;

    if (providerType == ProviderType.TWILIO) {
      ProviderConfig = {
        account_sid: accountSid,
        auth_token: authToken,
        sender_id: senderId,
      } as TwilioSmsConfig;
    } else if (providerType == ProviderType.VONAGE) {
      ProviderConfig = {
        api_key: apiKey,
        api_secret: apiSecret,
        sender_id: senderId,
      } as VonageSmsConfig;
    } else if (providerType == ProviderType.SMTP) {
      ProviderConfig = {
        sender: smtpSender,
        smtp_host: smtpHost,
        smtp_port: Number(smtpPort),
        smtp_password: smtpPassword,
        smtp_user_name: smtpUsername,
      } as SmtpConfig;
    } else if (providerType == ProviderType.FIREBASE) {
      ProviderConfig = JSON.parse(fcmConfig as string);
    } else if (providerType == ProviderType.WEBHOOK) {
      ProviderConfig = {
        url: webhookUrl,
        token: token,
      } as WebhookConfig;
    } else {
      toast.error('Select Provider type');
      return;
    }

    createProvider(
      {
        name: name,
        description: description,
        active: active,
        app_id: Number(app_id),
        type: notificationType,
        provider_type: providerType,
        config: ProviderConfig!,
        priority: priority,
        policy:policies,
      },
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error creating provider');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllSmsProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllEmailProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllPushProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllWebhookProvidersOfApp']);
          toast.success('Provider created successfully');
          clearForm();
        },
      }
    );
  }

  function handleEditProvider(policies?: object) {
    let ProviderConfig: object;

    if (providerType == ProviderType.TWILIO) {
      ProviderConfig = {
        account_sid: accountSid,
        auth_token: authToken,
        sender_id: senderId,
      } as TwilioSmsConfig;
    } else if (providerType == ProviderType.VONAGE) {
      ProviderConfig = {
        api_key: apiKey,
        api_secret: apiSecret,
        sender_id: senderId,
      } as VonageSmsConfig;
    } else if (providerType == ProviderType.SMTP) {
      ProviderConfig = {
        sender: smtpSender,
        smtp_host: smtpHost,
        smtp_port: Number(smtpPort),
        smtp_password: smtpPassword,
        smtp_user_name: smtpUsername,
      };
    } else if (providerType == ProviderType.FIREBASE) {
      ProviderConfig = JSON.parse(fcmConfig as string);
    } else if (providerType == ProviderType.WEBHOOK) {
      ProviderConfig = {
        url: webhookUrl,
        token: token,
      } as WebhookConfig;
    }

    editProvider(
      {
        name: name,
        description: description,
        active: active,
        app_id: Number(app_id),
        type: notificationType,
        provider_type: providerType,
        config: ProviderConfig!,
        priority: priority,
        policy:policies,
      },
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error editing provider');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllSmsProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllEmailProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllPushProvidersOfApp']);
          await queryClient.invalidateQueries(['getAllWebhookProvidersOfApp']);
          toast.success('Provider updated successfully');
          clearForm();
        },
      }
    );
  }

  const {mutate: editProvider} = useMutation(
    (data: CreateProviderInput) =>
      providerService.editProvider(data, EditProviderDetails?.ID),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  const {mutate: createProvider} = useMutation(
    (data: CreateProviderInput) => providerService.createProvider(data),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  return (
    <div className=''>
      <Form className='w-full max-w-lg pb-3' layout='vertical' form={form} onFinish={handleSubmit}>
        <Form.Item label='Notification Type'>
          <Select
            disabled={IsEdit}
            value={notificationType}
            onChange={(value) => {
              setProviderType('');
              setNotificationType(value);
            }}
            options={[
              {value: NotificationType.SMS, label: 'SMS'},
              {value: NotificationType.EMAIL, label: 'Email'},
              {value: NotificationType.PUSH, label: 'Push'},
              {value: NotificationType.WEBHOOK, label: 'Webhook'},
            ]}
          />
        </Form.Item>
        <Form.Item label='Provider Type'>
          <Select
            disabled={IsEdit}
            value={providerType}
            onChange={(value) => setProviderType(value)}
            options={
              notificationType === NotificationType.SMS
                ? [{value: ProviderType.TWILIO, label: 'Twilio'},
                  {value: ProviderType.VONAGE, label: 'Vonage'}]
                : notificationType === NotificationType.EMAIL
                  ? [{value: ProviderType.SMTP, label: 'SMTP'}]
                  : notificationType === NotificationType.PUSH
                    ? [{value: ProviderType.FIREBASE, label: 'FCM'}]
                    : notificationType === NotificationType.WEBHOOK
                      ? [{value: ProviderType.WEBHOOK, label: 'Webhook'}]
                      : []
            }
          />
        </Form.Item>
        <Form.Item label='Provider Name'>
          <Input value={name} onChange={(e) => setName(e.target.value)}/>
        </Form.Item>

        <Form.Item label='Description'>
          <Input.TextArea
            rows={4}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </Form.Item>

        {providerType === ProviderType.TWILIO && (
          <div className='border p-2'>
            <span className='font-semibold'>{ProviderType.TWILIO}</span>
            <Form.Item label='Account Sid'>
              <Input
                value={accountSid}
                onChange={(e) => setAccountSid(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='Auth Token'>
              <Input
                value={authToken}
                onChange={(e) => setAuthToken(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='Sender Id'>
              <Input
                value={senderId}
                onChange={(e) => setSenderId(e.target.value)}
              />
            </Form.Item>
          </div>
        )}

        {providerType === ProviderType.VONAGE && (
          <div className='border p-2'>
            <span className='font-semibold'>{ProviderType.VONAGE}</span>
            <Form.Item label='Api Key'>
              <Input
                value={apiKey}
                onChange={(e) => setApiKey(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='Api Secret'>
              <Input
                value={apiSecret}
                onChange={(e) => setApiSecret(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='Sender Id'>
              <Input
                value={senderId}
                onChange={(e) => setSenderId(e.target.value)}
              />
            </Form.Item>
          </div>
        )}

        {providerType === ProviderType.SMTP && (
          <div className='border p-2'>
            <span className='font-semibold'>{ProviderType.SMTP}</span>
            <Form.Item label='SMTP Host'>
              <Input
                value={smtpHost}
                onChange={(e) => setSmtpHost(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='SMTP Port'>
              <Input
                value={smtpPort}
                onChange={(e) => setSmtpPort(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='SMTP Username'>
              <Input
                value={smtpUsername}
                onChange={(e) => setSmtpUsername(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='SMTP Password'>
              <Input
                value={smtpPassword}
                onChange={(e) => setSmtpPassword(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='Sender'>
              <Input
                value={smtpSender}
                onChange={(e) => setSmtpSender(e.target.value)}
              />
            </Form.Item>
          </div>
        )}

        {providerType === ProviderType.FIREBASE && (
          <div className='border p-2'>
            <span className='font-semibold'>{ProviderType.FIREBASE}</span>
            <Form.Item label='Service Account JSON'>
              <Tooltip title="How to get Service Account JSON">
                <Typography.Link target="_blank" className={'float-right'}
                                 href="https://firebase.google.com/docs/cloud-messaging/auth-server">Need
                  Help?</Typography.Link>
              </Tooltip>
              <Input.TextArea
                rows={4}
                value={fcmConfig}
                onChange={(e) => setFcmConfig(e.target.value)}
              />
            </Form.Item>
          </div>
        )}

        {providerType === ProviderType.WEBHOOK && (
          <div className='border p-2'>
            <span className='font-semibold'>{ProviderType.WEBHOOK}</span>
            <Form.Item label='Url'>
              <Input
                value={webhookUrl}
                onChange={(e) => setWebhookUrl(e.target.value)}
              />
            </Form.Item>
            <Form.Item label='Bearer Token'>
              <Input
                value={token}
                onChange={(e) => setToken(e.target.value)}
              />
            </Form.Item>
          </div>
        )}

        {notificationType === NotificationType.SMS  && (
          <div className='border p-2 mt-2 mb-2'>
            <span className='font-semibold block pb-2'>
              <div className="flex items-end gap-1">
                Policy<
                Tooltip title={<div>All type of policy must satisfy and for each policy type minimum one condition must satisfy.
                <br/> Example: Receiver Country can be multiple but at least one must satisfy</div>}>
                <InfoCircleTwoTone />
              </Tooltip>
              </div>
            </span>
            <Form.List name="policy">
              {(fields, { add, remove }) => (
                <>
                  {fields.map((field, index) => (
                    <Policy form={form} key={index} field={field} remove={remove}/>
                  ))}

                  <Form.Item>
                    <Button type="dashed" onClick={() => add()} block icon={<PlusOutlined />}>
                      Add Policy
                    </Button>
                  </Form.Item>
                </>
              )}
            </Form.List>
          </div>
        )}

        <Form.Item label='Active' noStyle valuePropName='checked'>
          <span className="mr-1">Active: </span>
          <Switch
            checked={active}
            className='bg-gray-200'
            onChange={(value) => setActive(value)}
          />
        </Form.Item>

        <Form.Item className="lg:text-right">
          <Button
            htmlType="submit"
            type='primary'
            className='bg-blue-700 mt-3'
          >
            {IsEdit ? 'Update' : 'Create'}
          </Button>
        </Form.Item>
      </Form>
    </div>
  );
}

interface RowProps {
  field:  FormListFieldData;
  remove: (index: number) => void;
  form: any;
}

const Policy: React.FC<RowProps> = ({field, remove, form}) => {
  const [policyType, setPolicyType] = useState("receiver.country");

  const handleInputTypeChange = (value: "receiver.country" | "receiver.prefix") => {
    setPolicyType(value);
  };

  useEffect(()=>{
    const policy_type = form.getFieldValue(`policy`)[field.name].policyType;
    setPolicyType(policy_type)
  },[])

  return (
    <div className="w-full items-end flex gap-3">
      <Form.Item className="flex-1"
                 label="Policy Type"
                 name={[field.name, 'policyType']}
                 rules={[{ required: true, message: 'Missing policy type' }]}
                 initialValue={policyType}
      >
        <Select dropdownMatchSelectWidth={false} onChange={handleInputTypeChange} options={[
          {value:"receiver.country", label:"Receiver Country"},
          {value:"receiver.prefix", label:"Receiver Prefix"}
        ]}>
        </Select>
      </Form.Item>

      {policyType === "receiver.prefix" && (
        <Form.Item className="flex-1"
          label="Prefix"
          name={[field.name, "rule"]}
          initialValue=""
        >
          <InputNumber min={0} formatter={(value) => `${value}`.replace(/\./g, '')} className="w-full" placeholder="Prefix" />
        </Form.Item>
      )}

      {policyType === "receiver.country" && (
        <Form.Item className="flex-1"
          label="Country Code"
          name={[field.name, "rule"]}
          initialValue="880"
        >
          <Select dropdownMatchSelectWidth={false}>
            {Mobile_code.map((option) => (
              <Select.Option key={option.value} value={String(option.value)}>
                {option.label + " (+" + option.value+ ") "}
              </Select.Option>
            ))}
            </Select>
        </Form.Item>
      )}
      <Form.Item>
        <MinusCircleOutlined onClick={() => remove(field.name)} />
      </Form.Item>
    </div>
  );
}