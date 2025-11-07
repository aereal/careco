import { MonthlyReport } from '@/components/MonthlyReport';
import { FC } from 'react';

const Page: FC<PageProps<'/reports/[year]/[month]'>> = async ({ params }) => {
  const { year, month } = await params;
  return <MonthlyReport year={year} month={month} />;
};

export default Page;
