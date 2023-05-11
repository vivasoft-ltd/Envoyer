import { Button, Descriptions, List, Typography } from 'antd';
import { ExclamationCircleFilled } from '@ant-design/icons';
import { Modal } from 'antd';
import { useMutation, useQueryClient } from 'react-query';
import { toast } from 'react-toastify';
import React, { useState } from 'react';
import UserService from '@/services/userService';
import UserCreate from '@/components/user/UserCreate';
import Link from 'next/link';

export default function UserItem({
  userDetails,
  canEdit = true,
}: {
  userDetails: UserResponse;
  canEdit?: boolean;
}) {
  const { confirm } = Modal;
  const userService = new UserService();
  const queryClient = useQueryClient();

  const [showEditUserModal, setShowEditUserModal] = useState(false);
  const [showDetails, setShowDetails] = useState(false);

  const hideDetailsModal = () => {
    setShowDetails(false);
  };

  const hideModal = () => {
    setShowEditUserModal(false);
  };

  const { mutate: deleteUser } = useMutation(
    (id: number) => userService.deleteUser(id),
    {
      onSuccess: async (data) => {
        console.log(data);
      },
      onSettled: async () => {},
    }
  );

  function handleDeleteUser(id: number) {
    deleteUser(id, {
      onError: (err: any) => {
        toast.error(err?.message ? err?.message : 'Error deleting user');
      },
      onSuccess: async () => {
        toast.success('User deleted successfully');
        await queryClient.invalidateQueries(['getAllUsers']);
      },
    });
  }

  const showDeleteConfirm = (id: number) => {
    confirm({
      title: 'Confirm delete user',
      icon: <ExclamationCircleFilled />,
      content: 'Are you sure you want to delete this user?',
      okText: 'Delete',
      okType: 'danger',
      cancelText: 'Cancel',
      closable: true,
      onOk() {
        handleDeleteUser(id);
      },
      onCancel() {},
    });
  };

  return (
    <div style={{ border: '1px solid rgba(5, 5, 5, 0.06)' }}>
      <div>
        <Modal
          centered
          open={showEditUserModal}
          title='Edit User'
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className='pt-4'>
            <UserCreate
              IsEdit={true}
              EditUserDetails={userDetails}
              hideModal={hideModal}
            />
          </div>
        </Modal>

        <Modal
          centered
          open={showDetails}
          title='User Info'
          onCancel={hideDetailsModal}
          footer={null}
        >
          <div className='pt-4'>
            <Descriptions title='' layout='horizontal' column={1} bordered>
              <Descriptions.Item label='Username'>
                {userDetails?.user_name}
              </Descriptions.Item>
              <Descriptions.Item label='User Id'>
                {userDetails?.ID}
              </Descriptions.Item>
              <Descriptions.Item label='Role'>
                {userDetails?.role}
              </Descriptions.Item>
            </Descriptions>
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
              {userDetails.user_name}
            </Typography.Text>
          }
          description={userDetails.role}
        ></List.Item.Meta>
        {canEdit && (
          <div className='space-x-2'>
            <Button
              type='primary'
              ghost
              onClick={() => setShowEditUserModal(true)}
            >
              Edit
            </Button>
            <Button
              type='primary'
              danger
              ghost
              onClick={() => showDeleteConfirm(userDetails.ID)}
            >
              Delete
            </Button>
          </div>
        )}
      </List.Item>
    </div>
  );
}
