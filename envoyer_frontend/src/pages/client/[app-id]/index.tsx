import {useRouter} from 'next/router';
import React, {ReactNode, useState} from 'react';
import Admin from '@/layouts/Admin';
import {USER_ROLE} from '@/utils/Constants';
import {GetServerSidePropsContext} from 'next';
import {getSession} from 'next-auth/react';
import {Button, Layout, Modal, theme} from 'antd';
import ClientList from '@/components/client/ClientList';
import ClientCreate from '@/components/client/ClientCreate';

const {Header} = Layout;

const ClientId = () => {
  const router = useRouter();
  const [showCreateClientModal, setShowCreateClientModal] = useState(false);

  const hideModal = () => {
    setShowCreateClientModal(false);
  };

  const appId = router.query?.['app-id'];

  return (
    <div>
      <div className='text-right'>
        <Button
          type='primary'
          className='bg-blue-700'
          onClick={() => setShowCreateClientModal(true)}
        >
          Create Client
        </Button>
      </div>
      <div>
        <h2 className='text-2xl font-semibold pb-2'>Clients</h2>
        <ClientList appId={Number(appId)}/>
      </div>
      <div>
        <Modal
          centered
          open={showCreateClientModal}
          title='Create Client'
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className='pt-4'>
            <ClientCreate hideModal={hideModal}></ClientCreate>
          </div>
        </Modal>
      </div>
    </div>
  );
};

export default ClientId;

const HeaderComponent = () => {
  const {
    token: {colorBgContainer},
  } = theme.useToken();
  return (
    <Header style={{padding: 0, background: colorBgContainer}}>
      <div className='pl-10'>This is the Header from client</div>
    </Header>
  );
};

ClientId.pageOptions = {
  requiresAuth: true,
  getLayout: (children: ReactNode) => (
    <Admin /*header={<HeaderComponent/>}*/>{children}</Admin>
  ),
};

export async function getServerSideProps(context: GetServerSidePropsContext) {
  const session = await getSession(context);

  if (
    session?.role === USER_ROLE.SUPER_ADMIN ||
    (session?.role === USER_ROLE.DEVELOPER &&
      session?.app_id === Number(context.params?.['app-id']))
  ) {
    return {
      props: {},
    };
  } else {
    return {
      redirect: {
        destination: '/dashboard/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
