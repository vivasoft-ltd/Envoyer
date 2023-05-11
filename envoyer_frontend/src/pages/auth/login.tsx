import Spinner from '@/components/spinner';
import { getSession, signIn, signOut, useSession } from 'next-auth/react';
import { useRouter } from 'next/router';
import React, { useState } from 'react';
import { SendOutlined } from '@ant-design/icons';
import { Button } from 'antd';
import { toast } from 'react-toastify';
import { USER_ROLE } from '@/utils/Constants';

const LoginPage = () => {
  const { data: session, status } = useSession();
  const router = useRouter();
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    console.log('submit');
    const data = {
      redirect: false,
      username,
      password,
    };

    const { error, ok }: any = await signIn('credentials', data);

    if (!error && ok) {
      console.log('login success');
      // router.push('/');
    } else {
      toast.error(error);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    if (name === 'username') {
      setUsername(value);
    } else if (name === 'password') {
      setPassword(value);
    }
  };

  // const eventHero = useQuery(['eventHero'], () =>
  //   authServices.getAllBlogPosts()
  // );

  if (status === 'loading') return <Spinner.FullPage />;

  return (
    <div>
      <section className='h-screen'>
        <div className='px-6 h-full text-gray-800'>
          <div className='flex gap-3 text-xl whitespace-nowrap mt-4 ml-4'>
            <div className='transform -rotate-45 relative bottom-2'>
              <SendOutlined />
            </div>
            <div className='font-semibold'>Envoyer</div>
          </div>
          <div className='flex xl:justify-center lg:justify-between justify-center items-center flex-wrap h-full g-6'>
            <div className='grow-0 shrink-1 md:shrink-0 basis-auto xl:w-6/12 lg:w-6/12 md:w-9/12 mb-12 md:mb-0'>
              <img
                src='https://mdbcdn.b-cdn.net/img/Photos/new-templates/bootstrap-login-form/draw2.webp'
                className='w-full'
                alt='Sample image'
              />
            </div>
            <div className='xl:ml-20 xl:w-4/12 lg:w-5/12 md:w-8/12 mb-12 md:mb-0'>
              <form onSubmit={handleSubmit}>
                <div className='mb-6'>
                  <input
                    type='text'
                    name='username'
                    className='form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none'
                    placeholder='Username'
                    onChange={handleChange}
                  />
                </div>

                <div className='mb-6'>
                  <input
                    type='password'
                    name='password'
                    className='form-control block w-full px-4 py-2 text-xl font-normal text-gray-700 bg-white bg-clip-padding border border-solid border-gray-300 rounded transition ease-in-out m-0 focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none'
                    placeholder='Password'
                    onChange={handleChange}
                  />
                </div>

                <div className='text-center lg:text-left'>
                  <button
                    type='submit'
                    className='inline-block px-7 py-3 bg-blue-600 text-white font-medium text-sm leading-snug uppercase rounded shadow-md hover:bg-blue-700 hover:shadow-lg focus:bg-blue-700 focus:shadow-lg focus:outline-none focus:ring-0 active:bg-blue-800 active:shadow-lg transition duration-150 ease-in-out'
                  >
                    Login
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
};

export default LoginPage;

LoginPage.pageOptions = {
  redirectIfAuthenticated: true,
  requiresAuth: false,
};

export async function getServerSideProps(context: any) {
  const session = await getSession(context);
  const isSuperAdmin = session?.role === USER_ROLE.SUPER_ADMIN;
  const appId = session?.app_id;

  if (!session) {
    return {
      props: {},
    };
  } else {
    return {
      redirect: {
        destination: isSuperAdmin
          ? '/super-admin/dashboard'
          : `/dashboard/${appId}`,
        permanent: true,
      },
    };
  }
}
