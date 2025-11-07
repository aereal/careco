import { YMD } from '@/components/YMD';
import { getClient } from '@/get-client';
import { FC } from 'react';
import { GetRoot } from './query.get-root';

const Page: FC = async () => {
  const { data, error } = await getClient().query(GetRoot, { first: 5 });
  if (error) {
    return <>Error: {error.message}</>;
  }
  if (!data) {
    return null;
  }
  return (
    <div>
      <div>total: {data.totalStatistics.distanceKilometers}km</div>
      <ul>
        {data.recentDrivingRecords.map((record) => (
          <li key={record.recordedAt.toString()}>
            <YMD date={record.recordedAt} />: {record.distanceKilometers}km
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Page;
