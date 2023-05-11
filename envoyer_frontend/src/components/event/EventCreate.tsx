import { useMutation, useQueryClient } from 'react-query';
import React, { useEffect, useState } from 'react';
import { toast } from 'react-toastify';
import EventService from '@/services/eventService';
import { Button, Checkbox, Form, Input } from 'antd';
import TemplateVariables from '../templateVariables';
import { useRouter } from 'next/router';
import { transformToVariable } from '@/utils/stringFormat';

export default function EventCreate({
  hideModal,
  IsEdit = false,
  EditEventDetails,
}: {
  hideModal: () => void;
  IsEdit?: boolean;
  EditEventDetails?: EventResponse;
}) {
  const queryClient = useQueryClient();
  const eventService = new EventService();

  const router = useRouter();
  const app_id = router.query?.['app-id'] as string;

  const [nameError, setNameError] = useState('');
  const [variableError, setVariableError] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [variables, setVariables] = useState<string[] | undefined>([]);
  const [variable, setVariable] = useState('');

  useEffect(() => {
    if (IsEdit && EditEventDetails != undefined) {
      setName(EditEventDetails?.name);
      setDescription(EditEventDetails?.description);
      setVariables(EditEventDetails?.variables);
    }
  }, [IsEdit, EditEventDetails]);

  const clearForm = () => {
    setName('');
    setDescription('');
    setVariables([]);
    setVariable('');
    setNameError('');
    setVariableError('');
  };

  function isValidName(name: string) {
    let regexp = new RegExp('^[A-Za-z0-9_-]+$');
    return regexp.test(name);
  }

  function handleCreateEvent(e: React.FormEvent<HTMLButtonElement>) {
    if (!isValidName(name)) {
      setNameError(
        'Please enter a valid app name. Use letters, numbers, underscores or hyphens'
      );
      e.preventDefault();
      return;
    }
    e.preventDefault();

    let event = {
      name,
      description,
      variables: variables?.map((v) => transformToVariable(v)),
      app_id: Number(app_id),
    };
    createEvent(event, {
      onError: (err: any) => {
        toast.error(err?.message ? err?.message : 'Error creating event');
      },
      onSuccess: async () => {
        hideModal();
        await queryClient.invalidateQueries(['getAllEvents']);
        toast.success('Event created successfully');
        clearForm();
      },
    });
  }

  function handleEditEvent(e: React.FormEvent<HTMLButtonElement>) {
    if (!isValidName(name)) {
      setNameError(
        'Please enter a valid app name. Use letters, numbers, underscores or hyphens'
      );
      e.preventDefault();
      return;
    }
    e.preventDefault();
    editEvent(
      {
        name,
        description,
        variables: variables?.map((v) => transformToVariable(v)),
        app_id: Number(app_id),
      },
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error editing event');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllEvents']);
          toast.success('Event updated successfully');
          clearForm();
        },
      }
    );
  }

  const { mutate: editEvent } = useMutation(
    (data: CreateEventInput) =>
      eventService.editEvent(data, EditEventDetails?.ID),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {},
    }
  );

  const { mutate: createEvent } = useMutation(
    (data: CreateEventInput) => eventService.createEvent(data),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {},
    }
  );

  const handleKeyDown: React.KeyboardEventHandler<HTMLInputElement> = (e) => {
    if (e.key === 'Enter') {
      e.preventDefault();
      if (isValidName(variable) && !variables?.includes(variable)) {
        setVariables([...(variables || []), e.currentTarget.value]);
        setVariable('');
        setVariableError('');
      } else {
        setVariableError(
          'Please enter a valid and unique variable name. Use letters, numbers, underscores or hyphens'
        );
      }
    }
  };

  return (
    <div className=''>
      <Form className='w-full max-w-lg pb-3' layout='vertical'>
        <Form.Item
          label='Event Name'
          help={nameError}
          validateStatus={nameError ? 'error' : 'validating'}
        >
          <Input value={name} onChange={(e) => setName(e.target.value)} />
        </Form.Item>

        <Form.Item label='Description'>
          <Input.TextArea
            rows={4}
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          />
        </Form.Item>

        <Form.Item
          label='Variables'
          help={variableError}
          validateStatus={variableError ? 'error' : 'validating'}
        >
          <Input
            onKeyDown={handleKeyDown}
            value={variable}
            onChange={(e) => setVariable(e.target.value)}
            placeholder='Enter variable name and press enter'
          />
        </Form.Item>

        <div className='pt-3'>
          {variables && variables?.length > 0 && (
            <TemplateVariables
              variables={variables}
              setVariables={setVariables}
            />
          )}
        </div>

        <div className='text-center lg:text-right'>
          <button
            type='submit'
            onClick={(e) => {
              if (IsEdit) {
                handleEditEvent(e);
              } else {
                handleCreateEvent(e);
              }
            }}
            className='inline-block px-7 py-3 bg-blue-600 text-white font-medium text-sm leading-snug uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out'
          >
            {IsEdit ? 'Update' : 'Create'}
          </button>
        </div>
      </Form>
    </div>
  );
}
