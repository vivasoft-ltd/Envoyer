import React, {useEffect, useState} from 'react';
import {useMutation, useQueryClient} from 'react-query';
import AppServices from '@/services/appServices';
import {toast} from 'react-toastify';
import {Button} from "antd";

export default function AppCreate({
                                    hideModal,
                                    IsEdit = false,
                                    EditAppDetails,
                                  }: {
  hideModal: () => void;
  IsEdit?: boolean;
  EditAppDetails?: ListAppResponse;
}) {
  const queryClient = useQueryClient();
  const appServices = new AppServices();

  const [nameError, setNameError] = useState('');
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [active, setActive] = useState(false);

  useEffect(() => {
    if (IsEdit && EditAppDetails != undefined) {
      setActive(EditAppDetails?.active);
      setName(EditAppDetails?.name);
      setDescription(EditAppDetails?.description);
    }
  }, [IsEdit, EditAppDetails]);

  const clearForm = () => {
    setName('');
    setDescription('');
    setActive(false);
    setNameError('');
  };

  function isValidName(name: string) {
    let regexp = new RegExp('^[A-Za-z0-9_-]+$');
    return regexp.test(name);
  }

  function handleCreateApp(e: React.MouseEvent<HTMLAnchorElement, MouseEvent> | React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    if (!isValidName(name)) {
      setNameError(
        'Please enter a valid app name. Use letters, numbers, underscores or hyphens'
      );
      e.preventDefault();
      return;
    }
    e.preventDefault();
    createApp(
      {name, description, active},
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error creating app');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllApps']);
          toast.success('App created successfully');
          clearForm();
        },
      }
    );
  }

  function handleEditApp(e: React.MouseEvent<HTMLAnchorElement, MouseEvent> | React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();
    editApp(
      {name, description, active},
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error editing app');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllApps']);
          toast.success('App updated successfully');
          clearForm();
        },
      }
    );
  }

  const {mutate: editApp} = useMutation(
    (data: CreateAppInput) => appServices.editApp(data, EditAppDetails?.ID),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  const {mutate: createApp} = useMutation(
    (data: CreateAppInput) => appServices.createApp(data),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const {name, value, checked} = e.target;
    if (name === 'name') {
      setName(value);
    } else if (name === 'active') {
      setActive(checked);
    }
  };

  const handleTextAreaChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    const {name, value} = e.target;
    if (name === 'description') {
      setDescription(value);
    }
  };

  return (
    <div className=''>
      <form className='w-full max-w-lg pb-3'>
        <div className='flex flex-wrap -mx-3 mb-6'>
          <div className='w-full px-3 mb-6 md:mb-0'>
            <label
              className='block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2'
              htmlFor='grid-first-name'
            >
              App Name
            </label>
            <input
              className='appearance-none block w-full text-gray-700 border rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white disabled:bg-gray-200'
              id='grid-first-name'
              type='text'
              name='name'
              placeholder='Name'
              value={name}
              disabled={IsEdit}
              onChange={handleChange}
            />
            <p className='text-red-500 text-xs italic'>{nameError}</p>
          </div>
        </div>
        <div className='flex flex-wrap -mx-3 mb-6'>
          <div className='w-full px-3'>
            <label
              className='block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2'
              htmlFor='grid-last-name'
            >
              Description
            </label>
            <textarea
              className='appearance-none block w-full text-gray-700 border border-gray-200 rounded py-3 px-4 leading-tight focus:outline-none focus:bg-white focus:border-gray-500'
              id='grid-last-name'
              placeholder='Description'
              name='description'
              value={description}
              onChange={handleTextAreaChange}
            />
          </div>
        </div>
        <div className='flex flex-wrap -mx-3 mb-6'>
          <div className='w-full px-3'>
            <label className='relative inline-flex items-center cursor-pointer'>
              <input
                type='checkbox'
                value=''
                checked={active}
                name='active'
                className='sr-only peer'
                onChange={handleChange}
              />
              <div
                className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer dark:bg-gray-700 peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all dark:border-gray-600 peer-checked:bg-blue-600"></div>
              <span className='ml-3 block uppercase tracking-wide text-gray-700 text-xs font-bold'>
                Active
              </span>
            </label>
          </div>
        </div>
        <div className='text-center lg:text-right'>
          <Button
            type='primary'
            onClick={(e) => {
              if (IsEdit) {
                handleEditApp(e);
              } else {
                handleCreateApp(e);
              }
            }}
            className='bg-blue-700'
          >
            {IsEdit ? 'Update' : 'Create'}
          </Button>
        </div>
      </form>
    </div>
  );
}
