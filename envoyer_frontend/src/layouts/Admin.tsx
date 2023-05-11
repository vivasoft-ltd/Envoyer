import React, { ReactNode, useEffect, useState } from 'react';
import { Breadcrumb, Layout, Menu, theme } from 'antd';
import { signOut, useSession } from 'next-auth/react';
import Link from 'next/link';
import {
  DesktopOutlined,
  LogoutOutlined,
  PieChartOutlined,
  UserOutlined,
  SendOutlined,
  ProfileOutlined,
  FormOutlined,
  ApiOutlined, BarChartOutlined,
} from '@ant-design/icons';
import { useRouter } from 'next/router';
import { CollapseType } from 'antd/es/layout/Sider';
import { USER_ROLE } from '@/utils/Constants';
import { useQuery } from 'react-query';
import AppServices from '@/services/appServices';

const { Header, Content, Footer, Sider } = Layout;

const Admin = ({
  children,
  header,
}: {
  children: ReactNode;
  header?: JSX.Element;
}) => {
  const router = useRouter();
  console.log({ router });
  const [selectedMenu, setSelectedMenu] = useState<any>();

  useEffect(() => {
    setSelectedMenu(router.pathname.split('/')[1] || 'dashboard');
  }, [router.pathname]);

  const { data: session } = useSession();
  const [collapsed, setCollapsed] = useState(false);
  const appServices = new AppServices();
  const {
    token: { colorBgContainer },
  } = theme.useToken();

  const appId = session?.app_id || router.query?.['app-id'];

  const { data: appResp } = useQuery(['getApp'], () =>
    appServices.getApp(Number(appId))
  );

  const dashboardUrl = `/dashboard/${appId}`;

  const logoUrl = !session?.app_id
    ? '/super-admin/dashboard'
    : `/dashboard/${appId}`;

  const handleOnCollapse = (collapsed: boolean, type: CollapseType) => {
    setCollapsed(collapsed);
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider collapsible collapsed={collapsed} onCollapse={handleOnCollapse}>
        <Link
          href={logoUrl}
          onClick={() => {
            setSelectedMenu('dashboard');
          }}
        >
          <div className='p-4 m-4'>
            <div className='flex items-center justify-center gap-3 text-white text-xl whitespace-nowrap'>
              <div className='transform -rotate-45 relative bottom-2'>
                <SendOutlined />
              </div>
              {!collapsed && <span>Envoyer</span>}
            </div>
          </div>
        </Link>
        <Menu
          theme='dark'
          mode='inline'
          selectedKeys={[selectedMenu]}
          onSelect={(info) => setSelectedMenu(info.key)}
        >
          <Menu.Item key='dashboard'>
            <PieChartOutlined />
            <span>Dashboard</span>
            <Link href={dashboardUrl} />
          </Menu.Item>

          {session?.role !== USER_ROLE.ADMIN && (
            <Menu.Item key='user'>
              <UserOutlined />
              <span>User</span>
              <Link href={`/user/${appId}`} />
            </Menu.Item>
          )}

          {session?.role !== USER_ROLE.ADMIN && (
            <Menu.Item key='client'>
              <DesktopOutlined />
              <span>Client</span>
              <Link href={`/client/${appId}`} />
            </Menu.Item>
          )}

          {session?.role !== USER_ROLE.ADMIN && (
            <Menu.Item key='event'>
              <ProfileOutlined />
              <span>Event</span>
              <Link href={`/event/${appId}`} />
            </Menu.Item>
          )}

          <Menu.Item key='template'>
            <FormOutlined />
            <span>Template</span>
            <Link href={`/template/${appId}`} />
          </Menu.Item>

          <Menu.Item key='provider'>
            <ApiOutlined />
            <span>Provider</span>
            <Link href={`/provider/${appId}`} />
          </Menu.Item>

          {session?.role !== USER_ROLE.ADMIN && (
            <Menu.Item key='report'>
              <BarChartOutlined />
              <span>Report</span>
              <Link href={`/report/${appId}`} />
            </Menu.Item>
          )}

          <Menu.Item key='logout' onClick={() => signOut({ redirect: false })}>
            <LogoutOutlined />
            <span>Logout</span>
          </Menu.Item>
        </Menu>
      </Sider>
      <Layout className='site-layout'>
        {header || (
          <Header
            style={{
              padding: 0,
              background: colorBgContainer,
              // boxShadow: '0 0 10px 0 rgba(0,0,0,0.2)',
            }}
          >
            <div className='pl-10 font-semibold'>
              {appResp?.name ? appResp?.name : 'Envoyer'}
            </div>
          </Header>
        )}
        <Content style={{ margin: '0 16px' }}>
          {/* <Breadcrumb style={{ margin: '16px 0' }}>
            <Breadcrumb.Item>User</Breadcrumb.Item>
            <Breadcrumb.Item>Bill</Breadcrumb.Item>
          </Breadcrumb> */}
          <div
            style={{
              padding: 24,
              minHeight: 360,
              background: colorBgContainer,
              margin: '16px 0',
              borderRadius: 8,
              boxShadow: '0 0 10px 0 rgba(0, 0, 0, 0.1)',
            }}
          >
            {children}
          </div>
        </Content>
        <Footer style={{ textAlign: 'center' }}>
          Envoyer ©2023 Created by Vivasoft
        </Footer>
      </Layout>
    </Layout>
  );
};

export default Admin;
