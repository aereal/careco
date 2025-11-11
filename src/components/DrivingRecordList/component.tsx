import { YMD } from '@/components/YMD';
import { FragmentType, getFragmentData } from '@/graphql';
import { type FC } from 'react';
import { fragmentRecordList } from './fragment.record-list';

export const DrivingRecordList: FC<FragmentType<typeof fragmentRecordList>> = (
  props,
) => {
  const { nodes } = getFragmentData(fragmentRecordList, props);
  return (
    <table className='table'>
      <thead>
        <tr>
          <th>日付</th>
          <th>走行距離</th>
        </tr>
      </thead>
      <tbody>
        {nodes.map(({ recordedAt, distanceKilometers }) => (
          <tr key={recordedAt.toString()}>
            <td>
              <YMD date={recordedAt} />
            </td>
            <td>{distanceKilometers}km</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};
