'use client';

import { getClient } from '@/get-client';
import { runMutation } from '@/run-operation';
import { Result } from '@praha/byethrow';
import { format } from 'date-fns/format';
import { formatISO } from 'date-fns/formatISO';
import { useAtom, useAtomValue, useSetAtom } from 'jotai';
import { type ChangeEventHandler, type FC, type FormEventHandler } from 'react';
import * as atoms from './atoms';
import { modalOpen } from './atoms';
import { mutationRecordDrive } from './mutation.record';

const doRecordDrive = runMutation(getClient(), mutationRecordDrive);

export const RecordDialog: FC = () => {
  const [distance, setDistance] = useAtom(atoms.distance);
  const [date, setDate] = useAtom(atoms.date);
  const [active, setActive] = useAtom(atoms.active);
  const [memo, setMemo] = useAtom(atoms.memo);
  const setModalOpen = useSetAtom(atoms.modalOpen);
  const setError = useSetAtom(atoms.error);

  const handleChangeDistance: ChangeEventHandler<HTMLInputElement> = (e) => {
    setDistance(e.target.valueAsNumber);
  };
  const handleChangeDate: ChangeEventHandler<HTMLInputElement> = (e) => {
    const date = e.target.valueAsDate;
    if (!date) {
      return;
    }
    setDate(date);
  };
  const handleChangeMemo: ChangeEventHandler<HTMLTextAreaElement> = (e) => {
    setMemo(e.target.value);
  };
  const handleSubmit: FormEventHandler<HTMLFormElement> = async (e) => {
    e.preventDefault();
    setError(undefined);
    setActive(false);
    const ret = await doRecordDrive({
      distance,
      memo: memo === '' ? undefined : memo,
      date: formatISO(date),
    });
    const success =
      Result.isSuccess(ret) && ret.value.data
        ? ret.value.data.recordDrivingRecord
        : false;
    if (success) {
      setDistance(0);
      setMemo('');
      setModalOpen(false);
    }
    if (Result.isFailure(ret)) {
      setError(ret.error.message);
    }
    setActive(true);
  };

  return (
    <dialog className='modal' open={useAtomValue(modalOpen)}>
      <div className='modal-box'>
        <form onSubmit={handleSubmit}>
          <label className='input mt-1'>
            走行距離{' '}
            <input
              type='number'
              className='grow text-right'
              value={distance}
              onChange={handleChangeDistance}
              disabled={!active}
            />
            <span className='label'>km</span>
          </label>
          <label className='input my-1'>
            日時{' '}
            <input
              type='date'
              value={format(date, 'yyyy-MM-dd')}
              onChange={handleChangeDate}
              disabled={!active}
            />
          </label>
          <details className='collapse collapse-arrow my-1 p-1 bg-base-100 border-base-300 border'>
            <summary className='collapse-title text-sm pl-2 py-2'>メモ</summary>
            <div className='collapse-content px-2'>
              <textarea
                className='textarea h-24'
                value={memo}
                disabled={!active}
                onChange={handleChangeMemo}
              />
            </div>
          </details>
          <div className='modal-action'>
            <button className='btn btn-primary' disabled={!active}>
              記録する
            </button>
          </div>
        </form>
      </div>
      <form method='dialog' className='modal-backdrop'>
        <button>close record dialog</button>
      </form>
    </dialog>
  );
};
