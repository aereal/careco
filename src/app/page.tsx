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
    <div>
      <div>
        <SelectMonth />
      </div>
      <div>total: {totalStatistics.distanceKilometers}km</div>
      <ul>
        {recentDrivingRecords.map((record) => (
          <li key={record.recordedAt.toString()}>
            <YMD date={record.recordedAt} />: {record.distanceKilometers}km
          </li>
        ))}
      </ul>
      <RecordDialog />
    </div>
  );
};

export default Page;
