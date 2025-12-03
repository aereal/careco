import { DriveRecordChart } from '@/components/DriveRecordChart';
import { LoginButton } from '@/components/LoginButton';
import { LogoutButton } from '@/components/LogoutButton';
import { RecordDialogContainer } from '@/components/RecordDialog';
import { SelectMonth } from '@/components/SelectMonth';
import { TotalDrivingDistance } from '@/components/TotalDrivingDistance';
import { getClient } from '@/get-client';
import { client } from '@/lib/auth0';
import { readQuery } from '@/run-operation';
import { Result } from '@praha/byethrow';
import { FC } from 'react';
import { GetRoot } from './query.get-root';

const Page: FC = async () => {
  const session = await client.getSession();
  if (!session?.user) {
    return (
      <div className='max-w-2xl mx-auto'>
        <div className='p-4'>
          <LoginButton />
        </div>
      </div>
    );
  }
  const ret = await readQuery(getClient(), GetRoot)({ first: 30 });
  if (Result.isFailure(ret)) {
    return (
      <div className='max-w-2xl mx-auto'>
        <div className='p-4'>
          <p>Error: {ret.error.message}</p>
          <LoginButton />
        </div>
      </div>
    );
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
        <div className='my-4'>
          <h2 className='font-bold text-md mb-2'>月毎の記録を見る</h2>
          <SelectMonth />
        </div>
        <div className='my-8'>
          <h1 className='font-bold text-lg -mb-4'>最近の記録</h1>
          <div className='my-8 h-[60vh]'>
            <DriveRecordChart records={recentDrivingRecords} showMonth />
          </div>
        </div>
        <RecordDialogContainer />
        <LogoutButton />
      </div>
    </div>
  );
};

export default Page;
