import { FragmentType, getFragmentData } from '@/graphql';
import { type FC } from 'react';
import { fragmentTotalDistance } from './fragment.total-distance';

export const TotalDrivingDistance: FC<
  FragmentType<typeof fragmentTotalDistance>
> = (props) => {
  const { distanceKilometers } = getFragmentData(fragmentTotalDistance, props);
  return (
    <div className='stats shadow mt-8'>
      <div className='stat pl-4'>
        <h2 className='stat-title'>総走行距離</h2>
        <div className='stat-value'>{distanceKilometers}km</div>
      </div>
    </div>
  );
};
