import {useMutation, useQueryClient} from 'react-query';
import React, {useEffect, useState} from 'react';
import {toast} from 'react-toastify';
import {Button, Form, Input} from 'antd';
import {useRouter} from 'next/router';
import ClientService from "@/services/clientService";

export default function ClientCreate({
                                       hideModal,
                                       IsEdit = false,
                                       EditClientDetails,
                                     }: {
  hideModal: () => void;
  IsEdit?: boolean;
  EditClientDetails?: ClientResponse;
}) {
  const queryClient = useQueryClient();
  const clientService = new ClientService();

  const router = useRouter();
  const app_id = router.query?.['app-id'] as string;

  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  useEffect(() => {
    if (IsEdit && EditClientDetails != undefined) {
      setName(EditClientDetails?.name);
      setDescription(EditClientDetails?.description || '');
    }
  }, [IsEdit, EditClientDetails]);

  const clearForm = () => {
    setName('');
    setDescription('');
  };

  function handleCreateClient(e: React.MouseEvent<HTMLAnchorElement, MouseEvent> | React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();

    createClient(
      {name, description, app_id: Number(app_id)},
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error creating client');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllClients']);
          toast.success('Client created successfully');
          clearForm();
        },
      }
    );
  }

  function handleEditClient(e: React.MouseEvent<HTMLAnchorElement, MouseEvent> | React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();
    editClient(
      {name, description, app_id: Number(app_id)},
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error editing client');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllClients']);
          toast.success('Client updated successfully');
          clearForm();
        },
      }
    );
  }

  const {mutate: editClient} = useMutation(
    (data: CreateClientInput) =>
      clientService.editClient(data, EditClientDetails?.ID),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  const {mutate: createClient} = useMutation(
    (data: CreateClientInput) => clientService.createClient(data),
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
      <Form className='w-full max-w-lg pb-3' layout='vertical'>
        <Form.Item label='Client Name'>
          <Input
            value={name}
            onChange={(e) => setName(e.target.value)}
          />
        </Form.Item>

        <Form.Item label="Description">
          <Input.TextArea rows={4} value={description} onChange={(e) => setDescription(e.target.value)}/>
        </Form.Item>

        <div className='text-center lg:text-right'>
          <Button type="primary" className='bg-blue-700' onClick={(e) => {
            if (IsEdit) {
              handleEditClient(e);
            } else {
              handleCreateClient(e);
            }
          }}>
            {IsEdit ? 'Update' : 'Create'}
          </Button>
        </div>
      </Form>
    </div>
  );
}
