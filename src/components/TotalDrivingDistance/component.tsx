import { FragmentType, getFragmentData } from '@/graphql';
import { type FC } from 'react';
import { fragmentTotalDistance } from './fragment.total-distance';

export const TotalDrivingDistance: FC<
  FragmentType<typeof fragmentTotalDistance>
> = (props) => {
  const { distanceKilometers } = getFragmentData(fragmentTotalDistance, props);
  return (
    <div className='stats'>
      <div className='stat'>
        <h2 className='stat-title'>総走行距離</h2>
        <div className='stat-value'>{distanceKilometers}km</div>
      </div>
    </div>
  );
};
