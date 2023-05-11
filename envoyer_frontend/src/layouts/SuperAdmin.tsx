import React from "react";
import {Button, Layout, theme} from "antd";
import {LogoutOutlined, SendOutlined} from "@ant-design/icons";
import {signOut} from "next-auth/react";

const {Header, Content, Footer, Sider} = Layout;

export default function SuperAdmin({
                                     children,
                                   }: {
  children: React.ReactNode;
}) {
  const {
    token: {colorBgContainer},
  } = theme.useToken();

  return (
    <Layout style={{minHeight: '100vh'}}>
      <Header style={{padding: 0, background: colorBgContainer}}>
        <div className='flex text-xl whitespace-nowrap mt-4 ml-4 justify-between'>
          <div className='flex gap-3'>
            <div className='transform -rotate-45 relative bottom-2'>
              <SendOutlined/>
            </div>
            <div className='font-semibold'>Envoyer</div>
          </div>
          <div className='mr-4'>
            <Button type='default' className='flex items-center' icon={<LogoutOutlined/>}
                    onClick={() => signOut({redirect: false})}>
              Logout
            </Button>
          </div>
        </div>
      </Header>
      <Content style={{margin: '0 16px'}}>
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
      <Footer style={{textAlign: 'center'}}>
        Envoyer ©2023 Created by Vivasoft
      </Footer>
    </Layout>
  );
}
