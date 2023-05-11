import SuperAdmin from '@/layouts/SuperAdmin';
import {USER_ROLE} from '@/utils/Constants';
import {getSession} from 'next-auth/react';
import React, {ReactNode, useState} from 'react';
import AppCreate from '@/components/app/AppCreate';
import AppList from '@/components/app/AppList';
import {Button, Modal} from 'antd';
import {GetServerSidePropsContext} from 'next';
import {requireRole} from '@/utils/requireRole';

const SuperAdminDashboard = () => {
  const [showCreateAppModal, setShowCreateAppModal] = useState(false);

  const hideModal = () => {
    setShowCreateAppModal(false);
  };

  return (
    <div className=''>
      <div className='text-right'>
        <Button
          type="primary"
          onClick={() => setShowCreateAppModal(true)}
          className='bg-blue-700'
        >
          Create App
        </Button>
      </div>
      <div>
        <h2 className='text-2xl font-semibold pb-2'>Apps</h2>
        <AppList/>
      </div>
      <div>
        <Modal
          centered
          open={showCreateAppModal}
          title='Create App'
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className='pt-4'>
            <AppCreate hideModal={hideModal}></AppCreate>
          </div>
        </Modal>
      </div>
    </div>
  );
};

export default SuperAdminDashboard;

SuperAdminDashboard.layout = SuperAdmin;

SuperAdminDashboard.pageOptions = {
  role: USER_ROLE.SUPER_ADMIN,
  requiresAuth: true,
  getLayout: (children: ReactNode) => <SuperAdmin>{children}</SuperAdmin>,
};

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getSession(context);
  requireRole('super-admin');

  const appId = session?.app_id;

  if (session?.role === USER_ROLE.SUPER_ADMIN) {
    return {
      props: {},
    };
  } else {
    return {
      redirect: {
        destination: '/dashboard/' + appId,
        permanent: true,
      },
    };
  }
}
