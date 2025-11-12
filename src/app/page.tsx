import { DrivingRecordList } from '@/components/DrivingRecordList';
import { RecordDialogContainer } from '@/components/RecordDialog';
import { SelectMonth } from '@/components/SelectMonth';
import { TotalDrivingDistance } from '@/components/TotalDrivingDistance';
import { getClient } from '@/get-client';
import { readQuery } from '@/run-operation';
import { Result } from '@praha/byethrow';
import { FC } from 'react';
import { GetRoot } from './query.get-root';

const Page: FC = async () => {
  const ret = await readQuery(getClient(), GetRoot)({ first: 5 });
  if (Result.isFailure(ret)) {
    return <>Error: {ret.error.message}</>;
  }
  const {
    value: {
      data: { totalStatistics, recentDrivingRecords },
    },
  } = ret;
  return (
    <div className='max-w-2xl mx-auto'>
      <div className='p-4'>
        <TotalDrivingDistance {...totalStatistics} />
        <div className='my-8'>
          <h1 className='font-bold text-lg -mb-4'>最近の記録</h1>
          <DrivingRecordList {...recentDrivingRecords} />
          <div className='my-4'>
            <h2 className='font-bold text-md mb-2'>月毎の記録を見る</h2>
            <SelectMonth />
          </div>
        </div>
        <RecordDialogContainer />
      </div>
    </div>
  );
};

export default Page;
