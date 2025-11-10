import { YMD } from '@/components/YMD';
import { FragmentType, getFragmentData } from '@/graphql';
import { type FC } from 'react';
import { fragmentRecordList } from './fragment.record-list';

export const DrivingRecordList: FC<FragmentType<typeof fragmentRecordList>> = (
  props,
) => {
  const { nodes } = getFragmentData(fragmentRecordList, props);
  return (
    <ol className='mb-8 list rounded-box shadow-md'>
      {nodes.map((record) => (
        <li key={record.recordedAt.toString()} className='p-2'>
          <YMD date={record.recordedAt} />: {record.distanceKilometers}km
        </li>
      ))}
    </ol>
  );
};
