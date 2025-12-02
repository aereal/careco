import { type FC } from 'react';

export const LogoutButton: FC = () => (
  <a href='/auth/logout' className='btn' role='button'>
    Logout
  </a>
);
