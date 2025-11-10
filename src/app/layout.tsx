import { type Metadata } from 'next';
import { FC, PropsWithChildren } from 'react';
import './global.css';

export const metadata: Metadata = { title: 'Careco' };

const Layout: FC<PropsWithChildren> = ({ children }) => (
  <html lang='ja'>
    <body>{children}</body>
  </html>
);

export default Layout;
