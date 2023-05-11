import {useRouter} from 'next/router';
import React, {ReactNode, useState} from 'react';
import Admin from '@/layouts/Admin';
import {USER_ROLE} from '@/utils/Constants';
import {GetServerSidePropsContext} from 'next';
import {getSession, useSession} from 'next-auth/react';
import {Button, Modal} from "antd";
import UserList from "@/components/user/UserList";
import UserCreate from "@/components/user/UserCreate";

const UserId = () => {
  const {data: session} = useSession();
  const router = useRouter();
  const [showCreateUserModal, setShowCreateUserModal] = useState(false);

  const hideModal = () => {
    setShowCreateUserModal(false);
  };

  const appId = router.query?.['app-id'];

  return (
    <div>
      {session?.role === USER_ROLE.SUPER_ADMIN && (
        <div className='text-right'>
          <Button
            type='primary'
            className='bg-blue-700'
            onClick={() => setShowCreateUserModal(true)}
          >
            Create User
          </Button>
        </div>
      )}

      {session?.role === USER_ROLE.SUPER_ADMIN && (
        <div>
          <h2 className='text-2xl font-semibold pb-2'>Users</h2>
          <UserList appId={Number(appId)}/>
        </div>
      )}

      {session?.role === USER_ROLE.DEVELOPER && (
        <div>
          <h2 className='text-2xl font-semibold pb-2'>Users</h2>
          <UserList appId={Number(appId)} canEdit={false}/>
        </div>
      )}

      <div>
        <Modal
          centered
          open={showCreateUserModal}
          title='Create User'
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className='pt-4'>
            <UserCreate hideModal={hideModal}></UserCreate>
          </div>
        </Modal>
      </div>
    </div>
  );
};

export default UserId;

UserId.pageOptions = {
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
        destination: '/user/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
