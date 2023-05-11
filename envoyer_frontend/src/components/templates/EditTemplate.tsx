import {Button, Form, FormInstance, Input, Select, Spin, Switch} from 'antd';
import {useMutation, useQuery, useQueryClient} from 'react-query';
import TemplateService from '@/services/templateService';
import React, {useEffect, useRef, useState} from 'react';
import {Languages_list} from '@/utils/CountryList';
import {NotificationType} from '@/utils/Constants';
import EventService from '@/services/eventService';
import {useRouter} from 'next/router';
import {copyToClipboard} from '@/utils/copyToClipboard';
import {CopyOutlined} from '@ant-design/icons';
import {toast} from 'react-toastify';
import EmailEditor from 'react-email-editor';

const templateService = new TemplateService();

interface EditTemplateProps {
  template?: TemplateResponse;
  method: string;
  hideModal: () => void;
}

const EditTemplate = ({method, template, hideModal}: EditTemplateProps) => {
  const formRef = React.useRef<FormInstance>(null);
  const eventService = new EventService();
  const queryClient = useQueryClient();
  const router = useRouter();

  const emailEditorRef = useRef<any>(null);

  const [smsBody, setSmsBody] = useState(template?.message || '');
  const [description, setDescription] = useState(template?.description || '');
  const [isActive, setIsActive] = useState(template?.active || false);
  const [subject, setSubject] = useState(template?.email_subject || '');
  const [title, setTitle] = useState(template?.title || '');
  const [imageLink, setImageLink] = useState(template?.file || '');
  const [language, setLanguage] = useState(template?.language || 'en');

  const [isShowEmailEditor, setIsShowEmailEditor] = useState(false);

  useEffect(() => {
    setTimeout(() => {
      setIsShowEmailEditor(true);
    });
  }, [template]);

  const eventId = router.query?.['event-id'] as string;

  useEffect(() => {
    if (template) {
      setDescription(template?.description || '');
      setSmsBody(template?.message || '');
      setIsActive(template?.active || false);
      setSubject(template?.email_subject || '');
      setImageLink(template?.file || '');
      setTitle(template?.title || '')
      setLanguage(template?.language || 'en');
      emailEditorRef?.current?.editor?.loadDesign(
        JSON.parse(template?.markup || '')
      );
    }
  }, [template]);

  const {mutate: editTemplate} = useMutation(
    (data: UpdateTemplateInput) => templateService.editTemplate(data),
    {
      onSuccess: async (data) => {
        toast.success('Template updated successfully');
        await queryClient.invalidateQueries(['getAllTemplates']);
        hideModal();
        clearForm();
      },
      onError: async (e: Error) => {
        toast.error(e.message);
      },
      onSettled: async () => {
      },
    }
  );

  const handleUpdate = async (values: unknown) => {
    if (!template) return;
    const updateInput = {
      id: template.ID,
      markup: template.markup || '',
      event_id: Number(eventId),
      type: method,
      description,
      email_rendered_html: template.email_rendered_html || '',
      message: smsBody,
      active: isActive,
      email_subject: subject,
      title: title,
      file: imageLink,
      language: language,
    };
    if (method === NotificationType.EMAIL) {
      emailEditorRef?.current?.editor?.exportHtml((data: any) => {
        const {design, html} = data;
        if (design && html) {
          updateInput.markup = JSON.stringify(design);
          updateInput.email_rendered_html = html;
        }
        editTemplate(updateInput);
      });
      return;
    }
    editTemplate(updateInput);
  };

  const {data: getEven} = useQuery(['getEvent', eventId], () =>
    eventService.getEvent(Number(eventId))
  );

  const onLoad = () => {
    emailEditorRef.current?.editor?.loadDesign(
      JSON.parse(template?.markup || '')
    );
  };

  const onReady = () => {
    emailEditorRef.current?.editor?.loadDesign(
      JSON.parse(template?.markup || '')
    );
  };

  const clearForm = () => {
    formRef.current?.resetFields();
  };

  return (
    <div className='my-3'>
      <Form onFinish={handleUpdate} layout='vertical'>
        <Form.Item label='Description'>
          <Input.TextArea
            onChange={(e) => setDescription(e.target.value)}
            value={description}
          />
        </Form.Item>

        <h3 className='mb-1.5'>Usable Variables</h3>
        {getEven && (
          <div className='flex gap-5 pb-5'>
            {getEven?.variables?.map((variable) => {
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
        )}
        {method === NotificationType.EMAIL && (
          <Form.Item label='Subject'>
            <Input
              onChange={(event) => setSubject(event.target.value)}
              value={subject}
            />
          </Form.Item>
        )}

        {method === NotificationType.SMS && (
          <Form.Item isList label='SMS Body' className='w-full'>
            <Input.TextArea
              onChange={(e) => setSmsBody(e.target.value)}
              rows={4}
              value={smsBody}
            />
          </Form.Item>
        )}

        {method === NotificationType.PUSH && (
          <div>
            <Form.Item label='Title'>
              <Input
                onChange={(event) => setTitle(event.target.value)}
                value={title}
              />
            </Form.Item>
            <Form.Item isList label='Body' className='w-full'>
              <Input.TextArea
                onChange={(e) => setSmsBody(e.target.value)}
                rows={4}
                value={smsBody}
              />
            </Form.Item>
            <Form.Item label='Image URL (optional)'>
              <Input
                onChange={(event) => setImageLink(event.target.value)}
                value={imageLink}
              />
            </Form.Item>
          </div>
        )}

        {method === NotificationType.EMAIL && (
          <div>
            <h3 className='pb-3'>Email Body</h3>
            {isShowEmailEditor ? (
              <EmailEditor
                ref={emailEditorRef}
                onLoad={onLoad}
                onReady={onReady}
                projectId={138679}
              />
            ) : (
              <div className='h-[500px] flex justify-center items-center'>
                <Spin size='large'/>
              </div>
            )}
          </div>
        )}
        <div className='flex items-center justify-between'>
          <div className='flex items-center'>
            <label htmlFor='isActive' className='cursor-pointer pr-2'>
              Active
            </label>
            <Switch
              onChange={setIsActive}
              id='isActive'
              checked={isActive}
              className='bg-gray-500'
            />
            <div className='ml-5'>
              <label className='pr-2'>
                Language
              </label>
              <Select dropdownMatchSelectWidth={false} optionLabelProp="value" defaultValue={"en"}
                      options={Languages_list} value={language}
                      onChange={setLanguage}/>
            </div>
          </div>

          <Button
            type='primary'
            className='bg-blue-700 block mb-5 mt-5'
            htmlType='submit'
            onClick={() => {
              formRef.current?.submit();
            }}
          >
            Update
          </Button>
        </div>
      </Form>
    </div>
  );
};

export default EditTemplate;
