import { DrivingRecordList } from '@/components/DrivingRecordList';
import { RecordDialog } from '@/components/RecordDialog';
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
        <div>
          <SelectMonth />
        </div>
        <TotalDrivingDistance {...totalStatistics} />
        <DrivingRecordList {...recentDrivingRecords} />
        <RecordDialog />
      </div>
    </div>
  );
};

export default Page;
