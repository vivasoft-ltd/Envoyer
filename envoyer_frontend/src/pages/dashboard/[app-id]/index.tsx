import {USER_ROLE} from '@/utils/Constants';
import {GetServerSidePropsContext} from 'next';
import {getSession, useSession} from 'next-auth/react';
import React, {ReactNode} from 'react';
import Admin from '@/layouts/Admin';
import {useRouter} from 'next/router';
import {useQuery} from 'react-query';
import AppServices from '@/services/appServices';
import {Badge, Descriptions, Layout, theme} from "antd";
import UserInfo from "@/components/user/UserInfo";

const AppId = ({}) => {
  const router = useRouter();
  const {data: session} = useSession();
  const appServices = new AppServices();

  const appId = router.query?.['app-id'];

  const {data: appResp} = useQuery(['getApp'], () =>
    appServices.getApp(Number(appId))
  );

  return (
    <div>
      <Descriptions title="Application Information" layout="vertical" bordered>
        <Descriptions.Item label="App Name">{appResp?.name}</Descriptions.Item>
        <Descriptions.Item label="App Id">{appResp?.ID}</Descriptions.Item>
        <Descriptions.Item label="Status">
          <Badge status={appResp?.active ? "success" : "error"} text={appResp?.active ? "Active" : "Inactive"}/>
        </Descriptions.Item>
        <Descriptions.Item label="Description">
          {appResp?.description}
        </Descriptions.Item>
      </Descriptions>

      {session?.role === USER_ROLE.SUPER_ADMIN ?
        <Descriptions title="User Information" layout="vertical" bordered className="pt-5">
          <Descriptions.Item label="User Id">{session?.id}</Descriptions.Item>
          <Descriptions.Item label="Role">{session?.role}</Descriptions.Item>
        </Descriptions> :
        <UserInfo userId={session?.id as number}/>
      }
    </div>
  );
};

export default AppId;

AppId.pageOptions = {
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
        destination: '/dashboard/' + session?.app_id,
        permanent: true,
      },
    };
  }
}
