import { Button, List, Typography } from 'antd';
import { ExclamationCircleFilled } from '@ant-design/icons';
import { Modal } from 'antd';
import { useMutation, useQueryClient } from 'react-query';
import { toast } from 'react-toastify';
import React, { useState } from 'react';
import ClientService from '@/services/clientService';
import ClientCreate from '@/components/client/ClientCreate';
import Link from 'next/link';
import ClientInfo from '@/components/client/ClientInfo';

export default function ClientItem({
  clientDetails,
}: {
  clientDetails: ClientResponse;
}) {
  const { confirm } = Modal;
  const clientService = new ClientService();
  const queryClient = useQueryClient();

  const [showEditClientModal, setShowEditClientModal] = useState(false);
  const [showDetails, setShowDetails] = useState(false);

  const hideDetailsModal = () => {
    setShowDetails(false);
  };

  const hideModal = () => {
    setShowEditClientModal(false);
  };

  const { mutate: deleteClient } = useMutation(
    (id: number) => clientService.deleteClient(id),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {},
    }
  );

  function handleDeleteEvent(id: number) {
    deleteClient(id, {
      onError: (err: any) => {
        toast.error(err?.message ? err?.message : 'Error deleting client');
      },
      onSuccess: async () => {
        toast.success('Client deleted successfully');
        await queryClient.invalidateQueries(['getAllClients']);
      },
    });
  }

  const showDeleteConfirm = (id: number) => {
    confirm({
      title: 'Confirm delete client',
      icon: <ExclamationCircleFilled />,
      content: 'Are you sure you want to delete this client?',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      closable: true,
      onOk() {
        handleDeleteEvent(id);
      },
      onCancel() {},
    });
  };

  return (
    <div style={{ border: '1px solid rgba(5, 5, 5, 0.06)' }}>
      <div>
        <Modal
          centered
          open={showEditClientModal}
          title='Edit Client'
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className='pt-4'>
            <ClientCreate
              IsEdit={true}
              EditClientDetails={clientDetails}
              hideModal={hideModal}
            />
          </div>
        </Modal>
        <Modal
          centered
          open={showDetails}
          title='Client Info'
          onCancel={hideDetailsModal}
          footer={null}
          width={'50%'}
        >
          <div className='pt-4'>
            <ClientInfo clientDetails={clientDetails} />
          </div>
        </Modal>
      </div>

      <List.Item>
        <List.Item.Meta
          title={
            <Typography.Text
              className='hover:text-blue-500 cursor-pointer'
              onClick={() => setShowDetails(true)}
            >
              {clientDetails.name}
            </Typography.Text>
          }
          description={clientDetails.description}
        ></List.Item.Meta>
        <div className='space-x-2'>
          <Button
            type='primary'
            ghost
            onClick={() => setShowEditClientModal(true)}
          >
            Edit
          </Button>
          <Button
            type='primary'
            danger
            ghost
            onClick={() => showDeleteConfirm(clientDetails.ID)}
          >
            Delete
          </Button>
        </div>
      </List.Item>
    </div>
  );
}
