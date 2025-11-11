'use client';

import { useSetAtom } from 'jotai';
import { type FC, type MouseEventHandler } from 'react';
import { modalOpen } from './atoms';

export const OpenButton: FC = () => {
  const setModalOpen = useSetAtom(modalOpen);
  const handleClick: MouseEventHandler<HTMLButtonElement> = (e) => {
    e.preventDefault();
    setModalOpen((v) => !v);
  };
  return (
    <div className='fab'>
      <button
        className='btn btn-lg btn-circle btn-primary'
        onClick={handleClick}
      >
        +
      </button>
    </div>
  );
};
