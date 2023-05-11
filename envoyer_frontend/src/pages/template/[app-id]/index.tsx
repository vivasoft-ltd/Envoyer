import {useRouter} from 'next/router';
import React, {Fragment, ReactNode, useMemo, useRef, useState} from 'react';
import Admin from '@/layouts/Admin';
import {NotificationType, USER_ROLE} from '@/utils/Constants';
import {Languages_list} from '@/utils/CountryList';
import {GetServerSidePropsContext} from 'next';
import {getSession, useSession} from 'next-auth/react';
import {
  Button,
  Form,
  FormInstance,
  Input,
  List,
  Select,
  Switch,
} from 'antd';
import EventService from '@/services/eventService';
import {useMutation, useQuery, useQueryClient} from 'react-query';
import Link from 'next/link';
import EmailTemplate from '@/components/emailTemplate';
import {toast} from 'react-toastify';
import {CopyOutlined} from '@ant-design/icons';
import TemplateService from '@/services/templateService';
import {copyToClipboard} from '@/utils/copyToClipboard';

const {Option} = Select;

const TemplateId = () => {
  const [isShowTemplateEditor, setIsShowTemplateEditor] = useState(false);
  const formRef = React.useRef<FormInstance>(null);
  const router = useRouter();
  const {data: session} = useSession();
  const eventService = new EventService();
  const templateService = new TemplateService();

  const [subject, setSubject] = useState('');
  const [description, setDescription] = useState('');
  const [eventType, setEventType] = useState('');
  const [method, setMethod] = useState('');
  const [smsBody, setSmsBody] = useState('');
  const queryClient = useQueryClient();
  const emailEditorRef = useRef<any>(null);
  const [isActive, setIsActive] = useState(false);
  const [title, setTitle] = useState('');
  const [imageLink, setImageLink] = useState('');
  const [language, setLanguage] = useState('en');


  const appId = session?.app_id || (router.query?.['app-id'] as string);

  const {data: allEventsResp} = useQuery(['getAllEvents'], () =>
    eventService.getAllEvents(Number(appId))
  );

  const {mutate: createTemplate} = useMutation(
    (data: CreateTemplateInput) => templateService.createTemplate(data),
    {
      onSuccess: async (data) => {
        toast.success('Template created successfully');
        setIsShowTemplateEditor(false);
        formRef.current?.resetFields();
        setIsActive(false);
        setEventType('');
        setMethod('');
        setLanguage('en');
        await queryClient.invalidateQueries(['getAllTemplates']);
      },
      onSettled: async () => {
      },
    }
  );
  const selectedEvent = allEventsResp?.find(
    (event) => event.name === eventType
  );

  const createTemplateFunc = ({
                                design,
                                html,
                              }: {
    html?: string;
    design?: string;
  }) => {
    createTemplate({
      event_id: selectedEvent?.ID as number,
      type: method,
      description: description,
      message: smsBody,
      email_subject: subject,
      email_rendered_html: html,
      markup: JSON.stringify(design),
      active: isActive,
      title: title,
      file: imageLink,
      language: language,
    });
  };

  const onFinish = async (values: unknown) => {
    if (method === NotificationType.EMAIL) {
      emailEditorRef?.current?.editor.exportHtml((data: any) => {
        const {design, html} = data;
        createTemplateFunc({design, html});
      });
      return;
    }
    createTemplateFunc({});
  };

  const onReset = () => formRef.current?.resetFields();

  return (
    <Fragment>
      <div className='text-right'>
        <Button
          type='primary'
          className='bg-blue-700'
          onClick={() => {
            setIsShowTemplateEditor((prevState) => !prevState);
          }}
        >
          Create Template
        </Button>
      </div>
      <div>
        {isShowTemplateEditor && (
          <div className='my-3'>
            <Form
              ref={formRef}
              name='control-ref'
              onFinish={onFinish}
              className='max-w-3xl mx-auto'
              layout='vertical'
            >
              <div className='grid grid-cols-5 w-3/4 gap-x-3'>
                <Form.Item
                  label='Event Type'
                  name='Event Type'
                  rules={[{required: true}]}
                  className='col-span-2'
                >
                  <Select onChange={setEventType}>
                    {allEventsResp?.map((event) => {
                      return (
                        <Option value={event.name} key={event.ID}>
                          {event.name}
                        </Option>
                      );
                    })}
                  </Select>
                </Form.Item>
                <Form.Item
                  name='Notification Type'
                  label='Notification Type'
                  rules={[{required: true}]}
                  className='col-span-2'
                >
                  <Select onChange={setMethod}>
                    <Option value={NotificationType.EMAIL}>Email</Option>
                    <Option value={NotificationType.SMS}>SMS</Option>
                    <Option value={NotificationType.PUSH}>Push</Option>
                  </Select>
                </Form.Item>
                <Form.Item
                  name='language'
                  label='Language'
                  className='col-span-1'
                >
                  <Select dropdownMatchSelectWidth={false} optionLabelProp="value" defaultValue={"en"}
                          options={Languages_list} value={language}
                          onChange={setLanguage}/>
                </Form.Item>
              </div>

              <div style={{width: '75%'}}>
                <Form.Item name='description' label='Description'>
                  <Input.TextArea
                    onChange={(e) => setDescription(e.target.value)}
                  />
                </Form.Item>
                {selectedEvent && (
                  <>
                    <h3 className='mb-1.5'>Usable Variables</h3>

                    <div className='flex gap-5 flex-wrap pb-5'>
                      {selectedEvent?.variables?.map((variable) => {
                        return (
                          <Button
                            type='dashed'
                            key={variable}
                            size='large'
                            className='flex items-center'
                            onClick={() => {
                              copyToClipboard(variable);
                            }}
                          >
                            {variable}
                            <CopyOutlined/>
                          </Button>
                        );
                      })}
                    </div>
                  </>
                )}
                {method === 'email' && (
                  <Form.Item name='subject' label='Subject'>
                    <Input
                      onChange={(event) => setSubject(event.target.value)}
                    />
                  </Form.Item>
                )}
                {method === NotificationType.PUSH && (
                  <div>
                    <Form.Item name='title' label='Title'>
                      <Input
                        onChange={(event) => setTitle(event.target.value)}
                      />
                    </Form.Item>
                    <Form.Item name='sms-body' label='Body'>
                      <Input.TextArea
                        onChange={(e) => setSmsBody(e.target.value)}
                        rows={4}
                      />
                    </Form.Item>
                    <Form.Item name='imageLink' label='Image URL (optional)'>
                      <Input
                        onChange={(e) => setImageLink(e.target.value)}
                      />
                    </Form.Item>
                  </div>
                )}
              </div>

              {method === 'sms' && (
                <Form.Item
                  name='sms-body'
                  isList
                  label='SMS Body'
                  className='w-3/4'
                >
                  <Input.TextArea
                    onChange={(e) => setSmsBody(e.target.value)}
                    rows={4}
                  />
                </Form.Item>
              )}
            </Form>
            {method === 'email' && (
              <div className='mb-5'>
                <h3 className='text-base pb-3 max-w-3xl mx-auto'>Email Body</h3>
                <EmailTemplate ref={emailEditorRef}/>
              </div>
            )}
            <div className='max-w-3xl mx-auto'>
              <label htmlFor='isActive' className='cursor-pointer pr-2'>
                Active
              </label>
              <Switch
                onChange={setIsActive}
                id='isActive'
                checked={isActive}
                className='bg-gray-500'
              />
            </div>
            <Button
              type='primary'
              className='bg-blue-700 ml-auto block mb-5 mt-5'
              htmlType='submit'
              onClick={() => {
                formRef.current?.submit();
              }}
            >
              Create
            </Button>
          </div>
        )}

        <div className='pt-3'>
          <h2 className='text-2xl font-semibold pb-2'>All Events</h2>
          <List
            renderItem={(item) => (
              <div style={{border: '1px solid rgba(5, 5, 5, 0.06)'}}>
                <List.Item>
                  <List.Item.Meta
                    title={
                      <Link href={'/template/' + appId + '/event/' + item.ID}>
                        {item.name}
                      </Link>
                    }
                    description={item.description}
                  />
                </List.Item>
              </div>
            )}
            itemLayout='horizontal'
            dataSource={allEventsResp}
          />
        </div>
      </div>
    </Fragment>
  );
};

export default TemplateId;

TemplateId.pageOptions = {
  requiresAuth: true,
  getLayout: (children: ReactNode) => <Admin>{children}</Admin>,
};

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getSession(context);

  if (
    session?.role === USER_ROLE.SUPER_ADMIN ||
    session?.app_id === Number(context.params?.['app-id'])
  ) {
    return {
      props: {},
    };
  } else {
    return {
      redirect: {
        destination: '/template/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
