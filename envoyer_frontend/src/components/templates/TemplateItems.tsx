import {Badge, Card, Col, Modal, Popconfirm, Row} from 'antd';
import {useMutation, useQuery, useQueryClient} from 'react-query';
import TemplateService from '@/services/templateService';
import React, {useState} from 'react';
import parse from 'html-react-parser';
import {NotificationType} from '@/utils/Constants';
import {GrLanguage} from "react-icons/gr";
import EditTemplate from './EditTemplate';
import {toast} from 'react-toastify';
import {
  CheckCircleFilled,
  CloseCircleFilled,
  ExclamationCircleFilled,
  EyeOutlined,
  DeleteOutlined,
  EditOutlined,
} from '@ant-design/icons';
import Link from "next/link";

const templateService = new TemplateService();

const TemplateItems = ({
                         templates,
                         type,
                       }: {
  templates?: TemplateResponse[];
  type: string;
}) => {
  const queryClient = useQueryClient();

  const [showCreateEventModal, setShowCreateEventModal] = useState(false);
  const [showPreviewModal, setShowPreviewModal] = useState(false);
  const [selectedTemplate, setSelectedTemplate] = useState<TemplateResponse>();

  const hideModal = () => {
    setShowCreateEventModal(false);
    setShowPreviewModal(false);
  };

  const handleEdit = (template: TemplateResponse) => {
    setSelectedTemplate(template);
    setShowCreateEventModal(true);
  };

  const handlePreview = (template: TemplateResponse) => {
    setSelectedTemplate(template);
    setShowPreviewModal(true);
  };

  const {mutate: deleteTemplate} = useMutation(
    (id: number) => templateService.deleteTemplate(id),
    {
      onSuccess: async (data) => {
        console.log(data);
        toast.success('Template Deleted Successfully');
        await queryClient.invalidateQueries(['getAllTemplates']);
      },
      onError: async (error: Error) => {
        toast.error(error.message);
      },
    }
  );

  return (
    <div className='flex gap-5 flex-wrap'>
      {templates?.map((template) => {
        if (template.type !== type) return null;
        return (
          <div key={template.ID} className='w-96'>
            <Card
              actions={[
                // <SettingOutlined key='setting' />,
                <EyeOutlined
                  key='preview'
                  onClick={() => handlePreview(template)}
                />,
                <EditOutlined
                  key='edit'
                  onClick={() => handleEdit(template)}
                />,
                <Popconfirm
                  okButtonProps={{danger: true}}
                  title='Delete the Template?'
                  description='Are you sure to delete this Template?'
                  onConfirm={() => deleteTemplate(template.ID)}
                  // onCancel={cancel}
                  okText='Yes'
                  cancelText='No'
                  key='delete'
                >
                  <DeleteOutlined/>
                </Popconfirm>,
              ]}
              title={
                <div className='flex gap-2'>
                  {template.active ? (
                    <CheckCircleFilled style={{color: '#609966'}}/>
                  ) : (
                    <CloseCircleFilled style={{color: '#FF4d4f'}}/>
                  )}
                  <h1 className='truncate pr-5'>{template.description}</h1>
                  <div className='ml-auto flex items-center gap-1'>
                    <GrLanguage/>
                    <h1>{template.language}</h1>
                  </div>
                </div>
              }
            >
              <p className='line-clamp-4 h-[90px] whitespace-pre-wrap'>
                {type === NotificationType.EMAIL
                  ? template?.email_subject
                  : type === NotificationType.PUSH
                    ? template?.title
                    : template.message}
              </p>
            </Card>
          </div>
        );
      })}
      <Modal
        centered
        open={showCreateEventModal}
        title='Edit Template'
        onCancel={hideModal}
        footer={null}
        maskClosable={false}
        width={type === NotificationType.EMAIL ? '90%' : 600}
        destroyOnClose
      >
        <Col span={24}>
          <EditTemplate
            method={type}
            template={selectedTemplate}
            hideModal={hideModal}
          />
        </Col>
      </Modal>

      <Modal
        centered
        open={showPreviewModal}
        title={
          <div className='flex gap-2'>
            Template Preview
            <div className='ml-auto flex items-center relative right-8 gap-1'>
              <GrLanguage/>
              <h1>{selectedTemplate?.language}</h1>
            </div>
          </div>
        }
        onCancel={hideModal}
        footer={null}
        maskClosable={true}
        closable
        width={600}
      >
        <Col span={24}>
          {type === NotificationType.EMAIL && (
            <h1 className='font-semibold text-lg py-3'>
              Email Subject: {selectedTemplate?.email_subject}
            </h1>
          )}
          <p className='text-base pb-3'>
            Description: {selectedTemplate?.description}
          </p>

          {type === NotificationType.SMS ? (
            <p className='bg-[#e7e7e7] p-3 whitespace-pre-wrap'>
              {selectedTemplate?.message}
            </p>
          ) : type === NotificationType.PUSH ? (
            <div>
              <p className='pb-2 font-semibold whitespace-pre-wrap'>
                Title : {selectedTemplate?.title}
              </p>
              {selectedTemplate?.file && <p className='pb-2 font-semibold whitespace-pre-wrap'>
                  Image URL : <Link target="_blank" href={selectedTemplate?.file}>{selectedTemplate?.file}</Link>
              </p>}
              <p className='bg-[#e7e7e7] p-3 whitespace-pre-wrap'>
                {selectedTemplate?.message}
              </p>
            </div>
          ) : (
            <div className=''>
              {parse(selectedTemplate?.email_rendered_html || '')}
            </div>
          )}
        </Col>
      </Modal>
    </div>
  );
};

export default TemplateItems;
