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
    <div>
      <button className='btn btn-neutral' onClick={handleClick}>
        記録する
      </button>
    </div>
  );
};
