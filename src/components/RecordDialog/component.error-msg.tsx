'use client';

import { useAtom } from 'jotai';
import { type FC } from 'react';
import { error as errorAtom } from './atoms';

export const ErrorMsg: FC = () => {
  const [error, setError] = useAtom(errorAtom);
  if (!error) {
    return null;
  }
  const handleClickClose = (): void => {
    setError(undefined);
  };
  return (
    <div className='toast bottom-24'>
      <div className='alert alert-error' role='alert'>
        <span>{error}</span>
        <div>
          <button className='btn btn-xs' onClick={handleClickClose}>
            x
          </button>
        </div>
      </div>
    </div>
  );
};
