import {useRouter} from 'next/router';
import React, {ReactNode, useState} from 'react';
import Admin from '@/layouts/Admin';
import {USER_ROLE} from '@/utils/Constants';
import {GetServerSidePropsContext} from 'next';
import {getSession} from 'next-auth/react';
import {Button, Modal} from "antd";
import ProviderList from "@/components/provider/ProviderList";
import ProviderCreate from "@/components/provider/ProviderCreate";

const ProviderId = () => {
  const router = useRouter();
  const [showCreateProviderModal, setShowCreateProviderModal] = useState(false);

  const hideModal = () => {
    setShowCreateProviderModal(false);
  };

  const appId = router.query?.['app-id'];

  return (
    <div>
      <div className='text-right'>
        <Button
          type='primary'
          className='bg-blue-700'
          onClick={() => setShowCreateProviderModal(true)}
        >
          Create Provider
        </Button>
      </div>
      <div>
        <h2 className='text-2xl font-semibold pb-2'>Providers</h2>
        <ProviderList appId={Number(appId)}/>
      </div>
      <div>
        <Modal
          centered
          open={showCreateProviderModal}
          title='Create Provider'
          onCancel={hideModal}
          footer={null}
          maskClosable={false}
        >
          <div className='pt-4'>
            <ProviderCreate hideModal={hideModal}></ProviderCreate>
          </div>
        </Modal>
      </div>
    </div>
  );
};

export default ProviderId;

ProviderId.pageOptions = {
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
        destination: '/provider/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
