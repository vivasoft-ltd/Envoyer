import {useMutation, useQueryClient} from 'react-query';
import React, {useEffect, useState} from 'react';
import {toast} from 'react-toastify';
import {Button, Form, Input, Select} from 'antd';
import {useRouter} from 'next/router';
import UserService from "@/services/userService";
import {USER_ROLE} from "@/utils/Constants";

export default function UserCreate({
                                     hideModal,
                                     IsEdit = false,
                                     EditUserDetails,
                                   }: {
  hideModal: () => void;
  IsEdit?: boolean;
  EditUserDetails?: UserResponse;
}) {
  const queryClient = useQueryClient();
  const userService = new UserService();

  const router = useRouter();
  const app_id = router.query?.['app-id'] as string;

  const [username, setUserName] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState(USER_ROLE.ADMIN);

  useEffect(() => {
    if (IsEdit && EditUserDetails != undefined) {
      setUserName(EditUserDetails.user_name);
      setPassword(EditUserDetails?.password);
      setRole(EditUserDetails?.role);
    }
  }, [IsEdit, EditUserDetails]);

  const clearForm = () => {
    setRole(USER_ROLE.ADMIN);
    setUserName('');
    setPassword('');
  };

  function handleCreateUser(e: React.MouseEvent<HTMLAnchorElement, MouseEvent> | React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();

    createUser(
      {user_name: username, password, role, app_id: Number(app_id)},
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error creating user');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllUsers']);
          toast.success('User created successfully');
          clearForm();
        },
      }
    );
  }

  function handleEditUser(e: React.MouseEvent<HTMLAnchorElement, MouseEvent> | React.MouseEvent<HTMLButtonElement, MouseEvent>) {
    e.preventDefault();
    editUser(
      {user_name: username, password, role, app_id: Number(app_id)},
      {
        onError: (err: any) => {
          toast.error(err?.message ? err?.message : 'Error editing user');
        },
        onSuccess: async () => {
          hideModal();
          await queryClient.invalidateQueries(['getAllUsers']);
          toast.success('User updated successfully');
          clearForm();
        },
      }
    );
  }

  const {mutate: editUser} = useMutation(
    (data: CreateUserInput) =>
      userService.editUser(data, EditUserDetails?.ID),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {
      },
    }
  );

  const {mutate: createUser} = useMutation(
    (data: CreateUserInput) => userService.createUser(data),
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
        <Form.Item label='Username'>
          <Input
            value={username}
            onChange={(e) => setUserName(e.target.value)}
          />
        </Form.Item>

        <Form.Item label='Password'>
          <Input.Password
            value={password}
            type="password"
            onChange={(e) => setPassword(e.target.value)}
          />
        </Form.Item>

        <Form.Item label="Role">
          <Select
            value={role}
            onChange={(value) => setRole(value)}
            options={[
              {value: USER_ROLE.DEVELOPER, label: 'Developer'},
              {value: USER_ROLE.ADMIN, label: 'Application Owner'},
            ]}
          />
        </Form.Item>

        <div className='text-center lg:text-right'>
          <Button type="primary" className='bg-blue-700' onClick={(e) => {
            if (IsEdit) {
              handleEditUser(e);
            } else {
              handleCreateUser(e);
            }
          }}>
            {IsEdit ? 'Update' : 'Create'}
          </Button>
        </div>
      </Form>
    </div>
  );
}
