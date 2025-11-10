import { RecordDialog } from '@/components/RecordDialog';
import { SelectMonth } from '@/components/SelectMonth';
import { YMD } from '@/components/YMD';
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
        <div className='stats'>
          <div className='stat'>
            <h2 className='stat-title'>総走行距離</h2>
            <div className='stat-value'>
              {totalStatistics.distanceKilometers}km
            </div>
          </div>
        </div>
        <ol className='mb-8 list rounded-box shadow-md'>
          {recentDrivingRecords.map((record) => (
            <li key={record.recordedAt.toString()} className='p-2'>
              <YMD date={record.recordedAt} />: {record.distanceKilometers}km
            </li>
          ))}
        </ol>
        <RecordDialog />
      </div>
    </div>
  );
};

export default Page;
