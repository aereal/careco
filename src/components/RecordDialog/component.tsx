'use client';

import { getClient } from '@/get-client';
import { runMutation } from '@/run-operation';
import { Result } from '@praha/byethrow';
import { format } from 'date-fns/format';
import { formatISO } from 'date-fns/formatISO';
import { useAtom } from 'jotai';
import { ChangeEventHandler, FormEventHandler, type FC } from 'react';
import * as atoms from './atoms';
import { mutationRecordDrive } from './mutation.record';

const doRecordDrive = runMutation(getClient(), mutationRecordDrive);

export const RecordDialog: FC = () => {
  const [distance, setDistance] = useAtom(atoms.distance);
  const [date, setDate] = useAtom(atoms.date);
  const [active, setActive] = useAtom(atoms.active);
  const [memo, setMemo] = useAtom(atoms.memo);

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
    }
    setActive(true);
  };
  return (
    <form onSubmit={handleSubmit}>
      <label className='inline-block'>
        走行距離{' '}
        <input
          type='number'
          value={distance}
          onChange={handleChangeDistance}
          disabled={!active}
        />
      </label>
      <label className='inline-block'>
        日時{' '}
        <input
          type='date'
          value={format(date, 'yyyy-MM-dd')}
          onChange={handleChangeDate}
          disabled={!active}
        />
      </label>
      <textarea value={memo} disabled={!active} onChange={handleChangeMemo} />
      <button type='submit' disabled={!active}>
        記録
      </button>
    </form>
  );
};
